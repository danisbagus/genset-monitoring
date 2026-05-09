package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/config"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	redisClient *goredis.Client
	redisOnce   sync.Once
)

// NewRedis returns a singleton Redis client.
func NewRedis(cfg config.RedisConfig, log *zap.Logger) (*goredis.Client, error) {
	var initErr error

	redisOnce.Do(func() {
		client := goredis.NewClient(&goredis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("failed to connect to redis: %w", err)
			return
		}

		log.Info("Redis connected successfully",
			zap.String("host", cfg.Host),
			zap.String("port", cfg.Port),
		)

		redisClient = client
	})

	return redisClient, initErr
}

// Ping checks if the redis connection is alive.
func Ping(client *goredis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}
