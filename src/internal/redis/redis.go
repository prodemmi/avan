package redis

import (
	"avan/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Redis struct {
	Config config.Config
	client *redis.Client
}

const (
	MessageQueueStream = "avan:messages"
	MessageQueueGroup  = "avan:consumers"
)

func NewRedis(config config.Config, logger zap.Logger) Redis {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	if res := client.Ping(ctx); res.Err() != nil {
		logger.Fatal("failed to ping redis", zap.Error(res.Err()))
	}

	return Redis{Config: config, client: client}
}

func (r Redis) Client() *redis.Client {
	return r.client
}
