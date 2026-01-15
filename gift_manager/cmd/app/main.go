package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gift_manager/config"
	"gift_manager/internal/bootstrap"
	"gift_manager/internal/worker"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	producer := bootstrap.InitKafkaProducer(cfg)
	service := bootstrap.InitGiftService(storage, producer)
	api := bootstrap.InitGiftServiceAPI(service)

	// создаем воркер
	monthlyWorker := worker.NewMonthlyWorker(service)
	monthlyWorker.Start()
	defer monthlyWorker.Stop()

	// запускаем сервер
	go func() {
		bootstrap.AppRun(api, producer)
	}()

	// ждем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
}
