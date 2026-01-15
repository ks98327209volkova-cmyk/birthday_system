package bootstrap

import (
	"gift_manager/internal/producer/kafkaproducer"
	giftservice "gift_manager/internal/services/giftService"
	"gift_manager/internal/storage/pgstorage"
)

func InitGiftService(storage *pgstorage.PGstorage, producer *kafkaproducer.KafkaProducer) *giftservice.Service {
	return giftservice.NewService(
		storage,  // PersonRepository
		storage,  // GiftIdeaRepository
		storage,  // GiftHistoryRepository
		producer, // EventPublisher
	)
}
