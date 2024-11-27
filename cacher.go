package kvcacher

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// KVCacher abstracts KV cache
type KVCacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, b []byte) error
	Clear(ctx context.Context) error
}

// RedisErrOrNil returns nil if redis returned NIL error, otherwise passes
func RedisErrOrNil(err error) error {
	if errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}
