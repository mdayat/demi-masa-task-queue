package configs

import (
	"github.com/redis/go-redis/v9"
)

type Configs struct {
	Env   Env
	Db    Db
	Redis *redis.Client
}

func NewConfigs(env Env, db Db, redis *redis.Client) Configs {
	return Configs{
		Env:   env,
		Db:    db,
		Redis: redis,
	}
}
