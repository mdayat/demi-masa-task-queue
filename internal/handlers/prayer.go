package handlers

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/mdayat/demi-masa-task-queue/configs"
	"github.com/mdayat/demi-masa-task-queue/internal/services"
)

type PrayerHandler interface {
	PopulatePrayerSchedule(ctx context.Context, task *asynq.Task) error
	UpdateUncheckedPrayer(ctx context.Context, task *asynq.Task) error
}

const (
	PopulatePrayerScheduleType = "prayer:populate"
	UpdateUncheckedPrayerType  = "prayer:update"
)

type prayer struct {
	configs configs.Configs
	service services.PrayerServicer
}

func NewPrayerHandler(configs configs.Configs, service services.PrayerServicer) PrayerHandler {
	return &prayer{
		configs: configs,
		service: service,
	}
}

func (p prayer) PopulatePrayerSchedule(ctx context.Context, task *asynq.Task) error {
	return nil
}

func (p prayer) UpdateUncheckedPrayer(ctx context.Context, task *asynq.Task) error {
	return nil
}
