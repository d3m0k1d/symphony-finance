package hack

import (
	"errors"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
	"vtb-apihack-2025/client-pilot/le"
	"vtb-apihack-2025/client-pilot/payments"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/lib/client"
)

type Client struct {
	ac                       auth.ClientInterface
	lec                      le.ClientInterface
	pec                      pe.ClientInterface
	pc                       payments.ClientInterface
	authF                    auth.RequestEditorFn
	bankId, bankName, apiUrl string
}

func NewClient(apiUrl, bankId, bankName string) (client.Client, error) {
	var err error
	pec, err1 := pe.NewClientWithResponses(apiUrl)
	lec, err2 := le.NewClientWithResponses(apiUrl)
	_ = lec
	pc, err3 := payments.NewClientWithResponses(apiUrl)
	ac, err4 := auth.NewClientWithResponses(apiUrl)
	err = errors.Join(err1, err2, err3, err4)
	if err != nil {
		return nil, err
	}
	c := Client{
		ac:       ac,
		lec:      lec,
		pec:      pec,
		pc:       pc,
		bankId:   bankId,
		bankName: bankName,
		apiUrl:   apiUrl,
	}
	return c, nil
}
