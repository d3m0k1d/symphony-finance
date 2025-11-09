package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"vtb-apihack-2025/internal/client"
	"vtb-apihack-2025/internal/client/hack"
	envc "vtb-apihack-2025/internal/config/env"
	maps "vtb-apihack-2025/internal/storage/impl/map"
)

const cid = "team074-5"

func main() {
	if err := func() error {
		var err error
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cfg, err := envc.NewConfig()
		if err != nil {
			return err
		}

		bankId, err1 := cfg.BankID(ctx)
		bankName, err2 := cfg.BankName(ctx)
		banks, err3 := cfg.Banks(ctx)
		if err := errors.Join(err1, err2, err3); err != nil {
			return err
		}
		apiclients := make([]client.APIClient, len(banks))
		var bigerr error
		for i, bank := range banks {
			apiUrl, err1 := bank.ApiUrl(ctx)
			clientId, err2 := bank.CientId(ctx)
			clientSecret, err3 := bank.CientSecret(ctx)
			if err := errors.Join(err1, err2, err3); err != nil {
				bigerr = errors.Join(bigerr, fmt.Errorf("error getting config for bank %q: %w", apiUrl, err))
				continue
			}
			u, err := url.Parse(apiUrl)
			if err != nil {
				return err
			}
			cs := maps.NewCache()
			bc, err := hack.NewClient(apiUrl, bankId, bankName, u.Hostname(), cs, true)
			if err != nil {
				bigerr = errors.Join(bigerr, fmt.Errorf("error initializing bank %q: %w", apiUrl, err))
				continue
			}
			if err := bc.Authenticate(ctx, clientId, clientSecret); err != nil {
				bigerr = errors.Join(bigerr, fmt.Errorf("error authenticating bank %q: %w", apiUrl, err))
				continue
			}
			apiclients[i] = bc
		}
		if bigerr != nil {
			return bigerr
		}
		cli := apiclients[0].Client(cid)
		err = cli.RequestConsents(ctx)
		if err != nil {
			return err
		}
		accs, err := cli.Accounts(ctx)
		if err != nil {
			return err
		}
		if len(accs) == 0 {
			return errors.New("zero accounts")
		}
		ts, err := cli.TransactionsPage(ctx, accs[0].ID, 1, nil, nil)
		if err != nil {
			return err
		}
		for _, t := range ts {
			fmt.Printf("%#v\n", t)
		}

		// c.LoginAuthLoginPost(c, client.LoginRequest{}, )
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
