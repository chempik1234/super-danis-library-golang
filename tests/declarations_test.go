package tests

import (
	"github.com/chempik1234/super-danis-library-golang/v2/pkg/server"
	"github.com/chempik1234/super-danis-library-golang/v2/pkg/server/grpcserver"
	"github.com/chempik1234/super-danis-library-golang/v2/pkg/server/httpserver"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"testing"
)

type httpHandler struct{}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World"))
}

func TestInterfaces(t *testing.T) {
	grpcServer := server.NewGracefulServer[*net.Listener](
		grpcserver.NewGracefulServerImplementationGRPC(&grpc.Server{}))
	httpServer := server.NewGracefulServer[*http.Server](
		httpserver.NewGracefulServerImplementationHTTP(httpHandler{}))

	t.Log("successful interfaces", grpcServer, httpServer)
}
