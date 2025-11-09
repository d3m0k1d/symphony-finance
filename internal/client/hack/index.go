package hack

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
	"vtb-apihack-2025/client-pilot/le"
	"vtb-apihack-2025/client-pilot/payments"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/storage/interfaces"

	"moul.io/http2curl"
)

type ApiClient struct {
	// TODO: hide
	AccessToken string

	// TODO: use client mods to apply it to every request automagically.
	//       for authentication use plain http client and New*Request. also maybe do authn in ctor
	//       and make authn private. use pointer in authF to update tokens on expiry
	ac  auth.ClientInterface
	lec le.ClientInterface
	pec pe.ClientInterface
	pc  payments.ClientInterface

	// TODO: hide
	BankId,

	bankName, ApiUrl, consentReason string

	providerBankID string
	// TODO: hide
	CS             interfaces.ConsentStore
	authF          func(ctx context.Context, req *http.Request) error
	authTimer      *time.Timer
}

// ProviderBankID implements client.Client.
func (c *ApiClient) ProviderBankID() string {
	log.Println(c.providerBankID) 
	return c.providerBankID
}

func NewClient(apiUrl, bankId, bankName, providerBankID string, cs interfaces.ConsentStore, debug bool) (*ApiClient, error) {
	var err error
	debugF := func(ctx context.Context, req *http.Request) error {
		if debug {
			cmd, err := http2curl.GetCurlCommand(req)
			if err != nil {
				log.Println(err)
				return nil
			}
			log.Println(cmd)
		}
		return nil
	}
	pec, err1 := pe.NewClientWithResponses(apiUrl, pe.WithRequestEditorFn(debugF))
	lec, err2 := le.NewClientWithResponses(apiUrl, le.WithRequestEditorFn(debugF))
	pc, err3 := payments.NewClientWithResponses(apiUrl, payments.WithRequestEditorFn(debugF))
	ac, err4 := auth.NewClientWithResponses(apiUrl, auth.WithRequestEditorFn(debugF))
	err = errors.Join(err1, err2, err3, err4)
	if err != nil {
		return nil, err
	}
	c := ApiClient{
		ac:             ac,
		lec:            lec,
		pec:            pec,
		pc:             pc,
		BankId:         bankId,
		bankName:       bankName,
		ApiUrl:         apiUrl,
		authF:          nil,
		providerBankID: providerBankID,
		CS:             cs,
	}
	return &c, nil
}
