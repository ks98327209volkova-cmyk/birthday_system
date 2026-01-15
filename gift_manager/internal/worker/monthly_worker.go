package worker

import (
	"context"
	"log"
	"time"

	"gift_manager/internal/models"
)

// ServiceInterface - то что нужно воркеру от сервиса
type ServiceInterface interface {
	GetAllPersons(ctx context.Context) ([]*models.Person, error)
	SendPersonEvent(ctx context.Context, person *models.Person, eventType models.EventType) error
}

// MonthlyWorker отправляет ежемесячные напоминания
type MonthlyWorker struct {
	service ServiceInterface
	done    chan struct{}
}

// NewMonthlyWorker создает воркера
func NewMonthlyWorker(service ServiceInterface) *MonthlyWorker {
	return &MonthlyWorker{
		service: service,
		done:    make(chan struct{}),
	}
}

// Start запускает воркер
func (w *MonthlyWorker) Start() {
	go w.run()
	log.Println("monthly worker started")
}

// Stop останавливает воркер
func (w *MonthlyWorker) Stop() {
	close(w.done)
	log.Println("monthly worker stopped")
}

// run основной цикл
func (w *MonthlyWorker) run() {
	// таймер проверяет каждый день
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.checkDateAndSend()
		case <-w.done:
			return
		}
	}
}

// checkDateAndSend проверяет дату и отправляет если нужно
func (w *MonthlyWorker) checkDateAndSend() {
	now := time.Now()

	// если сегодня 1 число месяца
	if now.Day() == 1 {
		log.Printf("first day of month %s, sending reminders",
			now.Format("2006-01"))
		w.sendMonthlyReminders()
	}
}

// sendMonthlyReminders отправляет напоминания
func (w *MonthlyWorker) sendMonthlyReminders() {
	ctx := context.Background()

	// получаем всех людей
	persons, err := w.service.GetAllPersons(ctx)
	if err != nil {
		log.Printf("error getting persons: %v", err)
		return
	}

	if len(persons) == 0 {
		log.Println("no persons found")
		return
	}

	log.Printf("sending %d persons to kafka", len(persons))

	// отправляем событие для каждого человека
	for _, person := range persons {
		err := w.service.SendPersonEvent(ctx, person, models.EventTypeMonthlyReminder)
		if err != nil {
			log.Printf("error for person %d %s: %v",
				person.ID, person.Name, err)
		}
	}

	log.Println("monthly reminders sent")
}
