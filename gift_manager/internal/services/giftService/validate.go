package giftservice

import (
	"time"

	"gift_manager/internal/models"
)

func ValidatePerson(person *models.Person) error {
	if person.Name == "" {
		return ValidationError{Field: "name", Reason: "cannot be empty"}
	}

	if len(person.Name) > 255 {
		return ValidationError{Field: "name", Reason: "too long (max 255 characters)"}
	}

	if person.Birthday == "" {
		return ValidationError{Field: "birthday", Reason: "is required"}
	}

	birthday, err := time.Parse("2006-01-02", person.Birthday)
	if err != nil {
		return ValidationError{Field: "birthday", Reason: "invalid format, use YYYY-MM-DD"}
	}

	if birthday.After(time.Now()) {
		return ValidationError{Field: "birthday", Reason: "cannot be in the future"}
	}

	switch person.Category {
	case models.CategoryFamily, models.CategoryFriends, models.CategoryColleagues:
		// ok
	default:
		return ValidationError{Field: "category", Reason: "invalid value"}
	}

	return nil
}

func ValidateGiftIdea(idea *models.GiftIdea) error {
	if idea.PersonID <= 0 {
		return ValidationError{Field: "person_id", Reason: "must be positive"}
	}

	if idea.Idea == "" {
		return ValidationError{Field: "idea", Reason: "cannot be empty"}
	}

	if len(idea.Idea) > 1000 {
		return ValidationError{Field: "idea", Reason: "too long (max 1000 characters)"}
	}

	return nil
}

func ValidateGiftHistory(history *models.GiftHistory) error {
	if history.PersonID <= 0 {
		return ValidationError{Field: "person_id", Reason: "must be positive"}
	}

	if history.Gift == "" {
		return ValidationError{Field: "gift", Reason: "cannot be empty"}
	}

	if history.Year <= 0 {
		return ValidationError{Field: "year", Reason: "must be positive"}
	}

	if history.Year > int32(time.Now().Year()) {
		return ValidationError{Field: "year", Reason: "cannot be in the future"}
	}

	return nil
}

type ValidationError struct {
	Field  string
	Reason string
}

func (v ValidationError) Error() string {
	return v.Field + ": " + v.Reason
}
