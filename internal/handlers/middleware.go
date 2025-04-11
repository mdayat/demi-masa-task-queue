package handlers

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/mdayat/demi-masa-task-queue/configs"
	"github.com/rs/zerolog/log"
)

type MiddlewareHandler interface {
	Logger(next asynq.Handler) asynq.Handler
}

type middleware struct {
	configs configs.Configs
}

func NewMiddlewareHandler(configs configs.Configs) MiddlewareHandler {
	return &middleware{
		configs: configs,
	}
}

func (m middleware) Logger(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		logger := log.
			With().
			Str("task_id", task.ResultWriter().TaskID()).
			Str("task_type", task.Type()).
			Logger()

		return next.ProcessTask(logger.WithContext(ctx), task)
	})
}
