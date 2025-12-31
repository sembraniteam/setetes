package redisx

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sembraniteam/setetes/internal"
)

var log = slog.Default()

type (
	client struct {
		config internal.Config
	}

	Redis interface {
		Connect() (*redis.Client, error)
		Disconnect(rdb *redis.Client)
	}
)

func New() Redis {
	return client{config: *internal.Get()}
}

func (c client) Connect() (*redis.Client, error) {
	redisConf := c.config.Redis

	// `tls.Config` implementation is planned but has not been included in the
	// current scope.
	cl := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Username:        redisConf.Username,
		Password:        redisConf.Password,
		DB:              redisConf.Database,
		MaxRetries:      redisConf.MaxRetries,
		DialTimeout:     redisConf.DialTimeout * time.Second,
		ReadTimeout:     redisConf.ReadTimeout * time.Second,
		WriteTimeout:    redisConf.WriteTimeout * time.Second,
		PoolSize:        redisConf.PoolSize,
		PoolTimeout:     redisConf.PoolTimeout * time.Second,
		MinIdleConns:    redisConf.MinIdleConns,
		MaxIdleConns:    redisConf.MaxIdleConns,
		ConnMaxLifetime: redisConf.ConnMaxLifetime * time.Minute,
		ConnMaxIdleTime: redisConf.ConnMaxIdleTime * time.Minute,
	})

	_, err := cl.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c client) Disconnect(rdb *redis.Client) {
	if err := rdb.Close(); err != nil {
		log.Error(
			"Error disconnecting Redis client",
			slog.String("error", err.Error()),
		)
	}
}
