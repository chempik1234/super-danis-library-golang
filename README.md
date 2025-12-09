## Библиотека для меня

## LRU

```go
type OrderCache pkgports.Cache[string, models.Order]

func NewOrderCacheAdapterInMemoryLRU(capacity int) ports.OrderCache {
    return lru.NewCacheLRUInMemory[string, models.Order](capacity)
}

result, found, err := s.cache.Get(ctx, orderUID)
if err != nil {
	return models.Order{}, fmt.Errorf("error checking orders cache: %w", err)
}

go func() {
    cacheErr := s.cache.Set(ctx, result.OrderUID, result)
    if err != nil {
        logger.GetLoggerFromCtx(ctx).Error(ctx, "error caching order",
        zap.String("key", orderUID), zap.Error(cacheErr))
    }
}()
```

## Server

Пример

```go
package main

import (
	"context"
	"fmt"
	"github.com/chempik1234/super-danis-library-golang/pkg/server"
	"github.com/chempik1234/super-danis-library-golang/pkg/server/httpserver"
	"log"
	"net/http"
	"sync"
)

func main() {
	ctx, stopCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Ваши сервисы, работающие на контексте и wg

	// yourRouter := ...

	appServer := server.NewGracefulServer[*http.Server](
		httpserver.NewGracefulServerImplementationHTTP(yourRouter),
	)
	err := appServer.GracefulRun(ctx, 8080)

	//region shutdown
	if err != nil {
		log.Println(fmt.Errorf("http server error: %w", err).Error())
	}

	log.Println("server gracefully stopped")

	stopCtx()
	wg.Wait()
	log.Println("background operations gracefully stopped")
	//endregion
}

```