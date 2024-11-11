package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client redis.UniversalClient
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheNotFound
	}
	return val, err
}

func (c *redisCache) Set(ctx context.Context, key string, value []byte, exp time.Duration) error {
	return c.client.Set(ctx, key, value, exp).Err()
}

func (c *redisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func NewRedisCache(client redis.UniversalClient) Cache {
	return &redisCache{client: client}
}
