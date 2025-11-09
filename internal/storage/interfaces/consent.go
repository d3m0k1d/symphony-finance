package interfaces

import (
	"context"
	"time"
	"vtb-apihack-2025/client-pilot/pe"
)

type Consent struct {
	ID        string
	ExpiresAt time.Time
}
type ConsentStore interface {
	FirstValidFor(ctx context.Context, clientId string, perm pe.PermissionsType) (Consent, error)
	InsertConsent(ctx context.Context, clientId string, perms []pe.PermissionsType, cons Consent) error
	DropExpiredConsents(ctx context.Context) error
}
