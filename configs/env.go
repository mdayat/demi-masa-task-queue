package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseURL string
	RedisURL    string
}

func LoadEnv(filenames ...string) (Env, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return Env{}, err
	}

	env := Env{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}

	return env, nil
}
