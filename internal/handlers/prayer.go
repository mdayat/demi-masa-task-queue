package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mdayat/demi-masa-task-queue/configs"
	"github.com/mdayat/demi-masa-task-queue/internal/services"
	"github.com/mdayat/demi-masa-task-queue/repository"
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

type prayerName string

const (
	subuh  prayerName = "subuh"
	zuhur  prayerName = "zuhur"
	asar   prayerName = "asar"
	magrib prayerName = "magrib"
	isya   prayerName = "isya"
)

func (p prayer) PopulatePrayerSchedule(ctx context.Context, _ *asynq.Task) error {
	users, err := p.configs.Db.Queries.SelectUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to select users: %w", err)
	}

	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)

	firstDayOfNextTwoMonths := time.Date(nextMonth.Year(), nextMonth.Month()+1, 1, 0, 0, 0, 0, nextMonth.Location())
	lastDayOfNextMonth := firstDayOfNextTwoMonths.AddDate(0, 0, -1)
	numOfDaysOfNextMonth := lastDayOfNextMonth.Day()

	numOfPrayersDaily := 5
	prayers := make([]repository.InsertUserPrayersParams, 0, numOfDaysOfNextMonth*numOfPrayersDaily)

	var prayerName prayerName
	for day := 1; day <= numOfDaysOfNextMonth; day++ {
		for i := 1; i <= numOfPrayersDaily; i++ {
			switch i {
			case 1:
				prayerName = subuh
			case 2:
				prayerName = zuhur
			case 3:
				prayerName = asar
			case 4:
				prayerName = magrib
			case 5:
				prayerName = isya
			}

			prayers = append(prayers, repository.InsertUserPrayersParams{
				ID:     pgtype.UUID{},
				UserID: pgtype.UUID{},
				Name:   string(prayerName),
				Year:   int16(nextMonth.Year()),
				Month:  int16(nextMonth.Month()),
				Day:    int16(day),
			})
		}
	}

	for _, user := range users {
		for i := range prayers {
			prayers[i].ID = pgtype.UUID{Bytes: uuid.New(), Valid: true}
			prayers[i].UserID = user.ID
		}

		_, err = p.configs.Db.Queries.InsertUserPrayers(ctx, prayers)
		if err != nil {
			return fmt.Errorf("failed to bulk insert prayers: %w", err)
		}
	}

	return nil
}

func (p prayer) UpdateUncheckedPrayer(ctx context.Context, _ *asynq.Task) error {
	now := time.Now()
	lastTwoDays := now.AddDate(0, 0, -2)

	err := p.configs.Db.Queries.UpdateUncheckedPrayer(ctx, repository.UpdateUncheckedPrayerParams{
		Day:   int16(lastTwoDays.Day()),
		Month: int16(lastTwoDays.Month()),
		Year:  int16(lastTwoDays.Year()),
	})

	if err != nil {
		return fmt.Errorf("failed to update unchecked prayer to missed: %w", err)
	}

	return nil
}
