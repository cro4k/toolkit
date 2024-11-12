package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

const (
	snowflakeNodeMaxIndex        = 1023
	snowflakeNodeLeaseDuration   = 5 * time.Minute
	snowflakeNodeLeaseExpiration = time.Hour
	snowflakeNodeLeasePrefix     = "SNOWFLAKE_NODE_LEASE"
)

type (
	SnowflakeNodeRunner interface {
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}

	SnowflakeNode struct {
		*snowflake.Node
		service string
		id      string
		index   int64
		client  redis.UniversalClient
	}
)

func NewSnowflakeNodeWithEndpoint(ctx context.Context, client redis.UniversalClient, endpoint *Endpoint) (
	*SnowflakeNode, error,
) {
	return NewSnowflakeNode(ctx, client, endpoint.service, endpoint.nodeID)
}

// NewSnowflakeNode
// 基于 redis SetEx 实现自动分配服务节点序号，保证每个节点序号在服务内唯一，并将节点序号作为 NodeID 创建 snowflake.Node 对象。
// 注意：获取到节点序号后需要定期刷新（运行SnowflakeNode.Start()），保持当前节点对该序号的持有状态。
// 警告：节点序号范围为 0~1023，节点数量超过范围将无法分配序号。同时缓存有一小时过期时间，如果服务内所有节点一小时内累计异常重启次数超过1023次，
// 也可能导致无法分配序号。
func NewSnowflakeNode(ctx context.Context, client redis.UniversalClient, service, id string) (*SnowflakeNode, error) {
	index, err := resolveIndex(ctx, client, service, id)
	if err != nil {
		return nil, err
	}
	node, err := snowflake.NewNode(index)
	if err != nil {
		return nil, err
	}
	snode := &SnowflakeNode{
		Node:    node,
		service: service,
		id:      id,
		index:   index,
		client:  client,
	}
	return snode, nil
}

func resolveIndex(ctx context.Context, client redis.UniversalClient, service, id string) (int64, error) {
	var index int64
	for index < snowflakeNodeMaxIndex {
		key := fmt.Sprintf("%s:%s:%d", snowflakeNodeLeasePrefix, service, index)
		val, err := client.SetEx(ctx, key, id, snowflakeNodeLeaseExpiration).Result()
		if err != nil {
			return 0, err
		}
		if val == id {
			return index, nil
		}
		index++
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(10 * time.Millisecond):
			continue
		}
	}
	return 0, fmt.Errorf("node index is out of range")
}

func (s *SnowflakeNode) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(snowflakeNodeLeaseDuration):
			key := fmt.Sprintf("%s:%s:%d", snowflakeNodeLeasePrefix, s.service, s.index)
			s.client.Set(ctx, key, s.id, snowflakeNodeLeaseExpiration)
		}
	}
}

func (s *SnowflakeNode) Stop(ctx context.Context) error {
	key := fmt.Sprintf("%s:%s:%d", snowflakeNodeLeasePrefix, s.service, s.index)
	return s.client.Del(ctx, key).Err()
}

func (s *SnowflakeNode) Runner() SnowflakeNodeRunner {
	return s
}
