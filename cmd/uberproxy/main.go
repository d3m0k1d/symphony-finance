package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"vtb-apihack-2025/internal/api/uberproxy"
	bobc "vtb-apihack-2025/internal/config/bobsql"
	"vtb-apihack-2025/internal/mail"
	"vtb-apihack-2025/internal/mail/fake"
	email "vtb-apihack-2025/internal/mail/impl"
	otp "vtb-apihack-2025/internal/otp/impl"
	"vtb-apihack-2025/internal/storage/impl/bobsql"
	"vtb-apihack-2025/internal/storage/impl/redis"
	"vtb-apihack-2025/internal/storage/interfaces"
)

const debug = true

func main() {
	if err := func() error {
		var err error
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dsn := os.Getenv("SQLITE_DSN")
		cfg, err := bobc.NewConfig(dsn)
		if err != nil {
			return err
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
		users, err := bobsql.NewSqliteUserStore(dsn)
		if err != nil {
			return err
		}
		otper := otp.NewOtper(mailer, otpstore)
		// TODO: env config
		sv := uberproxy.NewServer(
			"i needa get laid",
			os.Getenv("CORS_ORIGIN"), otper, debug, users, cfg,
			func() (interfaces.ConsentStore, error) {
				return bobsql.NewSqliteConsentStore(dsn)
			},
		)
		err = sv.RefreshBanks(ctx)
		if err != nil {
			return err
		}
		sv.SetHandlers(http.DefaultServeMux)
		return http.ListenAndServe(":8089", nil)
	}(); err != nil {
		log.Println(err)
	}
}
