package otp

import (
	"context"
	"errors"
	"vtb-apihack-2025/internal/storage/interfaces"
)

var ErrOtpMismatch = errors.New("otp code is wrong")

type OTPAuthenticator interface {
	InitCodeAuth(ctx context.Context, email, session string) error
	CompleteCodeAuth(ctx context.Context, session string, code string) (user interfaces.UserIdentity, err error)
}
