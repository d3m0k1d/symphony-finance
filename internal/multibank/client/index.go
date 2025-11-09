package client

import (
	"context"
	"database/sql"
	"time"
	"vtb-apihack-2025/client-pilot/pe"
	"vtb-apihack-2025/internal/client"

	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type MultiBankClient struct {
	BankAPIClients []client.APIClient
}

type TransactionHistoryWithBankName struct {
	pe.TransactionHistory
	BankName string
}

func (mbc MultiBankClient) GetNCombinedTransactions(ctx context.Context, accId string, banknames []string, beginTime time.Time, n int) ([]TransactionHistoryWithBankName, error) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(cctx)
	ch := mbc.CombinedTransactions(ctx, accId, banknames, beginTime, eg)
	out := make([]TransactionHistoryWithBankName, 0, n)
	for v := range ch {
		out = append(out, v)
	}
	err := eg.Wait()
	return out, err

}

type Ch[T any] struct {
	closed bool
	ch     chan T
}

func (mbc MultiBankClient) CombinedTransactions(ctx context.Context, accId string, banknames []string, beginTime time.Time, eg *errgroup.Group) chan TransactionHistoryWithBankName {
	clients := lo.Filter(mbc.BankAPIClients, func(item client.APIClient, index int) bool {
		for _, b := range banknames {
			if item.ProviderBankID() == b {
				return true
			}
		}
		return false
	})
	chs := make(map[string]*Ch[pe.TransactionHistory], len(clients))
	for _, c := range clients {
		chs[c.ProviderBankID()] = &Ch[pe.TransactionHistory]{false, Channel(ctx, accId, c, eg, beginTime)}
	}
	// TODO: merge into single map
	trs := make(map[string][]pe.TransactionHistory, len(chs))
	lastTime := beginTime
	out := make(chan TransactionHistoryWithBankName)
	eg.Go(func() error {
		<-ctx.Done()
		close(out)
		return nil
	})
	eg.Go(func() error {
		for {
			for k, ch := range chs {
			RETRY:
				select {
				case tr, ok := <-ch.ch:
					if !ok {
						ch.closed = true
						continue
					}
					if tr.BookingDateTime.After(lastTime) {
						trs[k] = append(trs[k], tr)
					} else {
						goto RETRY
					}
				case <-ctx.Done():
					return nil
				}
			}
			for k, ch := range chs {
				if ch.closed && len(trs[k]) == 0 {
					delete(chs, k)
					delete(trs, k)
				}
			}
			if len(chs) == 0 {
				return sql.ErrNoRows
			}
			_, emptyFound := lo.Find(lo.Entries(trs), func(item lo.Entry[string, []pe.TransactionHistory]) bool {
				return len(item.Value) == 0
			})
			if emptyFound {
				continue
			}
			for {
				first := lo.MaxBy(
					lo.Map(
						lo.Entries(trs),
						func(entry lo.Entry[string, []pe.TransactionHistory], _ int) lo.Entry[string, pe.TransactionHistory] {
							return lo.Entry[string, pe.TransactionHistory]{Key: entry.Key, Value: entry.Value[0]}
						}),
					func(a lo.Entry[string, pe.TransactionHistory], b lo.Entry[string, pe.TransactionHistory]) bool {
						return a.Value.BookingDateTime.Before(b.Value.BookingDateTime)
					},
				)
				select {
				case out <- TransactionHistoryWithBankName{first.Value, first.Key}:
				case <-ctx.Done():
					return nil
				}
				lastTime = first.Value.BookingDateTime
				trs[first.Key] = trs[first.Key][1:]

			}
		}
		// return nil
	})
	return out
}
func closed[T any](ch chan T) bool {
	_, ok := <-ch
	return !ok

}

// TODO: mb better also send bank name through the channel so that it could be fanned in
// no, then it would slurp the ones i don't want without
func Channel(ctx context.Context, accId string, c client.APIClient, eg *errgroup.Group, now time.Time) chan pe.TransactionHistory {
	ch := make(chan pe.TransactionHistory)
	eg.Go(func() error {
		defer close(ch)
	OUTER:
		for i := 0; i > 1<<31-1; i++ {
			ts, err := c.TransactionsPage(ctx, accId, int32(i), &now, nil)
			if err != nil {
				return err
			}
			for _, t := range ts {
				select {
				case ch <- t:
				case <-ctx.Done():
					break OUTER
				}
			}
		}
		return nil
	})
	return ch
}
