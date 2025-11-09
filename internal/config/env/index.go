package env

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"vtb-apihack-2025/internal/config"

	"github.com/samber/lo"
)

var _ config.BankConfig = bankConf{}
var _ config.Config = conf{}

type conf struct {
	banks            []bankConf
	clientId         string
	clientSecret     string
	bankId, bankName string
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
	c      *conf
	apiurl string
}

// ApiUrl implements config.BankConfig.
func (b bankConf) ApiUrl(_ context.Context) (string, error) {
	return b.apiurl, nil
}

// CientId implements config.BankConfig.
func (b bankConf) CientId(_ context.Context) (string, error) {
	return b.c.clientId, nil
}

// CientSecret implements config.BankConfig.
func (b bankConf) CientSecret(_ context.Context) (string, error) {
	return b.c.clientSecret, nil
}

var ErrEnvNotFound = errors.New("Required environment variable not found")

// Banks implements config.Config.
func (c conf) Banks(_ context.Context) ([]config.BankConfig, error) {
	return lo.Map(c.banks, func(item bankConf, index int) config.BankConfig {
		return item
	}), nil
}
func EnvOrErr(k string) (v string, err error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrEnvNotFound, k)
	}
	return v, nil
}
func NewConfig() (config.Config, error) {
	banks, err1 := EnvOrErr("BANKS")
	clientId, err2 := EnvOrErr("CLIENT_ID")
	clientSecret, err3 := EnvOrErr("CLIENT_SECRET")
	bankId, err4 := EnvOrErr("MY_BANK_ID")
	bankName, err5 := EnvOrErr("MY_BANK_NAME")
	if err := errors.Join(err1, err2, err3, err4, err5); err != nil {
		return nil, err
	}
	banksS := strings.Split(banks, " ")
	c := conf{
		clientId:     clientId,
		clientSecret: clientSecret,
		bankId:       bankId,
		bankName:     bankName,
	}
	c.banks = lo.Map(banksS, func(item string, _ int) bankConf {
		return bankConf{
			apiurl: item,
			c:      &c,
		}
	})
	return c, nil
}
