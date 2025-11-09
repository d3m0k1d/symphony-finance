package bobsql

import (
	"context"
	"vtb-apihack-2025/bobgen/models"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	_ "modernc.org/sqlite"
)

func NewSqliteBankStore(dsn string) (*bankStore, error) {
	db, err := bob.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return &bankStore{exe: db}, nil
}

var _ interfaces.BankStore = &bankStore{}

type bankStore struct {
	exe bob.Executor
}

func MapBank(in *models.Bank) interfaces.Bank {
	return interfaces.Bank{
		BankID:          in.BankID,
		BankName:        in.BankName.GetOrZero(),
		BankDescription: in.BankName.GetOrZero(),
		BankAvatar:      in.BankName.GetOrZero(),
	}
}

// Find implements interfaces.BankStore.
func (b *bankStore) Find(ctx context.Context, id int64) (interfaces.Bank, error) {
	bank, err := models.FindBank(ctx, b.exe, id)
	if err != nil {
		return interfaces.Bank{}, err
	}
	return MapBank(bank), err
}

func (b *bankStore) All(ctx context.Context) ([]interfaces.Bank, error) {
	banks, err := models.Banks.Query().All(ctx, b.exe)
	if err != nil {
		return nil, err
	}
	return lo.Map(banks, func(item *models.Bank, _ int) interfaces.Bank {
		return MapBank(item)
	}), nil
}
