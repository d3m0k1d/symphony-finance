package maps

import (
	"context"
	"errors"
	"time"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/samber/lo"
)

// var _ interfaces.ConsentStore = Cache{}

var ErrConsentNotFound = errors.New("consent not found")

// FirstValidFor implements interfaces.ConsentStore.
func (c Cache) FirstValidFor(ctx context.Context, cid string, perm pe.PermissionsType) (cons interfaces.Consent, err error) {
	now := time.Now()
	cs := c.consents[cid][perm]
	co, ok := lo.Find(lo.Entries(cs), func(item lo.Entry[string, interfaces.Consent]) bool {
		return item.Value.ExpiresAt.Before(now)
	})
	if !ok {
		err = ErrConsentNotFound
		return
	}
	cons = interfaces.Consent{ID: co.Value.ID}
	return
}

// InsertConsent implements interfaces.ConsentStore.
func (c Cache) InsertConsent(ctx context.Context, cid string, perms []pe.PermissionsType, cons interfaces.Consent) error {
	for _, perm := range perms {
		if c.consents[cid] == nil {
			c.consents[cid] = make(map[pe.PermissionsType]map[string]interfaces.Consent, 1)
		}

		if c.consents[cid][perm] == nil {
			c.consents[cid][perm] = make(map[string]interfaces.Consent, len(perms))
		}
		c.consents[cid][perm][cons.ID] = cons
	}
	return nil
}
func (c Cache) DropExpiredConsents(ctx context.Context) error {
	now := time.Now()
	for _, v := range c.consents {
		for _, vv := range v {
			for kkk, vvv := range vv {
				if vvv.ExpiresAt.Before(now) {
					delete(vv, kkk)
				}
			}
		}
	}

	return nil
}
