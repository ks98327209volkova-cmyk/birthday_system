package bootstrap

import (
	"notification_service/internal/api/notificationserviceapi"
	"notification_service/internal/services/reminders"
)

func InitNotificationServiceAPI(service *reminders.Service) *notificationserviceapi.NotificationServiceAPI {
	return notificationserviceapi.NewNotificationServiceAPI(service)
}
