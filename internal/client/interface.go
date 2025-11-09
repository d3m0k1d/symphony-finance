package client

import (
	"context"
	"time"
	"vtb-apihack-2025/client-pilot/pe"
)

//	type RequestedConsent struct {
//		ExpiresAt time.Time
//		ClientId  string
//		ConsentID string
//	}
type APIClient interface {
	Authenticate(ctx context.Context, clientId, clientSecret string) error
	ProviderBankID() string
	Client(clientId string) Client
}
type Account struct {
	ID string
}

type Client interface {
	TransactionsPage(ctx context.Context, accId string, page int32, fromBookingDateTime *time.Time, toBookingDateTime *time.Time) ([]pe.TransactionHistory, error)
	RequestConsents(ctx context.Context) error
	Accounts(ctx context.Context) ([]Account, error)
}
