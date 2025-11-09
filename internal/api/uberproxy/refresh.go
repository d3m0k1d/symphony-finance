package uberproxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"vtb-apihack-2025/internal/client/hack"
	"vtb-apihack-2025/internal/config"

	"github.com/samber/lo"
)

// TODO: refresh using the existing data fetched on behalf of user in /banks instead of dumb polling
func (s server) RefreshBanks(ctx context.Context) error {
	banks, err := s.cfg.Banks(ctx)
	if err != nil {
		return err
	}
	bankId, err := s.cfg.BankID(ctx)
	if err != nil {
		return err
	}
	bankName, err := s.cfg.BankName(ctx)
	newmap := lo.FromEntries(lo.Map(banks, func(item config.BankConfig, _ int) lo.Entry[int64, config.BankConfig] {
		return lo.Entry[int64, config.BankConfig]{item.ID(), item}
	}))
	newkeys, _ := lo.Difference(lo.Keys(newmap), lo.Keys(s.apis))
	var bigerr error
	for _, k := range newkeys {
		bank := newmap[k]
		apiUrl, err1 := bank.ApiUrl(ctx)
		clientId, err2 := bank.CientId(ctx)
		clientSecret, err3 := bank.CientSecret(ctx)
		if err := errors.Join(err1, err2, err3); err != nil {
			bigerr = errors.Join(bigerr, fmt.Errorf("error getting config for bank %q: %w", apiUrl, err))
			continue
		}
		cs, err := s.makeConsentStore(bank.ID())
		if err != nil {
			return err
		}
		bc, err := hack.NewClient(apiUrl, bankId, bankName, bank.ID(), cs, true)
		if err != nil {
			bigerr = errors.Join(bigerr, fmt.Errorf("error initializing bank %q: %w", apiUrl, err))
			continue
		}
		if err := bc.Authenticate(ctx, clientId, clientSecret); err != nil {
			bigerr = errors.Join(bigerr, fmt.Errorf("error authenticating bank %q: %w", apiUrl, err))
			continue
		}
		s.apis[bc.ProviderBankID()] = bc
	}
	return nil
}
func (s server) PollBanks(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := s.RefreshBanks(ctx); err != nil {
			log.Println(err)
		}
	}
}
