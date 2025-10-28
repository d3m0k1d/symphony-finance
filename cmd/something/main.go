package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"vtb-apihack-2025/client-pilot/pe"
)

func main() {
	auth := func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJkZW1vLTAwMyIsInR5cGUiOiJjbGllbnQiLCJiYW5rIjoic2VsZiIsImV4cCI6MTc2MTcwNTE3MH0.Lb8WdPeoOSruor_77q-XPlkfVKI3-drqZCx5lTqPTOc")
		return nil
	}
	if err := func() error {
		ctx := context.Background()
		c, err := pe.NewClient(os.Getenv("BANK"))
		if err != nil {
			return err
		}
		// c.LoginAuthLoginPost(c, client.LoginRequest{}, )
		resp, err := c.GetAccountsaccountIdTransactions(ctx, "14", &pe.GetAccountsaccountIdTransactionsParams{}, auth)
		if err != nil {
			return err
		}
		respp, err := pe.ParseGetAccountsaccountIdTransactionsResponse(resp)
		if err != nil {
			return err
		}
		// io.Copy(os.Stdout, resp.Body)
		for _, t := range *respp.JSON200.Data.Transaction {
			fmt.Println(t.Amount)
		}
		// client.GetTransactionsAccountsAccountIdTransactionsGetResponse

		return nil
	}(); err != nil {
		log.Println(err)
	}
}
