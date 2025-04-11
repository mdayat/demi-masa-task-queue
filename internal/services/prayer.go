package services

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/mdayat/demi-masa-task-queue/configs"
)

type PrayerServicer interface {
	EnqueuePopulatePrayerSchedule(pattern string, payload interface{}) (*asynq.TaskInfo, error)
	EnqueueUpdateUncheckedPrayer(pattern string, payload interface{}) (*asynq.TaskInfo, error)
}

type prayer struct {
	configs configs.Configs
}

func NewPrayerService(configs configs.Configs) PrayerServicer {
	return &prayer{
		configs: configs,
	}
}

func (p prayer) EnqueuePopulatePrayerSchedule(pattern string, payload interface{}) (*asynq.TaskInfo, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	now := time.Now()
	firstDayOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	lastDayOfThisMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

	task := asynq.NewTask(pattern, payloadJSON)
	return p.configs.AsynqClient.Enqueue(task, asynq.ProcessIn(lastDayOfThisMonth.Sub(now)))
}

func (p prayer) EnqueueUpdateUncheckedPrayer(pattern string, payload interface{}) (*asynq.TaskInfo, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	now := time.Now()
	midnightOfThisDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	midnightOfOvertomorrow := midnightOfThisDay.AddDate(0, 0, +2)

	task := asynq.NewTask(pattern, payloadJSON)
	return p.configs.AsynqClient.Enqueue(task, asynq.ProcessIn(midnightOfOvertomorrow.Sub(now)))
}
