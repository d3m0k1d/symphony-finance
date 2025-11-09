package interfaces

import "context"

type OTPSession struct {
	Hash  []byte
	Email string
}
type OtpSessionStore interface {
	StoreCode(ctx context.Context, session, email string, hash []byte) error
	RetrieveCode(ctx context.Context, session string) (result OTPSession, err error)
	// TODO: this is unused
	// not sure whether code guessing accounting should be done inside the store
	// or outside
	DropCode(ctx context.Context, session string) (err error)
}
