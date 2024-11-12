package cluster

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Endpoint struct {
	service string
	host    string
	port    uint16

	nodeID string

	metadata Metadata
	healthy  *HealthyConfig
}

func (e *Endpoint) Service() string    { return e.service }
func (e *Endpoint) Host() string       { return e.host }
func (e *Endpoint) Port() uint16       { return e.port }
func (e *Endpoint) NodeID() string     { return e.nodeID }
func (e *Endpoint) Metadata() Metadata { return e.metadata }

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
		service:  service,
		host:     host,
		port:     port,
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
		err = GRPCHealthCheck(ctx, e.healthy, e.service, grpc.WithTransportCredentials(insecure.NewCredentials()))
	default:
		return fmt.Errorf("unknown protocol: %v", e.healthy.Protocol.protocol())
	}
	return err
}
