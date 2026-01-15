package giftservice

import (
	"context"
	"testing"

	"gift_manager/internal/models"
	"gift_manager/internal/services/giftService/mocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type GiftHistoryServiceSuite struct {
	suite.Suite
	ctx             context.Context
	personRepo      *mocks.PersonRepository
	giftIdeaRepo    *mocks.GiftIdeaRepository
	giftHistoryRepo *mocks.GiftHistoryRepository
	eventPublisher  *mocks.EventPublisher
	service         *Service
}

func (s *GiftHistoryServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.personRepo = mocks.NewPersonRepository(s.T())
	s.giftIdeaRepo = mocks.NewGiftIdeaRepository(s.T())
	s.giftHistoryRepo = mocks.NewGiftHistoryRepository(s.T())
	s.eventPublisher = mocks.NewEventPublisher(s.T())

	s.service = NewService(
		s.personRepo,
		s.giftIdeaRepo,
		s.giftHistoryRepo,
		s.eventPublisher,
	)
}

func (s *GiftHistoryServiceSuite) TestAddGiftHistory_Success() {
	history := &models.GiftHistory{
		PersonID: 123,
		Gift:     "Сертификат в книжный магазин",
		Year:     2023,
	}

	person := &models.Person{
		ID:       123,
		Name:     "Слепов Павел",
		Birthday: "2004-09-24",
		Category: models.CategoryFamily,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(123)).Return(person, nil)
	s.giftHistoryRepo.On("AddGiftHistory", s.ctx, history).Return(nil)

	err := s.service.AddGiftHistory(s.ctx, history)

	assert.NoError(s.T(), err)
}

func (s *GiftHistoryServiceSuite) TestAddGiftHistory_ValidationError() {
	testCases := []struct {
		name    string
		history *models.GiftHistory
	}{
		{
			name: "empty gift",
			history: &models.GiftHistory{
				PersonID: 123,
				Gift:     "",
				Year:     2023,
			},
		},
		{
			name: "zero person ID",
			history: &models.GiftHistory{
				PersonID: 0,
				Gift:     "Книга",
				Year:     2023,
			},
		},
		{
			name: "zero year",
			history: &models.GiftHistory{
				PersonID: 123,
				Gift:     "Книга",
				Year:     0,
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.service.AddGiftHistory(s.ctx, tc.history)
			assert.Error(t, err)
		})
	}
}

func (s *GiftHistoryServiceSuite) TestAddGiftHistory_PersonNotFound() {
	history := &models.GiftHistory{
		PersonID: 123,
		Gift:     "Сертификат в книжный магазин",
		Year:     2023,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(123)).Return(nil, errors.New("person not found"))

	err := s.service.AddGiftHistory(s.ctx, history)

	assert.Error(s.T(), err)
}

func (s *GiftHistoryServiceSuite) TestAddGiftHistory_RepositoryError() {
	history := &models.GiftHistory{
		PersonID: 123,
		Gift:     "Сертификат в книжный магазин",
		Year:     2023,
	}

	person := &models.Person{
		ID:       123,
		Name:     "Слепов Павел",
		Birthday: "2004-09-24",
		Category: models.CategoryFamily,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(123)).Return(person, nil)
	s.giftHistoryRepo.On("AddGiftHistory", s.ctx, history).Return(errors.New("database error"))

	err := s.service.AddGiftHistory(s.ctx, history)

	assert.Error(s.T(), err)
}

func (s *GiftHistoryServiceSuite) TestGetGiftHistoryByPersonID_Success() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	expectedHistory := []*models.GiftHistory{
		{
			ID:       1,
			PersonID: int64(personID),
			Gift:     "Духи",
			Year:     2023,
		},
		{
			ID:       2,
			PersonID: int64(personID),
			Gift:     "Шарф",
			Year:     2022,
		},
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftHistoryRepo.On("GetGiftHistoryByPersonID", s.ctx, int64(personID)).Return(expectedHistory, nil)

	result, err := s.service.GetGiftHistoryByPersonID(s.ctx, personID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedHistory, result)
}

func (s *GiftHistoryServiceSuite) TestGetGiftHistoryByPersonID_ZeroID() {
	personID := uint64(0)

	result, err := s.service.GetGiftHistoryByPersonID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *GiftHistoryServiceSuite) TestGetGiftHistoryByPersonID_PersonNotFound() {
	personID := uint64(999)
	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(nil, errors.New("person not found"))

	result, err := s.service.GetGiftHistoryByPersonID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *GiftHistoryServiceSuite) TestGetGiftHistoryByPersonID_RepositoryError() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftHistoryRepo.On("GetGiftHistoryByPersonID", s.ctx, int64(personID)).Return(nil, errors.New("database error"))

	result, err := s.service.GetGiftHistoryByPersonID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *GiftHistoryServiceSuite) TestGetGiftHistoryByPersonID_EmptyList() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftHistoryRepo.On("GetGiftHistoryByPersonID", s.ctx, int64(personID)).Return([]*models.GiftHistory{}, nil)

	result, err := s.service.GetGiftHistoryByPersonID(s.ctx, personID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
}

func TestGiftHistoryServiceSuite(t *testing.T) {
	suite.Run(t, new(GiftHistoryServiceSuite))
}
