package models

import (
	"fmt"
	"time"
)

// PersonEvent - событие из Kafka от GiftManager
type PersonEvent struct {
	PersonID  int64
	Name      string
	Birthday  string
	EventType EventType
	Timestamp time.Time
}

// EventType - тип события
type EventType string

const (
	EventTypePersonAdded     EventType = "person_added"
	EventTypePersonUpdated   EventType = "person_updated"
	EventTypePersonDeleted   EventType = "person_deleted"
	EventTypeMonthlyReminder EventType = "monthly_reminder"
)

// Reminder - напоминание о дне рождения
type Reminder struct {
	ID           string
	PersonID     int64
	PersonName   string
	Birthday     string
	ReminderDate string
	DaysUntil    int
	CreatedAt    time.Time
}

// Key - генерирует ключ для Redis
// Формат: reminders:YYYY-MM:personID:days
func (r *Reminder) Key() string {
	return fmt.Sprintf("reminders:%s:%d:%d",
		r.ReminderDate[:7], // первые 7 символов "YYYY-MM"
		r.PersonID,
		r.DaysUntil)
}
