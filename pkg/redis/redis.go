package redis

import (
	"context"
	"fmt"
	"github.com/chempik1234/super-danis-library-golang/v2/pkg/logger"
	"github.com/go-redis/redis/v8"
	"time"
)

// Config is the redis connection config struct
type Config struct {
	Addr       string `env:"ADDR" env-default:"localhost:6379"`
	Password   string `env:"PASSWORD" env-default:""`
	DB         int    `env:"DB" env-default:"0"`
	TTLSeconds int    `env:"TTL_SECONDS" env-default:"0"`

	Timeout TimeoutConfig `yaml:"timeout" env-prefix:"TIMEOUT_"`

	Retries RetriesConfig `yaml:"retries" env-prefix:"RETRIES_"`

	Pool PoolConfig `yaml:"pool" env-prefix:"POOL_"`
}

// PoolConfig - pool for redis.Config
type PoolConfig struct {
	Size               int `yaml:"size" env:"SIZE" env-default:"3"`
	MinIdleConnections int `yaml:"min_idle_connections" env:"MIN_IDLE_CONNECTIONS" env-default:"2"`
}

// TimeoutConfig - timeouts for redis.Config
type TimeoutConfig struct {
	DialMilliseconds  int `yaml:"dial_milliseconds" env:"DIAL_MILLISECONDS" env-default:"5000"`
	ReadMilliseconds  int `yaml:"read_milliseconds" env:"READ_MILLISECONDS" env-default:"1000"`
	WriteMilliseconds int `yaml:"write_milliseconds" env:"WRITE_MILLISECONDS" env-default:"1000"`
}

// RetriesConfig - retries for redis.Config
type RetriesConfig struct {
	MaxRetries int `env:"MAX_RETRIES" env-default:"3"`
}

// New - create new redis conn with using given config
//
// Doesn't use all options
//
//	client, err := New(ctx, config)
//	if err != nil {
//	   logger.GetLoggerFromCtx(ctx).Error(ctx, "aw hell no")
//	   return
//	}
//	defer DeferDisconnect(ctx, client)
func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.Retries.MaxRetries,
		DialTimeout:  time.Duration(cfg.Timeout.DialMilliseconds) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.Timeout.ReadMilliseconds) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.Timeout.WriteMilliseconds) * time.Millisecond,
		PoolSize:     cfg.Pool.Size,
		MinIdleConns: cfg.Pool.MinIdleConnections,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error pinging redis: %w", err)
	}
	return redisClient, nil
}

// DeferDisconnect - call in defer after getting client.
//
//	client, err := New(ctx, config)
//	if err != nil {
//	   logger.GetLoggerFromCtx(ctx).Error(ctx, "aw hell no")
//	   return
//	}
//	defer DeferDisconnect(ctx, client)
func DeferDisconnect(ctx context.Context, client *redis.Client) {
	if err := client.Close(); err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Error(ctx, "failed to disconnect from mongodb client")
	}
}
