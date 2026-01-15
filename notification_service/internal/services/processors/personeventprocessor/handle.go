package personeventprocessor

import (
	"context"

	"notification_service/internal/models"
)

func (p *PersonEventProcessor) Handle(ctx context.Context, event *models.PersonEvent) error {
	return p.notificationService.HandlePersonEvent(ctx, event)
}
