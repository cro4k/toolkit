package clients

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	HeaderClientID  = "x-client-id"
	HeaderVersion   = "x-client-version"
	HeaderClientIP  = "x-client-ip"
	HeaderUA        = "x-client-ua"
	HeaderRequestID = "x-request-id"
)

type ClientInfo struct {
	ID        string `json:"id"`
	IP        string `json:"ip"`
	UA        string `json:"ua"`
	Version   string `json:"v"`
	RequestID string `json:"rid"`
}

func (c *ClientInfo) MD() metadata.MD {
	return metadata.New(map[string]string{
		HeaderClientID:  c.ID,
		HeaderVersion:   c.Version,
		HeaderClientIP:  c.IP,
		HeaderUA:        c.UA,
		HeaderRequestID: c.RequestID,
	})
}

type clientInfoContextKey struct{}

func FromContext(ctx context.Context) *ClientInfo {
	info, ok := ctx.Value(clientInfoContextKey{}).(*ClientInfo)
	if !ok {
		return nil
	}
	return info
}

func WithContext(ctx context.Context, info *ClientInfo) context.Context {
	return context.WithValue(ctx, clientInfoContextKey{}, info)
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(next)
	}
}

func MiddlewareFunc() func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware(handler)
	}
}

func middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := &ClientInfo{
			ID:        def(r.Header.Get(HeaderClientID)),
			IP:        ClientIP(r),
			UA:        r.UserAgent(),
			Version:   r.Header.Get(HeaderVersion),
			RequestID: def(r.Header.Get(HeaderRequestID)),
		}
		ctx := WithContext(r.Context(), info)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return handler(withGRPCServerContext(ctx), req)
	}
}

type ServerStream struct {
	grpc.ServerStream
}

func (s *ServerStream) Context() context.Context {
	return withGRPCServerContext(s.ServerStream.Context())
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &ServerStream{ss})
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(withGRPCClientContext(ctx), method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(withGRPCClientContext(ctx), desc, cc, method, opts...)
	}
}

func withGRPCServerContext(ctx context.Context) context.Context {
	if client := FromContext(ctx); client == nil {
		md, _ := metadata.FromIncomingContext(ctx)
		client = &ClientInfo{
			ID:        defMD(md, HeaderClientID, uuid.New().String()),
			IP:        defMD(md, HeaderClientIP, ""),
			UA:        defMD(md, HeaderUA, ""),
			Version:   defMD(md, HeaderVersion, ""),
			RequestID: defMD(md, HeaderRequestID, uuid.New().String()),
		}
		return WithContext(ctx, client)
	}
	return ctx
}

func withGRPCClientContext(ctx context.Context) context.Context {
	client := FromContext(ctx)
	if client != nil {
		return metadata.NewOutgoingContext(ctx, client.MD())
	}
	return ctx
}

func def(a string) string {
	if a != "" {
		return a
	}
	return uuid.New().String()
}

func defMD(md metadata.MD, key, def string) string {
	if md == nil {
		return def
	}
	a := md.Get(key)
	if len(a) == 0 || a[0] == "" {
		return def
	}
	return a[0]
}
