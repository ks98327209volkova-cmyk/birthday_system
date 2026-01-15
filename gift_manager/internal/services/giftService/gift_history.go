package giftservice

import (
	"context"

	"gift_manager/internal/models"

	"github.com/pkg/errors"
)

func (s *Service) AddGiftHistory(ctx context.Context, history *models.GiftHistory) error {
	if err := ValidateGiftHistory(history); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	_, err := s.personRepo.GetPersonByID(ctx, history.PersonID)
	if err != nil {
		return errors.Wrap(err, "person not found")
	}

	return s.giftHistoryRepo.AddGiftHistory(ctx, history)
}

func (s *Service) GetGiftHistoryByPersonID(ctx context.Context, personID uint64) ([]*models.GiftHistory, error) {
	if personID == 0 {
		return nil, errors.New("invalid person ID")
	}

	_, err := s.personRepo.GetPersonByID(ctx, int64(personID))
	if err != nil {
		return nil, errors.Wrap(err, "person not found")
	}

	return s.giftHistoryRepo.GetGiftHistoryByPersonID(ctx, int64(personID))
}
