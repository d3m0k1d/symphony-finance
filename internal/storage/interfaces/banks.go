package interfaces

import "context"

type Bank struct {
	BankID          int64
	BankName        string
	BankDescription string
	BankAvatar      string
}
type BankStore interface {
	Find(ctx context.Context, id int64) (Bank, error)
	All(ctx context.Context) ([]Bank, error)
}
