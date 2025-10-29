package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"vtb-apihack-2025/client-pilot/le"
	"vtb-apihack-2025/client-pilot/payments"
	"vtb-apihack-2025/client-pilot/pe"
)

func main() {
	auth := func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkZW1vLTAwMyIsInR5cGUiOiJjbGllbnQiLCJiYW5rIjoic2VsZiIsImV4cCI6MTc2MTcwNTE3MH0.Lb8WdPeoOSruor_77q-XPlkfVKI3-drqZCx5lTqPTOc")
		return nil
	}
	if err := func() error {
		var err error
		ctx := context.Background()
		apiurl := os.Getenv("BANK")
		pec, err1 := pe.NewClientWithResponses(apiurl)
		lec, err2 := le.NewClientWithResponses(apiurl)
		pc, err3 := payments.NewClientWithResponses(apiurl)
		err = errors.Join(err1, err2, err3)
		if err != nil {
			return err
		}
		_ = lec
		// c.LoginAuthLoginPost(c, client.LoginRequest{}, )
		{
			respp, err:=pec.GetAccountsaccountIdTransactionsWithResponse(ctx, "14", &pe.GetAccountsaccountIdTransactionsParams{}, auth) 
			if err != nil {
				return err
			}
			// io.Copy(os.Stdout, resp.Body)
			for _, t := range *respp.JSON200.Data.Transaction {
				fmt.Println(t.Amount)
			}
		}
		{
			respp, err := pc.CreatePaymentWithResponse(ctx, &payments.CreatePaymentParams{}, payments.PaymentRequest{}, auth)
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
