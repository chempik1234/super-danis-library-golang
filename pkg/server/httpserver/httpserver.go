package httpserver

import (
	"context"
	"fmt"
	"net/http"
)

// GracefulServerImplementationHTTP - реализация запуска/остановки http.Server - вызывается server.GracefulServer
//
// Создаётся на основе http.Handler и обслуживает порт, указываемый при запуске (NewInstance)
type GracefulServerImplementationHTTP struct {
	router http.Handler
}

// NewGracefulServerImplementationHTTP - создаёт новый GracefulServerImplementationHTTP
func NewGracefulServerImplementationHTTP(router http.Handler) *GracefulServerImplementationHTTP {
	return &GracefulServerImplementationHTTP{router: router}
}

// NewInstance - реализация создания сервера на порту для http.Server - вызывается server.GracefulServer
func (s *GracefulServerImplementationHTTP) NewInstance(port int) (*http.Server, error) {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}, nil
}

// ListenInstance - реализация работы сервера http.Server - начало прослушки
func (s *GracefulServerImplementationHTTP) ListenInstance(instance *http.Server) error {
	return instance.ListenAndServe()
}

// ShutdownInstance - реализация завершения сервера для http.Server - вызывается server.GracefulServer
func (s *GracefulServerImplementationHTTP) ShutdownInstance(ctx context.Context, instance *http.Server) error {
	return instance.Shutdown(ctx)
}
