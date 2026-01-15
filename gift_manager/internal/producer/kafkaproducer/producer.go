package kafkaproducer

import (
	"context"
	"encoding/json"
	"gift_manager/internal/models"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishPersonEvent(ctx context.Context, event *models.PersonEvent) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: topic,
		},
	}
}

func (p *KafkaProducer) PublishPersonEvent(ctx context.Context, event *models.PersonEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Printf("Publishing event to Kafka: %s", event.EventType)

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: eventJSON,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
