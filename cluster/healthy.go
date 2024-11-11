package cluster

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Protocol interface {
	protocol() protocol

	Is(v protocol) bool
}

// unexported type is aims to avoid this scene:
//
//	protocol("any-text-user-input-can-be-convert-to-this-type")
//
// if you need the type define, just use Protocol.
// example:
//
//	package another
//
//	func hello() {
//		var p Protocol
//		p = cluster.HTTP
//	}
//
// and now, the user can only use the defined enums.
type protocol string

func (p protocol) protocol() protocol { return p }

func (p protocol) Is(v protocol) bool {
	return p == v
}

const (
	HTTP protocol = "http"
	GRPC protocol = "grpc"
)

type HealthyConfig struct {
	Protocol Protocol
	Target   string
	Options  Metadata
}

type HealthChecker interface {
	Healthy(ctx context.Context) error
}

type HealthyConfigOption func(*HealthyConfig)

func WithHealthyConfigOptions(key string, val ...string) HealthyConfigOption {
	return func(config *HealthyConfig) {
		config.Options.Add(key, val...)
	}
}

func NewHealthyConfig(protocol Protocol, target string, options ...HealthyConfigOption) *HealthyConfig {
	c := &HealthyConfig{
		Protocol: protocol,
		Target:   target,
		Options:  Metadata{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func HTTPHealthCheck(ctx context.Context, healthyConfig *HealthyConfig) error {
	method := healthyConfig.Options.Get("method")
	if method == "" {
		method = http.MethodGet
	}
	request, err := http.NewRequestWithContext(ctx, method, healthyConfig.Target, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func GRPCHealthCheck(ctx context.Context, healthyConfig *HealthyConfig, service string, options ...grpc.DialOption) error {
	cc, err := grpc.DialContext(ctx, healthyConfig.Target, options...)
	if err != nil {
		return err
	}
	defer cc.Close()
	client := grpc_health_v1.NewHealthClient(cc)
	response, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: service})
	if err != nil {
		return err
	}
	if response.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return errors.New(response.Status.String())
	}
	return nil
}
