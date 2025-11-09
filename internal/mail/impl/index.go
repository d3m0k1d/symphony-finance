package impl

import (
	"context"
	"fmt"
	"vtb-apihack-2025/internal/mail"

	gomail "github.com/wneessen/go-mail"
)

// import "net/mail"

type mailer struct {
	addr, username, password, sender string
}

// SendCode implements mail.Mailer.
func NewMailer(addr, username, password, sender string) (mailer, error) {
	mlr := mailer{addr, username, password, sender}
	return mlr, nil
}
func (m mailer) SendCode(ctx context.Context, rcpt string, code string) error {
	msg := gomail.NewMsg()
	if err := msg.From(m.sender); err != nil {
		return err
	}
	if err := msg.To(rcpt); err != nil {
		return err
	}
	msg.Subject("Подтвердите вход в Symphony")
	msg.SetBodyString(gomail.TypeTextPlain, fmt.Sprintf("Код подтверждения входа в symphony: %q\nНикому не говорите этот код!", code))
	cli, err := gomail.NewClient(m.addr,
		gomail.WithSMTPAuth(gomail.SMTPAuthAutoDiscover), gomail.WithTLSPortPolicy(gomail.TLSMandatory),
		gomail.WithUsername(m.username), gomail.WithPassword(m.password),
	)
	if err != nil {
		return err
	}
	if err := cli.DialAndSendWithContext(ctx, msg); err != nil {
		return err
	}
	return nil
}

var _ mail.Mailer = mailer{}
