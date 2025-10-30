package client

import (
	"context"
	"time"
)

type RequestedConsent struct {
	ExpiresAt time.Time
	ClientId  string
}
type Client interface {
	RequestConsent(ctx context.Context, clientId string) (RequestedConsent, error)
	Authenticate(ctx context.Context, clientId, clientSecret string) error
}
