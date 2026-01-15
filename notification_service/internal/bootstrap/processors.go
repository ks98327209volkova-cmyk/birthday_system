package bootstrap

import (
	"notification_service/internal/services/processors/personeventprocessor"
	"notification_service/internal/services/reminders"
)

func InitPersonEventProcessor(service *reminders.Service) *personeventprocessor.PersonEventProcessor {
	return personeventprocessor.NewPersonEventProcessor(service)
}
