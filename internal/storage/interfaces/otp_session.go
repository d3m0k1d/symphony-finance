package interfaces

import "context"

type OTPSession struct {
	Hash  []byte
	Email string
}
type OtpSessionStore interface {
	StoreCode(ctx context.Context, session, email string, hash []byte) error
	RetrieveCode(ctx context.Context, session string) (result OTPSession, err error)
}
