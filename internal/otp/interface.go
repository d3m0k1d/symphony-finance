package otp

import (
	"context"
	"errors"
)

var ErrOtpMismatch = errors.New("otp code is wrong")

type UserIdentity struct {
	Email string
}
type OTPAuthenticator interface {
	InitCodeAuth(ctx context.Context, email, session string) error
	CompleteCodeAuth(ctx context.Context, session string, code string) (user UserIdentity, err error)
}
