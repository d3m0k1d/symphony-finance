package interfaces

import "context"

type UserBank struct {
	Bank
	MyBankClientId string
}
type User interface {
	AddBank(ctx context.Context, bankId int64, clientId string) error
	Banks(ctx context.Context) ([]UserBank, error)
}
type UserIdentity struct {
	Email string
}
type UserStore interface {
	Find(ctx context.Context, uid string) (User, error)
	Create(ctx context.Context, u UserIdentity) error
}
