package mongodb

import (
	"context"
	"fmt"
	"github.com/chempik1234/super-danis-library-golang/v2/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Config - required to connect to MongoDB with this library
type Config struct {
	Hosts       []string `env:"HOSTS" yaml:"hosts" envDefault:"mongodb:27017"`
	MinPoolSize uint64   `env:"MIN_POOL_SIZE" yaml:"min_pool_size" envDefault:"1"`
	MaxPoolSize uint64   `env:"MAX_POOL_SIZE" yaml:"max_pool_size" envDefault:"10"`

	UserName    string `env:"USERNAME" yaml:"username" envDefault:"root"`
	Password    string `env:"PASSWORD" yaml:"password" envDefault:"root"`
	PasswordSet bool   `env:"PASSWORD_SET" yaml:"password_set" envDefault:"false"`

	RetryWrites bool `env:"RETRY_WRITES" yaml:"retry_writes" envDefault:"true"`
	RetryReads  bool `env:"RETRY_READS" yaml:"retry_reads" envDefault:"true"`
}

// New - connect to MongoDB with given config
//
// tries to PING after connecting
//
//	client, err := New(ctx, config)
//	if err != nil {
//	   logger.GetLoggerFromCtx(ctx).Error(ctx, "aw hell no")
//	   return
//	}
//	defer DeferDisconnect(ctx, client)
func New(ctx context.Context, config Config) (*mongo.Client, error) {
	client, err := mongo.Connect(
		options.Client().SetMinPoolSize(config.MinPoolSize),
		options.Client().SetMaxPoolSize(config.MaxPoolSize),
		options.Client().SetHosts(config.Hosts),
		options.Client().SetAuth(options.Credential{
			Username:    config.UserName,
			Password:    config.Password,
			PasswordSet: config.PasswordSet,
		}),
		options.Client().SetRetryWrites(config.RetryWrites),
		options.Client().SetRetryReads(config.RetryReads),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}
	return client, nil
}

// DeferDisconnect - call in defer after getting client.
//
//	client, err := New(ctx, config)
//	if err != nil {
//	   logger.GetLoggerFromCtx(ctx).Error(ctx, "aw hell no")
//	   return
//	}
//	defer DeferDisconnect(ctx, client)
func DeferDisconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Error(ctx, "failed to disconnect from mongodb client")
	}
}
