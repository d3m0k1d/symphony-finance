package fake

import (
	"context"
	"log"
	"vtb-apihack-2025/internal/mail"
)

// import "net/mail"

type mailer struct {
}

// SendCode implements mail.Mailer.
func NewMailer() (mailer, error) {
	mlr := mailer{}
	return mlr, nil
}
func (m mailer) SendCode(ctx context.Context, rcpt string, code string) error {
	log.Printf("Confirmation code for %q: %q\n", rcpt, code)
	return nil
}

var _ mail.Mailer = mailer{}
