package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"vtb-apihack-2025/internal/api/uberproxy"
	"vtb-apihack-2025/internal/client/hack"
	envc "vtb-apihack-2025/internal/config/env"
	"vtb-apihack-2025/internal/mail"
	"vtb-apihack-2025/internal/mail/fake"
	email "vtb-apihack-2025/internal/mail/impl"
	otp "vtb-apihack-2025/internal/otp/impl"
	maps "vtb-apihack-2025/internal/storage/impl/map"
	"vtb-apihack-2025/internal/storage/impl/redis"

	"github.com/samber/lo"
)

const debug = true

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
		apiclients := make([]*hack.ApiClient, len(banks))
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
		var mailer mail.Mailer
		if debug {
			mailer, err = fake.NewMailer()
		} else {
			mailer, err = email.NewMailer(os.Getenv("SMTP_ADDR"), os.Getenv("SMTP_LOGIN"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SENDER_MAIL"))
		}
		if err != nil {
			return err
		}
		otpstore, err := redis.NewOtpCache(os.Getenv("REDIS_ADDR"))
		if err != nil {
			return err
		}
		otper := otp.NewOtper(mailer, otpstore)
		uberproxy.NewServer(
			"i needa get laid",
			lo.FromEntries(lo.Map(apiclients, func(item *hack.ApiClient, _ int) lo.Entry[string, *hack.ApiClient] {
				return lo.Entry[string, *hack.ApiClient]{item.ProviderBankID(), item}
			})),
			os.Getenv("CORS_ORIGIN"), otper, debug,
		).SetHandlers(http.DefaultServeMux) // TODO: env config
		return http.ListenAndServe(":8089", nil)
	}(); err != nil {
		log.Println(err)
	}
}
