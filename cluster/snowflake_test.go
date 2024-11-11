package cluster

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type client struct {
	redis.UniversalClient

	storage sync.Map
}

func (c *client) SetEx(_ context.Context, key string, value any, _ time.Duration) *redis.StatusCmd {
	val, _ := c.storage.LoadOrStore(key, value)
	return redis.NewStatusResult(val.(string), nil)
}

func (c *client) Set(_ context.Context, key string, value any, _ time.Duration) *redis.StatusCmd {
	c.storage.Store(key, value)
	return redis.NewStatusResult("", nil)
}

func (c *client) Del(_ context.Context, keys ...string) *redis.IntCmd {
	for _, key := range keys {
		c.storage.Delete(key)
	}
	return redis.NewIntResult(int64(len(keys)), nil)
}

func TestSnowflakeNode(t *testing.T) {
	ctx := context.Background()
	cli := &client{}
	node1, err := NewSnowflakeNode(ctx, cli, "service", "node-001")
	if err != nil {
		t.Error(err)
		return
	}
	if node1.index != 0 {
		t.Error("node1.index != 0")
	}
	node2, err := NewSnowflakeNode(ctx, cli, "service", "node-002")
	if err != nil {
		t.Error(err)
		return
	}
	if node2.index != 1 {
		t.Error("node2.index != 1")
	}
	if err = node2.Stop(ctx); err != nil {
		t.Error(err)
	}
	node3, err := NewSnowflakeNode(ctx, cli, "service", "node-003")
	if err != nil {
		t.Error(err)
		return
	}
	if node3.index != 1 {
		t.Error("node2.index != 1")
	}
}
