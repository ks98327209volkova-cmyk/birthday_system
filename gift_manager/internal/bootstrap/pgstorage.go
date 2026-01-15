package bootstrap

import (
	"fmt"
	"gift_manager/config"
	"gift_manager/internal/storage/pgstorage"
	"log"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGstorage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)

	storage, err := pgstorage.NewPGStorage(connectionString, cfg.Database.BucketQuantity)
	if err != nil {
		log.Panicf("ошибка инициализации БД: %v", err)
	}
	return storage
}
