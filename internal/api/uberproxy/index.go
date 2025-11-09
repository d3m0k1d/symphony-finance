package uberproxy

import (
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/client/hack"
	"vtb-apihack-2025/internal/otp"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"moul.io/http2curl"
)

// var proxyRoutes = []string{}

type server struct {
	jwtSecret []byte
	jwtMethod jwt.SigningMethod
	apis      map[string]*hack.ApiClient

	corsOrigin string
	otp        otp.OTPAuthenticator

	isdebug bool
}

func NewServer(jwtSecret string, apis map[string]*hack.ApiClient, origin string, otp otp.OTPAuthenticator, isdebug bool) server {
	return server{
		jwtSecret:  []byte(jwtSecret),
		jwtMethod:  jwt.SigningMethodHS256,
		corsOrigin: origin,

		apis:    apis,
		otp:     otp,
		isdebug: isdebug,
	}
}

// can panic on bad patterns
func (s server) SetHandlers(mux *http.ServeMux) {
	// for _, v := range proxyRoutes {
	muxx := http.NewServeMux()
	mux.Handle("/", s.corsMiddleware(muxx))
	muxx.HandleFunc("GET /banks", wrapHandler(func(w http.ResponseWriter, r *http.Request) (err error) {
		type bank struct {
			BankID          string `json:"bank_id"`
			BankName        string `json:"bank_name"`
			BankDescription string `json:"bank_description"`
			BankAvatar      string `json:"bank_avatar"`
		}
		type resp struct {
			Banks []bank `json:"banks"`
		}
		return writeJsonBody(w,
			resp{
				Banks: lo.Map(lo.Values(s.apis), func(item *hack.ApiClient, _ int) bank {
					return bank{
						BankID: item.ProviderBankID(),

						// TODO
						BankName:        "",
						BankDescription: "",
						BankAvatar:      "",
					}
				})})
	}),
	)
	muxx.HandleFunc("POST /my/banks", wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		type input struct {
			ClientID string `json:"client_id"`
			BankID   string `json:"bank_id"`
		}
		var err error
		in, err := readJsonBody[input](w, r)
		if err != nil {
			return err
		}

		return s.apis[in.BankID].Client(in.ClientID).RequestConsents(r.Context())
	}),
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
		return nil
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
		u, err := s.otp.CompleteCodeAuth(r.Context(), in.SessionID, in.Code)
		if err != nil {
			if errors.Is(err, otp.ErrOtpMismatch) {
				w.WriteHeader(http.StatusUnauthorized)
				return nil
			}
			return err
		}
		tok, err := s.newJWTString(u.Email)
		if err != nil {
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
	muxx.HandleFunc("/", wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		toks, found := strings.CutPrefix(r.Header.Get("authorization"), "Bearer ")
		if !found {
			return errors.New("authentication header is not a token bearer")
		}
		bankid := r.Header.Get("x-bank-id")
		bank, ok := s.apis[bankid]
		if !ok {
			log.Printf("%+v\n", s.apis)
			return errors.New("fuck it 19857")
		}
		uid, err := s.parseJwt(toks)
		if err != nil {
			return err
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
		newreq.Header.Set("authorization", "Bearer "+s.apis[bankid].AccessToken)
		newreq.Header.Set("x-requesting-bank", bank.ProviderBankID())
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
	}),
	)
	// }
	// return nil
}
