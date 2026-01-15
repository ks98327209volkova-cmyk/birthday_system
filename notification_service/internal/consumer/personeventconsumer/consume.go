package personeventconsumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"notification_service/internal/models"

	"github.com/segmentio/kafka-go"
)

func (c *PersonEventConsumer) Consume(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "NotificationService_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("read message error: %v", err)
			continue
		}

		var event *models.PersonEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("parse message error: %v", err)
			continue
		}

		log.Printf("received event: %s for person %d", event.EventType, event.PersonID)

		err = c.processor.Handle(ctx, event)
		if err != nil {
			log.Printf("handle event error: %v", err)
		}
	}
}
