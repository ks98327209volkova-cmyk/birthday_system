package bootstrap

import (
	"notification_service/config"
	"notification_service/internal/services/reminders"
	"notification_service/internal/storage/redisstorage"
)

func InitNotificationService(storage *redisstorage.RedisStorage, cfg *config.Config) *reminders.Service {
	return reminders.NewService(storage, cfg.Service.ReminderDays)
}
