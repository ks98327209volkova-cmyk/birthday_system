package giftservice

import (
	"context"
	"time"

	"gift_manager/internal/models"
)

// интерфейсы репозиториев storage слой
type PersonRepository interface {
	AddPerson(ctx context.Context, person *models.Person) (int64, error)
	GetPersonByID(ctx context.Context, id int64) (*models.Person, error)
	UpdatePerson(ctx context.Context, person *models.Person) error
	DeletePerson(ctx context.Context, id int64) error
	GetAllPersons(ctx context.Context) ([]*models.Person, error)
}

type GiftIdeaRepository interface {
	AddGiftIdea(ctx context.Context, idea *models.GiftIdea) error
	GetGiftIdeasByPersonID(ctx context.Context, personID int64) ([]*models.GiftIdea, error)
	DeleteGiftIdea(ctx context.Context, personID, ideaID int64) error
}

type GiftHistoryRepository interface {
	AddGiftHistory(ctx context.Context, history *models.GiftHistory) error
	GetGiftHistoryByPersonID(ctx context.Context, personID int64) ([]*models.GiftHistory, error)
}

// интерфейс для kafka producer
type EventPublisher interface {
	PublishPersonEvent(ctx context.Context, event *models.PersonEvent) error
}

// ServiceInterface для воркера
type ServiceInterface interface {
	GetAllPersons(ctx context.Context) ([]*models.Person, error)
	SendPersonEvent(ctx context.Context, person *models.Person, eventType models.EventType) error
}

type Service struct {
	personRepo      PersonRepository
	giftIdeaRepo    GiftIdeaRepository
	giftHistoryRepo GiftHistoryRepository
	eventPublisher  EventPublisher
}

func NewService(
	personRepo PersonRepository,
	giftIdeaRepo GiftIdeaRepository,
	giftHistoryRepo GiftHistoryRepository,
	eventPublisher EventPublisher,
) *Service {
	return &Service{
		personRepo:      personRepo,
		giftIdeaRepo:    giftIdeaRepo,
		giftHistoryRepo: giftHistoryRepo,
		eventPublisher:  eventPublisher,
	}
}

// SendPersonEvent отправляет событие в кафку
// используется воркером для ежемесячных напоминаний
func (s *Service) SendPersonEvent(ctx context.Context, person *models.Person, eventType models.EventType) error {
	event := &models.PersonEvent{
		PersonID:  person.ID,
		Name:      person.Name,
		Birthday:  person.Birthday,
		EventType: eventType,
		Timestamp: time.Now(),
	}
	return s.eventPublisher.PublishPersonEvent(ctx, event)
}
