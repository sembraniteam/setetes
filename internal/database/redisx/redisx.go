package redisx

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sembraniteam/setetes/internal"
)

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
	cl := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Username:        redisConf.Username,
		Password:        redisConf.Password,
		DB:              redisConf.Database,
		MaxRetries:      redisConf.MaxRetries,
		DialTimeout:     redisConf.DialTimeout,
		ReadTimeout:     redisConf.ReadTimeout,
		WriteTimeout:    redisConf.WriteTimeout,
		PoolSize:        redisConf.PoolSize,
		PoolTimeout:     redisConf.PoolTimeout,
		MinIdleConns:    redisConf.MinIdleConns,
		MaxIdleConns:    redisConf.MaxIdleConns,
		MaxActiveConns:  redisConf.MaxActiveConns,
		ConnMaxLifetime: redisConf.ConnMaxLifetime,
		ConnMaxIdleTime: redisConf.ConnMaxIdleTime,
		TLSConfig: &tls.Config{
			ServerName:         redisConf.Host,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
			ClientAuth:         tls.RequireAndVerifyClientCert,
		},
	})

	_, err := cl.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c client) Disconnect(rdb *redis.Client) {
	defer func(rdb *redis.Client) {
		if err := rdb.Close(); err != nil {
			fmt.Printf("error closing Redis client: %v\n", err)
		}
	}(rdb)
}
