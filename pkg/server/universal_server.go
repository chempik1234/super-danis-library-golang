package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

// UniversalServer - это порт, под который нужно подогнать сервер для GracefulServer
//
// Например:
// HTTP - метод для создания нового HTTPServer, создания нового и остановки
type UniversalServer[ServerObjectType any] interface {
	// NewInstance - быстрое создание сервера на порту
	NewInstance(port int) (ServerObjectType, error)
	// ListenInstance - запуск сервера в работу, блокирует
	ListenInstance(instance ServerObjectType) error
	// ShutdownInstance - остановка данного сервера (созданного ранее)
	ShutdownInstance(ctx context.Context, instance ServerObjectType) error
}

// GracefulServer - общая структура для сервера, который можно завершить плавно.
//
// В него можно впихнуть HTTP, gRPC или любой другой сервер, и он будет запускаться и завершаться
type GracefulServer[ServerObjectType any] struct {
	serverManager UniversalServer[ServerObjectType]
}

// NewGracefulServer создаёт GracefulServer с данным сервером
func NewGracefulServer[ServerObjectType any](server UniversalServer[ServerObjectType]) *GracefulServer[ServerObjectType] {
	return &GracefulServer[ServerObjectType]{serverManager: server}
}

// GracefulRun запускает GracefulServer и плавно завершает при os.Interrupt или естественной ошибке
//
// “ctx context.Context“ тоже вызывает Graceful Shutdown
//
// 1. Создать структуру сервера. При каждом запуске новая
//
// 2. Подготовить каналы с сигналами
//
// 3. Запустить фоном слушатель для shutdown - именно он и закрывает сервер в нормальных условиях
//
// 4. Запустить сам сервер и ждать, пока он не схлопнется
//
// 5. Подождать завершение слушателя и выйти
func (s *GracefulServer[_]) GracefulRun(ctx context.Context, port int) error {
	// шаг 1.
	serverInstance, err := s.serverManager.NewInstance(port)
	if err != nil {
		return fmt.Errorf("could not create server instance: %w", err)
	}

	// шаг 1.1. Каналы с сигналами о том, что 1) вышел сервер 2) вышла горутина, слушаящая os.Interrupt и ctx
	serverStopped := make(chan bool, 1)
	signalListenerExited := make(chan bool, 1)

	// шаг 2.
	go listenSignal(ctx, s.serverManager, serverInstance, serverStopped, signalListenerExited)

	// шаг 3.
	err = s.serverManager.ListenInstance(serverInstance)
	serverStopped <- true

	// подождём пока не завершится слушатель, и потом выйдем
	<-signalListenerExited // этот сигнал всегда идёт после

	if err != nil {
		return fmt.Errorf("error while listening %v port '%d': %w", serverInstance, port, err)
	}

	return nil
}

func listenSignal[T any](ctx context.Context, serverManager UniversalServer[T], serverObject T, serverStopped <-chan bool, funcExited chan<- bool) {
	// шаг 1. Graceful shutdown через сигнал - прикручиваем к основному контексту
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// шаг 2. Есть 2 источника сигналов
	//
	// 1) signalCtx проверяет: настал os.Signal или "вызывающая сторона" закрыла контекст
	// 2) serverStopped проверяет, что сервер уже отрубился без нас

	select {
	// Если не <1>, а уже <2>, просто выйдем
	case <-serverStopped:
		break
	// Если всё-таки <1>, то плавно отрубаемся (где-то на фоне случится <2>)
	case <-signalCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		someErr := serverManager.ShutdownInstance(shutdownCtx, serverObject)
		if someErr != nil {
			log.Fatalf("error while shutting down http serverManager: %v", someErr)
		}
	}

	funcExited <- true
}
