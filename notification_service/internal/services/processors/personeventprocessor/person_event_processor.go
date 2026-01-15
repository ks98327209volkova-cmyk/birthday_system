package personeventprocessor

import (
	"context"

	"notification_service/internal/models"
)

type notificationService interface {
	HandlePersonEvent(ctx context.Context, event *models.PersonEvent) error
}

type PersonEventProcessor struct {
	notificationService notificationService
}

func NewPersonEventProcessor(notificationService notificationService) *PersonEventProcessor {
	return &PersonEventProcessor{
		notificationService: notificationService,
	}
}
