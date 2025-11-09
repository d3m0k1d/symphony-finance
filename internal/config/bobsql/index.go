package env

import (
	"context"
	"errors"
	"fmt"
	"os"
	"vtb-apihack-2025/bobgen/models"
	"vtb-apihack-2025/internal/config"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	_ "modernc.org/sqlite"
)

var _ config.BankConfig = bankConf{}
var _ config.Config = &conf{}

type conf struct {
	banks            []bankConf
	clientId         string
	clientSecret     string
	bankId, bankName string
	exe              bob.Executor
}

// BankID implements config.Config.
func (c conf) BankID(ctx context.Context) (string, error) {
	return c.bankId, nil
}

// BankName implements config.Config.
func (c conf) BankName(ctx context.Context) (string, error) {
	return c.bankName, nil
}

type bankConf struct {
	c *conf
	m *models.Bank
}

// Description implements config.BankConfig.
func (b bankConf) Description() string {
	return b.m.BankDescription.GetOrZero()
}

// Name implements config.BankConfig.
func (b bankConf) Name() string {
	return b.m.BankName.GetOrZero()
}

// ApiUrl implements config.BankConfig.
func (b bankConf) ApiUrl(_ context.Context) (url string, err error) {
	return b.m.BankAPIURL, nil
}

// CientId implements config.BankConfig.
func (b bankConf) CientId(_ context.Context) (string, error) {
	return b.c.clientId, nil
}

// CientSecret implements config.BankConfig.
func (b bankConf) CientSecret(_ context.Context) (string, error) {
	return b.c.clientSecret, nil
}
func (b bankConf) ID() int64 {
	return b.m.BankID
}

var ErrEnvNotFound = errors.New("Required environment variable not found")

// Banks implements config.Config.
func (c *conf) Banks(ctx context.Context) ([]config.BankConfig, error) {
	banks, err := models.Banks.Query().All(ctx, c.exe)
	if err != nil {
		return nil, err
	}
	return lo.Map(banks, func(item *models.Bank, _ int) config.BankConfig {
		return bankConf{c, item}
	}), nil
}
func EnvOrErr(k string) (v string, err error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrEnvNotFound, k)
	}
	return v, nil
}
func NewConfig(dsn string) (*conf, error) {
	clientId, err2 := EnvOrErr("CLIENT_ID")
	clientSecret, err3 := EnvOrErr("CLIENT_SECRET")
	bankId, err4 := EnvOrErr("MY_BANK_ID")
	bankName, err5 := EnvOrErr("MY_BANK_NAME")
	if err := errors.Join(err2, err3, err4, err5); err != nil {
		return nil, err
	}
	db, err := bob.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	c := conf{
		clientId:     clientId,
		clientSecret: clientSecret,
		bankId:       bankId,
		bankName:     bankName,
		exe:          db,
	}
	return &c, nil
}
