package bootstrap

import (
	"fmt"

	"notification_service/config"
	"notification_service/internal/consumer/personeventconsumer"
	"notification_service/internal/services/processors/personeventprocessor"
)

func InitPersonEventConsumer(cfg *config.Config, processor *personeventprocessor.PersonEventProcessor) *personeventconsumer.PersonEventConsumer {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return personeventconsumer.NewPersonEventConsumer(processor, brokers, cfg.Kafka.Topic)
}
