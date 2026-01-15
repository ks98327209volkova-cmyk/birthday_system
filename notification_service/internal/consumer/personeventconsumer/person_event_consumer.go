package personeventconsumer

import (
	"context"

	"notification_service/internal/models"
)

type personEventProcessor interface {
	Handle(ctx context.Context, event *models.PersonEvent) error
}

type PersonEventConsumer struct {
	processor   personEventProcessor
	kafkaBroker []string
	topicName   string
}

func NewPersonEventConsumer(processor personEventProcessor, kafkaBroker []string, topicName string) *PersonEventConsumer {
	return &PersonEventConsumer{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
	}
}
