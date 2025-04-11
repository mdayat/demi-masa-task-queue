package configs

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func NewAsynqClient(redis *redis.Client) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redis.Options().Addr,
		Username: redis.Options().Username,
		Password: redis.Options().Password,
		DB:       redis.Options().DB,
	})
}
