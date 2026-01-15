package reminders

import (
	"fmt"
	"time"

	"notification_service/internal/models"
)

// рассчитывает напоминания на текущий месяц
func (s *Service) calculateRemindersForCurrentMonth(event *models.PersonEvent) []*models.Reminder {
	var reminders []*models.Reminder
	now := time.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	// Парсим birthday из строки
	birthday, err := time.Parse("2006-01-02", event.Birthday)
	if err != nil {
		// Если ошибка - возвращаем пустой список
		return reminders
	}

	birthdayThisYear := time.Date(
		currentYear,
		birthday.Month(),
		birthday.Day(),
		0, 0, 0, 0, time.UTC,
	)

	if birthdayThisYear.Before(now) {
		birthdayThisYear = birthdayThisYear.AddDate(1, 0, 0)
	}

	for _, daysBefore := range s.reminderDays {
		reminderDate := birthdayThisYear.AddDate(0, 0, -daysBefore)

		if reminderDate.Month() == currentMonth &&
			reminderDate.Year() == currentYear &&
			!reminderDate.Before(now) {

			reminderID := fmt.Sprintf("%d_%d_%s",
				event.PersonID,
				daysBefore,
				reminderDate.Format("20060102"))

			reminder := &models.Reminder{
				ID:           reminderID,
				PersonID:     event.PersonID,
				PersonName:   event.Name,
				Birthday:     event.Birthday,
				ReminderDate: reminderDate.Format("2006-01-02"),
				DaysUntil:    daysBefore,
				CreatedAt:    time.Now(),
			}

			reminders = append(reminders, reminder)
		}
	}

	return reminders
}
