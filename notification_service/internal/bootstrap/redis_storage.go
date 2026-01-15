package bootstrap

import (
	"log"

	"notification_service/config"
	"notification_service/internal/storage/redisstorage"
)

func InitRedisStorage(cfg *config.Config) *redisstorage.RedisStorage {
	storage, err := redisstorage.NewRedisStorage(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
		31,
	)
	if err != nil {
		log.Panicf("redis init error: %v", err)
	}
	return storage
}
