package mail

import "context"

type Mailer interface {
	SendCode(ctx context.Context, email string, code string) error
}
