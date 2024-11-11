package cluster

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Endpoint struct {
	Service string
	Host    string
	Port    uint16

	nodeID string

	metadata Metadata
	healthy  *HealthyConfig
}

type EndpointOption func(*Endpoint)

func WithMetadata(key string, val ...string) EndpointOption {
	return func(e *Endpoint) {
		e.metadata.Add(key, val...)
	}
}

func WithHealthy(healthy *HealthyConfig) EndpointOption {
	return func(e *Endpoint) {
		e.healthy = healthy
	}
}

func WithNodeID(nodeID string) EndpointOption {
	return func(e *Endpoint) {
		e.nodeID = nodeID
	}
}

func NewEndpoint(service, host string, port uint16, options ...EndpointOption) *Endpoint {
	endpoint := &Endpoint{
		Service:  service,
		Host:     host,
		Port:     port,
		nodeID:   uuid.New().String(),
		metadata: Metadata{},
	}
	for _, option := range options {
		option(endpoint)
	}
	return endpoint
}

func (e *Endpoint) Healthy(ctx context.Context) error {
	if e.healthy == nil {
		return nil
	}
	var err error
	switch e.healthy.Protocol.protocol() {
	case HTTP:
		err = HTTPHealthCheck(ctx, e.healthy)
	case GRPC:
		err = GRPCHealthCheck(ctx, e.healthy, e.Service, grpc.WithTransportCredentials(insecure.NewCredentials()))
	default:
		return fmt.Errorf("unknown protocol: %v", e.healthy.Protocol.protocol())
	}
	return err
}
