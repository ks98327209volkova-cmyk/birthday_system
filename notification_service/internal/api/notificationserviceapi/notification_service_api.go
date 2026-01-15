package notificationserviceapi

import (
	"context"

	"notification_service/internal/models"
	notificationpb "notification_service/internal/pb/notification_api"
)

type notificationService interface {
	GetAllReminders(ctx context.Context) ([]*models.Reminder, error)
}

type NotificationServiceAPI struct {
	notificationpb.UnimplementedNotificationServiceServer
	notificationService notificationService
}

func NewNotificationServiceAPI(notificationService notificationService) *NotificationServiceAPI {
	return &NotificationServiceAPI{
		notificationService: notificationService,
	}
}
