package bobsql

import (
	"context"
	"vtb-apihack-2025/bobgen/models"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/aarondl/opt/omit"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	_ "modernc.org/sqlite"
)

func NewSqliteUserStore(dsn string) (*userStore, error) {
	db, err := bob.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return &userStore{exe: db}, nil
}

var _ interfaces.User = &user{}

type user struct {
	m *models.User
	s *userStore
}

// Banks implements interfaces.User.
func (u *user) Banks(ctx context.Context) ([]interfaces.UserBank, error) {
	banks, err := u.m.UserEmailUserBanks(models.Preload.UserBank.Bank()).All(ctx, u.s.exe)
	if err != nil {
		return nil, err
	}
	return lo.Map(banks, func(m *models.UserBank, _ int) interfaces.UserBank {
		return interfaces.UserBank{
			Bank:           MapBank(m.R.Bank),
			MyBankClientId: m.ClientID,
		}
	}), nil
}

// AddBank implements interfaces.User.
func (u user) AddBank(ctx context.Context, bankId int64, clientId string) error {

	return u.m.InsertUserEmailUserBanks(ctx, u.s.exe, &models.UserBankSetter{
		BankID: omit.From(bankId), ClientID: omit.From(clientId),
		// UserEmail: u.m.UserEmail,
	})
}

var _ interfaces.UserStore = &userStore{}

type userStore struct {
	exe bob.Executor
}

// Create implements interfaces.UserStore.
func (s *userStore) Create(ctx context.Context, u interfaces.UserIdentity) error {
	_, err := models.Users.Insert(&models.UserSetter{
		UserEmail: omit.From(u.Email),
	}, im.OnConflict(models.Users.Columns.UserEmail).DoNothing()).Exec(ctx, s.exe)
	if err != nil {
		return err
	}
	return nil
}

// Find implements interfaces.UserStore.
func (s *userStore) Find(ctx context.Context, uid string) (interfaces.User, error) {
	usr, err := models.FindUser(ctx, s.exe, uid)
	if err != nil {
		return nil, err
	}
	return &user{usr, s}, nil
}
