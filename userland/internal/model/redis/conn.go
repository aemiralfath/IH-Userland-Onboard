package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func NewRedis(config RedisConfig) (*redis.Client, error) {

	connectionCfg := &redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	}

	client := redis.NewClient(connectionCfg)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
