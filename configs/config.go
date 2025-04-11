package configs

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Configs struct {
	Env         Env
	Db          Db
	Redis       *redis.Client
	AsynqClient *asynq.Client
}

func NewConfigs(env Env, db Db, redis *redis.Client, asynqClient *asynq.Client) Configs {
	return Configs{
		Env:         env,
		Db:          db,
		Redis:       redis,
		AsynqClient: asynqClient,
	}
}
