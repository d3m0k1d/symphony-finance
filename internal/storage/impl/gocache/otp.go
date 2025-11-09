package gocache

import (
	"context"
	"time"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

var _ interfaces.OtpSessionStore = otpCache{}

type otpCache struct {
	c cache.CacheInterface[[]byte]
}

func NewOtpCache() (*otpCache, error) {
	r, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 10_000,
		MaxCost:     0,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	s := ristretto_store.NewRistretto(r)
	c := cache.New[[]byte](s)
	return &otpCache{
		c: c,
	}, nil
}

// RetrieveCode implements interfaces.OtpSessionStore.
func (c otpCache) RetrieveCode(ctx context.Context, session string) (hash []byte, err error) {
	return c.c.Get(ctx, session)
}

// StoreCode implements interfaces.OtpSessionStore.
func (c otpCache) StoreCode(ctx context.Context, session string, hash []byte) error {
	return c.c.Set(ctx, session, hash, store.WithExpiration(10*time.Minute))
}
