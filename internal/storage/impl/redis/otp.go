package redis

import (
	"context"
	"time"
	"vtb-apihack-2025/internal/storage/interfaces"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

var _ interfaces.OtpSessionStore = otpCache{}

type otpCache struct {
	c *marshaler.Marshaler
}

// DropCode implements interfaces.OtpSessionStore.
func (c otpCache) DropCode(ctx context.Context, session string) (err error) {
	return c.c.Delete(ctx, session)
}

func NewOtpCache(addr string) (*otpCache, error) {
	s := redis_store.NewRedis(redis.NewClient(&redis.Options{Addr: addr}), store.WithExpiration(10*time.Minute))
	c := cache.New[any](s)
	m := marshaler.New(c)
	return &otpCache{
		c: m,
	}, nil
}

// RetrieveCode implements interfaces.OtpSessionStore.
func (c otpCache) RetrieveCode(ctx context.Context, session string) (hash interfaces.OTPSession, err error) {
	_, err = c.c.Get(ctx, session, &hash)
	return
}

// StoreCode implements interfaces.OtpSessionStore.
func (c otpCache) StoreCode(ctx context.Context, session, email string, hash []byte) error {
	return c.c.Set(ctx, session, interfaces.OTPSession{
		Hash:  hash,
		Email: email,
	})
}
