package kvcacher

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis creates new redis-based cacher
func NewRedis(redisClient redis.UniversalClient, ttl time.Duration, prefix string) *Redis {
	return &Redis{cli: redisClient, ttl: ttl, prefix: prefix}
}

// Redis represents in-memory cacher
type Redis struct {
	sync.Mutex
	cli    redis.UniversalClient
	ttl    time.Duration
	prefix string
}

// Get item from cache
func (m *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	m.Lock()
	defer m.Unlock()

	b, err := m.cli.Get(ctx, m.prefix+key).Bytes()
	switch {
	case errors.Is(err, redis.Nil):
		return nil, nil
	case err == nil:
		return b, nil
	}
	return nil, err
}

// Set item into cache
func (m *Redis) Set(ctx context.Context, key string, b []byte) error {
	m.Lock()
	defer m.Unlock()
	return RedisErrOrNil(m.cli.Set(ctx, m.prefix+key, b, m.ttl).Err())
}

// Clear the cache
func (m *Redis) Clear(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	iterator := m.cli.Scan(ctx, 0, m.prefix+"*", 0).Iterator()
	for iterator.Next(ctx) {
		if err := RedisErrOrNil(m.cli.Del(ctx, iterator.Val()).Err()); err != nil {
			return err
		}
	}
	return RedisErrOrNil(iterator.Err())
}
