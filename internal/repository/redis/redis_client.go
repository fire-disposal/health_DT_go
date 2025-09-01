// redis_client.go
package redis

import (
	"context"
	"sync"

	"github.com/fire-disposal/health_DT_go/config"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// InitRedisClient 初始化 Redis 客户端（单例）
func InitRedisClient(cfg *config.RedisConfig) *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
			PoolSize: 20, // 高并发连接池
		})
	})
	return client
}

// GetRedisClient 获取 Redis 客户端
func GetRedisClient() *redis.Client {
	return client
}

// PingRedis 检查连接
func PingRedis(ctx context.Context) error {
	return client.Ping(ctx).Err()
}
