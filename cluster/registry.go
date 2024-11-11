package cluster

import "context"

type Registry interface {
	Register(ctx context.Context, endpoint *Endpoint) error
	Deregister(ctx context.Context, endpoint *Endpoint) error
}
