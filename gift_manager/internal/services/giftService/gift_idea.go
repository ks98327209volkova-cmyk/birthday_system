package giftservice

import (
	"context"

	"gift_manager/internal/models"

	"github.com/pkg/errors"
)

func (s *Service) AddGiftIdea(ctx context.Context, idea *models.GiftIdea) error {
	if err := ValidateGiftIdea(idea); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	_, err := s.personRepo.GetPersonByID(ctx, idea.PersonID)
	if err != nil {
		return errors.Wrap(err, "person not found")
	}

	return s.giftIdeaRepo.AddGiftIdea(ctx, idea)
}

func (s *Service) GetGiftIdeasByPersonID(ctx context.Context, personID uint64) ([]*models.GiftIdea, error) {
	if personID == 0 {
		return nil, errors.New("invalid person ID")
	}

	_, err := s.personRepo.GetPersonByID(ctx, int64(personID))
	if err != nil {
		return nil, errors.Wrap(err, "person not found")
	}

	return s.giftIdeaRepo.GetGiftIdeasByPersonID(ctx, int64(personID))
}

func (s *Service) DeleteGiftIdea(ctx context.Context, personID, ideaID uint64) error {
	if personID == 0 || ideaID == 0 {
		return errors.New("invalid IDs")
	}

	_, err := s.personRepo.GetPersonByID(ctx, int64(personID))
	if err != nil {
		return errors.Wrap(err, "person not found")
	}

	return s.giftIdeaRepo.DeleteGiftIdea(ctx, int64(personID), int64(ideaID))
}
