package uberproxy

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/client/hack"
	"vtb-apihack-2025/internal/config"
	"vtb-apihack-2025/internal/otp"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"moul.io/http2curl"
)

// var proxyRoutes = []string{}

type server struct {
	jwtSecret []byte
	jwtMethod jwt.SigningMethod

	// TODO: don't leak implementation
	apis map[int64]*hack.ApiClient

	corsOrigin string
	otp        otp.OTPAuthenticator

	users interfaces.UserStore
	// banks interfaces.BankStore
	cfg config.Config

	isdebug          bool
	makeConsentStore func(bankId int64) (interfaces.ConsentStore, error)
}

func NewServer(jwtSecret string, origin string, otp otp.OTPAuthenticator, isdebug bool, users interfaces.UserStore, cfg config.Config, makeConsentStore func(bankId int64) (interfaces.ConsentStore, error)) *server {
	s := &server{
		jwtSecret:        []byte(jwtSecret),
		jwtMethod:        jwt.SigningMethodHS256,
		corsOrigin:       origin,
		apis:             map[int64]*hack.ApiClient{},
		otp:              otp,
		isdebug:          isdebug,
		users:            users,
		cfg:              cfg,
		makeConsentStore: makeConsentStore,
	}
	return s
}

// can panic on bad patterns
func (s server) SetHandlers(mux *http.ServeMux) {
	// for _, v := range proxyRoutes {
	muxx := http.NewServeMux()
	mux.Handle("/", s.corsMiddleware(muxx))
	muxx.HandleFunc("GET /banks", wrapHandler(func(w http.ResponseWriter, r *http.Request) (err error) {
		type bank struct {
			BankID          int64  `json:"bank_id"`
			BankName        string `json:"bank_name"`
			BankDescription string `json:"bank_description"`
			BankAvatar      string `json:"bank_avatar"`
		}
		type resp struct {
			Banks []bank `json:"banks"`
		}
		banks, err := s.cfg.Banks(r.Context())
		if err != nil {
			return err
		}
		return writeJsonBody(w,
			resp{
				Banks: lo.Map(banks, func(item config.BankConfig, _ int) bank {
					return bank{
						BankID:          item.ID(),
						BankName:        item.Name(),
						BankDescription: item.Description(),

						// TODO
						// BankAvatar:      item.Avatar(),
					}
				})})
	}),
	)
	muxx.HandleFunc("GET /my/banks", wrapHandler(s.withJwtUser(func(w http.ResponseWriter, r *http.Request, uid string) (err error) {
		type bank struct {
			BankID          int64  `json:"bank_id"`
			MyBankClientId  string `json:"my_bank_client_id"`
			BankName        string `json:"bank_name"`
			BankDescription string `json:"bank_description"`
			BankAvatar      string `json:"bank_avatar"`
		}
		type resp struct {
			Banks []bank `json:"banks"`
		}
		u, err := s.users.Find(r.Context(), uid)
		if err != nil {
			return err
		}
		banks, err := u.Banks(r.Context())
		if err != nil {
			return err
		}
		return writeJsonBody(w,
			resp{
				Banks: lo.Map(banks, func(item interfaces.UserBank, _ int) bank {
					return bank{
						BankID:          item.BankID,
						BankName:        item.BankName,
						BankAvatar:      item.BankAvatar,
						MyBankClientId:  item.MyBankClientId,
						BankDescription: item.BankDescription,
					}
				})})
	})),
	)
	muxx.HandleFunc("POST /my/banks", wrapHandler(s.withJwtUser(func(w http.ResponseWriter, r *http.Request, uid string) error {
		type input struct {
			ClientID string `json:"client_id"`
			BankID   int64  `json:"bank_id"`
		}
		var err error
		in, err := readJsonBody[input](w, r)
		if err != nil {
			return err
		}
		u, err := s.users.Find(r.Context(), uid)
		if err != nil {
			return err
		}
		err = u.AddBank(r.Context(), in.BankID, in.ClientID)
		if err != nil {
			return err
		}
		return s.apis[in.BankID].Client(in.ClientID).RequestConsents(r.Context())
	})),
	)
	muxx.HandleFunc("POST /auth/begin", wrapHandler(func(w http.ResponseWriter, r *http.Request) (err error) {
		type input struct {
			Email string `json:"email"`
		}
		in, err := readJsonBody[input](w, r)
		if err != nil {
			return err
		}
		ses := rand.Text()
		err = s.otp.InitCodeAuth(r.Context(), in.Email, ses)
		if err != nil {
			return err
		}

		return writeJsonBody(w, struct {
			SessionID string `json:"session_id"`
		}{
			SessionID: ses,
		})
	}),
	)
	muxx.HandleFunc("POST /auth/complete", wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		type input struct {
			SessionID string `json:"session_id"`
			Code      string `json:"code"`
		}
		in, err := readJsonBody[input](w, r)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", in)
		u, err := s.otp.CompleteCodeAuth(r.Context(), in.SessionID, in.Code)
		if err != nil {
			if errors.Is(err, otp.ErrOtpMismatch) {
				w.WriteHeader(http.StatusUnauthorized)
				return nil
			}
			return err
		}
		err = s.users.Create(r.Context(), u)
		if err != nil {
			fmt.Println(8725)
			return err
		}
		tok, err := s.newJWTString(u.Email)
		if err != nil {
			fmt.Println(98796)
			return err
		}
		output := struct {
			Token string `json:"token"`
		}{
			Token: tok,
		}
		return writeJsonBody(w, output)
	}),
	)
	muxx.HandleFunc("/",
		wrapHandler(
			s.withJwtUser(func(w http.ResponseWriter, r *http.Request, uid string) error {
				bankids := r.Header.Get("x-bank-id")
				bankid, err := strconv.ParseInt(bankids, 10, 64)
				if err != nil {
					return err
				}
				bank, ok := s.apis[bankid]
				if !ok {
					log.Printf("%+v\n", s.apis)
					return errors.New("19857 ")
				}
				// cli := bank.Client(uid)
				cons, err := bank.CS.FirstValidFor(r.Context(), uid, pe.ReadAccountsBasic)
				if err != nil {
					return err
				}
				newurl := bank.ApiUrl + "/" + r.URL.Path + "?" + r.URL.RawQuery
				newreq, err := http.NewRequestWithContext(r.Context(), r.Method, newurl, r.Body)
				if err != nil {
					return err
				}
				myBankId, err := s.cfg.BankID(r.Context())
				if err != nil {
					return err
				}
				newreq.Header.Set("authorization", "Bearer "+s.apis[bankid].AccessToken)
				newreq.Header.Set("x-requesting-bank", myBankId)
				newreq.Header.Set("x-consent-id", cons.ID) // secutity is someone else's problem now
				if s.isdebug {
					cmd, err := http2curl.GetCurlCommand(newreq)
					if err != nil {
						return err
					}
					log.Println(cmd)
				}
				resp, err := http.DefaultClient.Do(newreq)
				if err != nil {
					return err
				}
				for name, values := range resp.Header {
					for _, value := range values {
						w.Header().Add(name, value)
					}
				}
				_, err = io.Copy(w, resp.Body)
				if err != nil {
					return err
				}
				return nil
			})),
	)
	// }
	// return nil
}
