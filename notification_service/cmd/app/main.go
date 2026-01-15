package main

import (
	"fmt"
	"os"

	"notification_service/config"
	"notification_service/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("config error: %v", err))
	}

	storage := bootstrap.InitRedisStorage(cfg)
	service := bootstrap.InitNotificationService(storage, cfg)
	processor := bootstrap.InitPersonEventProcessor(service)
	consumer := bootstrap.InitPersonEventConsumer(cfg, processor)
	api := bootstrap.InitNotificationServiceAPI(service)

	bootstrap.AppRun(api, consumer)
}
