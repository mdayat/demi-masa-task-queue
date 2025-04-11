package main

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/mdayat/demi-masa-task-queue/configs"
	"github.com/mdayat/demi-masa-task-queue/internal/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	logger := log.With().Caller().Logger()

	env, err := configs.LoadEnv()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	ctx := context.TODO()
	db, err := configs.NewDb(ctx, env.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	defer db.Conn.Close()

	redis, err := configs.NewRedis(env.RedisURL)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	defer redis.Close()

	asynqClient := configs.NewAsynqClient(redis)
	defer asynqClient.Close()

	configs := configs.NewConfigs(env, db, redis, asynqClient)
	server, mux, err := worker.NewWorkerServer(configs)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	if err := server.Run(mux); err != nil {
		logger.Fatal().Err(err).Send()
	}
}
