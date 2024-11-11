package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrCacheNotFound = errors.New("cache not found")
)

type (
	Cache interface {
		Get(ctx context.Context, key string) ([]byte, error)
		Set(ctx context.Context, key string, value []byte, exp time.Duration) error
		Del(ctx context.Context, key string) error
	}

	ObjectCache interface {
		Get(ctx context.Context, key string, dst any) error
		Set(ctx context.Context, key string, val any, exp time.Duration) error
		Del(ctx context.Context, key string) error
	}
)

type cacheWithOptions struct {
	cache Cache

	ignoreNotFound bool
	prefix         string
}

func (c *cacheWithOptions) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.cache.Get(ctx, c.prefix+key)
	if c.ignoreNotFound && errors.Is(err, ErrCacheNotFound) {
		err = nil
	}
	return value, err
}
func (c *cacheWithOptions) Set(ctx context.Context, key string, value []byte, exp time.Duration) error {
	return c.cache.Set(ctx, c.prefix+key, value, exp)
}

func (c *cacheWithOptions) Del(ctx context.Context, key string) error {
	return c.cache.Del(ctx, c.prefix+key)
}

type Option func(o *cacheWithOptions)

func WithPrefix(prefix string) Option {
	return func(o *cacheWithOptions) {
		o.prefix = prefix
	}
}

func WithIgnoreNotFound() Option {
	return func(o *cacheWithOptions) {
		o.ignoreNotFound = true
	}
}

func With(cache Cache, options ...Option) Cache {
	o := &cacheWithOptions{cache: cache}
	for _, option := range options {
		option(o)
	}
	return o
}

type objectCache struct {
	Cache
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func (c *objectCache) Get(ctx context.Context, key string, dst any) error {
	data, err := c.Cache.Get(ctx, key)
	if err != nil {
		return err
	}
	return c.unmarshal(data, dst)
}

func (c *objectCache) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	data, err := c.marshal(val)
	if err != nil {
		return err
	}
	return c.Cache.Set(ctx, key, data, exp)
}

func JSONCache(cache Cache) ObjectCache {
	return &objectCache{
		Cache:     cache,
		marshal:   json.Marshal,
		unmarshal: json.Unmarshal,
	}
}

func NewObjectCache(cache Cache, marshal func(any) ([]byte, error), unmarshal func([]byte, any) error) ObjectCache {
	return &objectCache{
		Cache:     cache,
		marshal:   marshal,
		unmarshal: unmarshal,
	}
}
