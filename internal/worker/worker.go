package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mdayat/demi-masa-task-queue/configs"
	"github.com/mdayat/demi-masa-task-queue/internal/handlers"
	"github.com/mdayat/demi-masa-task-queue/internal/services"
)

func NewWorkerServer(configs configs.Configs) (*asynq.Server, *asynq.ServeMux, error) {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     configs.Redis.Options().Addr,
			Username: configs.Redis.Options().Username,
			Password: configs.Redis.Options().Password,
			DB:       configs.Redis.Options().DB,
		},
		asynq.Config{Concurrency: 10},
	)

	middleware := handlers.NewMiddlewareHandler(configs)
	mux := asynq.NewServeMux()
	mux.Use(middleware.Logger)

	prayerService := services.NewPrayerService(configs)
	_, err := prayerService.EnqueuePopulatePrayerSchedule(services.PopulatePrayerScheduleType, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to enqueue %s: %w", services.PopulatePrayerScheduleType, err)
	}

	_, err = prayerService.EnqueueUpdateUncheckedPrayer(services.UpdateUncheckedPrayerType, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to enqueue %s: %w", services.UpdateUncheckedPrayerType, err)
	}

	prayerHandler := handlers.NewPrayerHandler(configs, prayerService)
	mux.HandleFunc(services.PopulatePrayerScheduleType, prayerHandler.PopulatePrayerSchedule)
	mux.HandleFunc(services.UpdateUncheckedPrayerType, prayerHandler.UpdateUncheckedPrayer)

	return server, mux, nil
}
