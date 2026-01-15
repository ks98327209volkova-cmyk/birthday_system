package models

import "time"

// Category - категория человека
type Category string

// константы для категорий
const (
	CategoryFamily     Category = "family"
	CategoryFriends    Category = "friends"
	CategoryColleagues Category = "colleagues"
)

// Person - человек, у которого день рождения
type Person struct {
	ID       int64
	Name     string
	Birthday string
	Category Category
}

// GiftIdea - идея подарка для конкретного человека
type GiftIdea struct {
	ID       int64
	PersonID int64
	Idea     string
	Notes    string
}

// GiftHistory - история уже подаренных подарков
type GiftHistory struct {
	ID       int64
	PersonID int64
	Gift     string
	Year     int32
}

type EventType string

const (
	EventTypePersonAdded     EventType = "person_added"
	EventTypePersonUpdated   EventType = "person_updated"
	EventTypePersonDeleted   EventType = "person_deleted"
	EventTypeMonthlyReminder EventType = "monthly_reminder"
)

type PersonEvent struct {
	PersonID  int64
	Name      string
	Birthday  string
	EventType EventType
	Timestamp time.Time
}
