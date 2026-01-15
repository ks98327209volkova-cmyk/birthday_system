package bootstrap

import (
	"fmt"
	"gift_manager/config"
	"gift_manager/internal/producer/kafkaproducer"
)

func InitKafkaProducer(cfg *config.Config) *kafkaproducer.KafkaProducer {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return kafkaproducer.NewProducer(brokers, cfg.Kafka.Topic)
}
