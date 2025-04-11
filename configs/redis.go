package configs

import (
	"github.com/redis/go-redis/v9"
)

func NewRedis(REDIS_URL string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: REDIS_URL,
	})
}
