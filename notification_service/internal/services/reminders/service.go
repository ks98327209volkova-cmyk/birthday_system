package reminders

import (
	"context"
	"fmt"
	"log"
	"time"

	"notification_service/internal/models"
)

// интерфейс хранилища
type ReminderStorage interface {
	SaveReminder(ctx context.Context, reminder *models.Reminder) error
	GetAllReminders(ctx context.Context) ([]*models.Reminder, error)
	DeleteAllReminders(ctx context.Context) error
	DeleteRemindersByPerson(ctx context.Context, personID int64) error
}

type Service struct {
	storage      ReminderStorage
	reminderDays []int  // [30, 14, 7]
	currentMonth string // формат "2024-04"
}

func NewService(storage ReminderStorage, reminderDays []int) *Service {
	return &Service{
		storage:      storage,
		reminderDays: reminderDays,
		currentMonth: "",
	}
}

// основной метод обработки событий из Kafka
func (s *Service) HandlePersonEvent(ctx context.Context, event *models.PersonEvent) error {
	if err := validatePersonEvent(event); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	log.Printf("Processing %s for person %d", event.EventType, event.PersonID)

	switch event.EventType {
	case models.EventTypePersonDeleted:
		return s.handlePersonDeleted(ctx, event)

	case models.EventTypeMonthlyReminder:
		return s.handleMonthlyReminder(ctx, event)

	default:
		return s.handlePersonUpsert(ctx, event)
	}
}

func (s *Service) GetAllReminders(ctx context.Context) ([]*models.Reminder, error) {
	return s.storage.GetAllReminders(ctx)
}

func (s *Service) handlePersonDeleted(ctx context.Context, event *models.PersonEvent) error {
	return s.storage.DeleteRemindersByPerson(ctx, event.PersonID)
}

func (s *Service) handleMonthlyReminder(ctx context.Context, event *models.PersonEvent) error {
	if err := s.maybeClearRedisForNewMonth(ctx); err != nil {
		return fmt.Errorf("clear redis for new month: %w", err)
	}

	return s.calculateAndSaveReminders(ctx, event)
}

func (s *Service) handlePersonUpsert(ctx context.Context, event *models.PersonEvent) error {
	if err := s.storage.DeleteRemindersByPerson(ctx, event.PersonID); err != nil {
		return fmt.Errorf("delete old reminders: %w", err)
	}

	return s.calculateAndSaveReminders(ctx, event)
}

func (s *Service) maybeClearRedisForNewMonth(ctx context.Context) error {
	currentMonth := time.Now().Format("2006-01")

	if s.currentMonth != currentMonth {
		log.Printf("New month detected: %s, clearing Redis", currentMonth)

		if err := s.storage.DeleteAllReminders(ctx); err != nil {
			return err
		}

		s.currentMonth = currentMonth
	}

	return nil
}

func (s *Service) calculateAndSaveReminders(ctx context.Context, event *models.PersonEvent) error {
	reminders := s.calculateRemindersForCurrentMonth(event)

	for _, reminder := range reminders {
		if err := s.storage.SaveReminder(ctx, reminder); err != nil {
			return fmt.Errorf("save reminder: %w", err)
		}
	}

	if len(reminders) > 0 {
		log.Printf("Saved %d reminders for person %d", len(reminders), event.PersonID)
	}

	return nil
}

func validatePersonEvent(event *models.PersonEvent) error {
	if event.PersonID <= 0 {
		return fmt.Errorf("person_id must be positive")
	}
	if event.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// проверяем что birthday это валидная дата
	if event.Birthday == "" {
		return fmt.Errorf("birthday is required")
	}
	_, err := time.Parse("2006-01-02", event.Birthday)
	if err != nil {
		return fmt.Errorf("invalid birthday format, use YYYY-MM-DD")
	}
	return nil
}
