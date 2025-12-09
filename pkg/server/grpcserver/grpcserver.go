package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

// GracefulServerImplementationGRPC - реализация запуска/остановки grpc.Server - вызывается server.GracefulServer
//
// Создаётся на основе grpc.Server и обслуживает порт, указываемый при запуске (NewInstance)
//
// ## ВНИМАНИЕ!!!
// Поддерживает 1 инстанс т.к. создавать и регистрировать grpc.Server под каждый инстанс очень сложно для вас, сори
//
// Пример:
//
//	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AddLogMiddleware))
//
//	test.RegisterOrderServiceServer(server, gRPCServer)
//
//	appServer := server.NewGracefulServer[*net.Listener](
//		grpcserver.NewGracefulServerImplementationGRPC(
//			gRPCServer
//		),
//	)
type GracefulServerImplementationGRPC struct {
	server *grpc.Server
}

// NewGracefulServerImplementationGRPC - создаёт новый GracefulServerImplementationGRPC
func NewGracefulServerImplementationGRPC(server *grpc.Server) *GracefulServerImplementationGRPC {
	return &GracefulServerImplementationGRPC{server: server}
}

// NewInstance - реализация создания сервера на порту для grpc.Server - вызывается server.GracefulServer
func (s *GracefulServerImplementationGRPC) NewInstance(port int) (*net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, fmt.Errorf("err creating listener for graceful gRPC: %w", err)
	}
	return &lis, nil
}

// ListenInstance - реализация работы сервера grpc.Server - начало прослушки
func (s *GracefulServerImplementationGRPC) ListenInstance(instance net.Listener) error {
	return s.server.Serve(instance)
}

// ShutdownInstance - реализация завершения сервера для grpc.Server - вызывается server.GracefulServer
func (s *GracefulServerImplementationGRPC) ShutdownInstance(context.Context, *grpc.Server) error {
	s.server.GracefulStop()
	return nil
}
