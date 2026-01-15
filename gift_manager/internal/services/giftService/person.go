package giftservice

import (
	"context"
	"time"

	"gift_manager/internal/models"

	"github.com/pkg/errors"
)

func (s *Service) AddPerson(ctx context.Context, person *models.Person) (*models.Person, error) {
	if err := ValidatePerson(person); err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	id, err := s.personRepo.AddPerson(ctx, person)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save person")
	}

	personWithID := *person
	personWithID.ID = id

	event := &models.PersonEvent{
		PersonID:  id,
		Name:      person.Name,
		Birthday:  person.Birthday,
		EventType: models.EventTypePersonAdded,
		Timestamp: time.Now(),
	}

	s.eventPublisher.PublishPersonEvent(ctx, event)

	return &personWithID, nil
}

func (s *Service) GetPersonByID(ctx context.Context, id uint64) (*models.Person, error) {
	if id == 0 {
		return nil, errors.New("invalid person ID")
	}

	return s.personRepo.GetPersonByID(ctx, int64(id))
}

func (s *Service) UpdatePerson(ctx context.Context, person *models.Person) error {
	if person.ID <= 0 {
		return errors.New("invalid person ID")
	}

	if err := ValidatePerson(person); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	if err := s.personRepo.UpdatePerson(ctx, person); err != nil {
		return errors.Wrap(err, "failed to update person")
	}

	event := &models.PersonEvent{
		PersonID:  person.ID,
		Name:      person.Name,
		Birthday:  person.Birthday,
		EventType: models.EventTypePersonUpdated,
		Timestamp: time.Now(),
	}

	return s.eventPublisher.PublishPersonEvent(ctx, event)
}

func (s *Service) DeletePerson(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("invalid person ID")
	}

	person, err := s.personRepo.GetPersonByID(ctx, int64(id))
	if err != nil {
		return errors.Wrap(err, "person not found")
	}

	if err := s.personRepo.DeletePerson(ctx, int64(id)); err != nil {
		return errors.Wrap(err, "failed to delete person")
	}

	event := &models.PersonEvent{
		PersonID:  int64(id),
		Name:      person.Name,
		Birthday:  person.Birthday,
		EventType: models.EventTypePersonDeleted,
		Timestamp: time.Now(),
	}

	return s.eventPublisher.PublishPersonEvent(ctx, event)
}

func (s *Service) GetAllPersons(ctx context.Context) ([]*models.Person, error) {
	return s.personRepo.GetAllPersons(ctx)
}