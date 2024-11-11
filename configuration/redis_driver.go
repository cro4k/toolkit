package configuration

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	client redis.UniversalClient
	prefix string
}

type RedisDriverOption func(*RedisDriver)

func WithRedisDriverPrefix(prefix string) RedisDriverOption {
	return func(r *RedisDriver) {
		r.prefix = prefix
	}
}

func NewRedisDriver(client redis.UniversalClient, options ...RedisDriverOption) *RedisDriver {
	r := &RedisDriver{client: client}
	for _, option := range options {
		option(r)
	}
	return r
}

func (d *RedisDriver) Load(ctx context.Context, key string) ([]byte, string, error) {
	key = d.prefix + key
	data, err := d.client.Get(ctx, key).Bytes()
	return data, "", err
}
