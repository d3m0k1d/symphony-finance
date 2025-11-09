package config

import "context"

type BankConfig interface {
	ApiUrl(ctx context.Context) (string, error)
	CientId(ctx context.Context) (string, error)
	CientSecret(ctx context.Context) (string, error)
}
type Config interface {
	Banks(ctx context.Context) ([]BankConfig, error)
	BankID(ctx context.Context) (string, error)
	BankName(ctx context.Context) (string, error)
}
