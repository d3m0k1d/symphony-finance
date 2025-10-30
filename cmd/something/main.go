package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	auth "vtb-apihack-2025/client-pilot/auth/hack"
	"vtb-apihack-2025/client-pilot/le"
	"vtb-apihack-2025/client-pilot/payments"
	"vtb-apihack-2025/client-pilot/pe"

	"moul.io/http2curl"
)

func main() {
	if err := func() error {
		var err error
		ctx := context.Background()
		apiurl := os.Getenv("BANK")
		pec, err1 := pe.NewClientWithResponses(apiurl)
		lec, err2 := le.NewClientWithResponses(apiurl)
		_ = lec
		pc, err3 := payments.NewClientWithResponses(apiurl)
		ac, err4 := auth.NewClientWithResponses(apiurl)
		err = errors.Join(err1, err2, err3, err4)
		if err != nil {
			return err
		}
		clientId := os.Getenv("CLIENT_ID")
		client_secret := os.Getenv("CLIENT_SECRET")
		aresp, err := ac.CreateBankTokenAuthBankTokenPostWithResponse(ctx, &auth.CreateBankTokenAuthBankTokenPostParams{
			ClientId:     clientId,
			ClientSecret: client_secret,
		})
		if err != nil {
			return err
		}
		if false {
			curl, err := http2curl.GetCurlCommand(aresp.HTTPResponse.Request)
			if err != nil {
				return err
			}
			log.Println(curl)
		}

		if aresp.StatusCode() != 200 {
			return fmt.Errorf("Auth error: exit code: %d\nbody:\n%s", aresp.StatusCode(), aresp.Body)
		}
		authF := func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", "Bearer "+aresp.JSON200.AccessToken)
			return nil
		}
		bank_id := os.Getenv("MY_BANK_ID")
		{
			b := auth.ConsentRequestBody{
				ClientId:           "team074-5",
				Permissions:        []string{"ReadAccountsDetail", "ReadBalances", "ReadTransactionsDetail"},
				Reason:             sp("WE NEEd YOUR dATA!!!!!"),
				RequestingBank:     &bank_id,
				RequestingBankName: sp(os.Getenv("MY_BANK_NAME")),
			}
			resp, err := ac.RequestConsentAccountConsentsRequestPost(
				ctx,
				&auth.RequestConsentAccountConsentsRequestPostParams{
					XRequestingBank: &bank_id,
				},
				b,
				authF,
			)
			time.Now().Local().MarshalJSON()
			curl, _ := http2curl.GetCurlCommand(resp.Request)
			log.Println(curl)
			bj, _ := json.Marshal(b)
			log.Printf("%s\n", bj)
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
				return fmt.Errorf("Auth2 error: exit code: %d\nbody:\n%s", resp.StatusCode, resp.Body)
			}
			respp, err := auth.ParseRequestConsentAccountConsentsRequestPostResponse(resp)
			if err != nil {
				return err
			}
			fmt.Printf("consent request response: %v\n", respp.JSON200)
		}
		// c.LoginAuthLoginPost(c, client.LoginRequest{}, )
		if false {
			respp, err := pec.GetAccountsaccountIdTransactionsWithResponse(ctx, "14", &pe.GetAccountsaccountIdTransactionsParams{}, authF)
			if err != nil {
				return err
			}
			// io.Copy(os.Stdout, resp.Body)
			for _, t := range *respp.JSON200.Data.Transaction {
				fmt.Println(t.Amount)
			}
		}
		if false {
			respp, err := pc.CreatePaymentWithResponse(ctx, &payments.CreatePaymentParams{}, payments.PaymentRequest{}, authF)
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", respp.JSON201)
		}
		// resp,err:=
		// client.GetTransactionsAccountsAccountIdTransactionsGetResponse

		return nil
	}(); err != nil {
		log.Println(err)
	}
}
func sp(s string) *string {
	return &s
}
