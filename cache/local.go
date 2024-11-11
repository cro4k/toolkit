package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type (
	cacheElement struct {
		key  string
		data []byte
		exp  time.Time
	}

	localCache struct {
		list    list.List
		storage sync.Map
		maxLen  int
	}
)

func (c *localCache) Get(_ context.Context, key string) ([]byte, error) {
	val, ok := c.storage.Load(key)
	if !ok {
		return nil, ErrCacheNotFound
	}
	e, _ := val.(*list.Element)
	element, _ := e.Value.(*cacheElement)
	if element.exp.Before(time.Now()) {
		c.list.Remove(e)
		c.storage.Delete(key)
		return nil, ErrCacheNotFound
	}
	return element.data, nil
}

func (c *localCache) Set(_ context.Context, key string, value []byte, exp time.Duration) error {
	val, ok := c.storage.Load(key)
	if ok {
		e, _ := val.(*list.Element)
		element, _ := e.Value.(*cacheElement)
		element.data = value
		element.exp = time.Now().Add(exp)
		return nil
	}
	if c.list.Len() >= c.maxLen {
		e := c.list.Front()
		element, _ := e.Value.(*cacheElement)
		c.list.Remove(e)
		c.storage.Delete(element.key)
	}
	element := &cacheElement{key: key, data: value, exp: time.Now().Add(exp)}
	c.storage.Store(key, c.list.PushBack(element))
	return nil
}

func (c *localCache) Del(_ context.Context, key string) error {
	e, ok := c.storage.Load(key)
	if !ok {
		return nil
	}
	c.storage.Delete(key)
	c.list.Remove(e.(*list.Element))
	return nil
}

func NewLocalCache(max int) Cache {
	return &localCache{maxLen: max}
}
