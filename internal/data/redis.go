package data

import (
	"github.com/redis/go-redis/v9"
)

// NewRedisClient 创建 Redis 客户端实例
func NewRedisClient(d *Data) *redis.Client {
	return d.redis
}
