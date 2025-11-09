package bobsql

import (
	"context"
	"time"
	"vtb-apihack-2025/bobgen/models"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/types"
	_ "modernc.org/sqlite"
)

func NewSqliteConsentStore(dsn string) (*consentStore, error) {
	db, err := bob.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return &consentStore{exe: db}, nil
}

var _ interfaces.ConsentStore = &consentStore{}

type consentStore struct {
	exe    bob.Executor
	bankId int64
}

// DropExpiredConsents implements interfaces.ConsentStore.
func (c *consentStore) DropExpiredConsents(ctx context.Context) (err error) {
	_, err = models.Consents.Delete(models.DeleteWhere.Consents.ConsentExpiry.LTE(types.Time{time.Now()})).Exec(ctx, c.exe)
	return
}

// FirstValidFor implements interfaces.ConsentStore.
func (c *consentStore) FirstValidFor(ctx context.Context, clientId string, perm pe.PermissionsType) (v interfaces.Consent, err error) {
	cons, err := models.Consents.Query(models.SelectWhere.Consents.ConsentExpiry.GT(types.Time{time.Now()})).One(ctx, c.exe)
	if err != nil {
		return
	}
	v = interfaces.Consent{
		ID:        cons.ConsentID,
		ExpiresAt: cons.ConsentExpiry.Time,
	}
	return
}

// InsertConsent implements interfaces.ConsentStore.
func (c *consentStore) InsertConsent(ctx context.Context, clientId string, perms []pe.PermissionsType, cons interfaces.Consent) (err error) {
	_, err = models.Consents.Insert(&models.ConsentSetter{
		ConsentID:     omit.From(cons.ID),
		ConsentExpiry: omit.From(types.Time{cons.ExpiresAt}),
	}).Exec(ctx, c.exe)
	if err != nil {
		return err
	}
	_, err = models.UserBankConsents.Insert(&models.UserBankConsentSetter{
		ConsentID: omit.From(cons.ID),
		BankID:    omit.From(c.bankId),
		UserEmail: omit.From(cons.ID),
	}).Exec(ctx, c.exe)
	if err != nil {
		return err
	}
	return
}
