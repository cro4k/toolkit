package pprof

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/google/gops/agent"
)

func RunAsync() {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Errorf("start pprof on error: %w", err))
	}
	log.Println("start pprof service on:", ln.Addr())
	go func() {
		_ = http.Serve(ln, nil)
	}()

	if err = agent.Listen(agent.Options{ShutdownCleanup: false}); err != nil {
		panic(fmt.Errorf("start pprof on error: %w", err))
	}
}
