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

type GiftIdeaServiceSuite struct {
	suite.Suite
	ctx             context.Context
	personRepo      *mocks.PersonRepository
	giftIdeaRepo    *mocks.GiftIdeaRepository
	giftHistoryRepo *mocks.GiftHistoryRepository
	eventPublisher  *mocks.EventPublisher
	service         *Service
}

func (s *GiftIdeaServiceSuite) SetupTest() {
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

func (s *GiftIdeaServiceSuite) TestAddGiftIdea_Success() {
	idea := &models.GiftIdea{
		PersonID: 123,
		Idea:     "Плойка",
		Notes:    "Давно хотела",
	}

	person := &models.Person{
		ID:       123,
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(123)).Return(person, nil)
	s.giftIdeaRepo.On("AddGiftIdea", s.ctx, idea).Return(nil)

	err := s.service.AddGiftIdea(s.ctx, idea)

	assert.NoError(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestAddGiftIdea_ValidationError() {
	testCases := []struct {
		name string
		idea *models.GiftIdea
	}{
		{
			name: "empty idea",
			idea: &models.GiftIdea{
				PersonID: 123,
				Idea:     "",
				Notes:    "Примечание",
			},
		},
		{
			name: "zero person ID",
			idea: &models.GiftIdea{
				PersonID: 0,
				Idea:     "Подарок",
				Notes:    "Примечание",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.service.AddGiftIdea(s.ctx, tc.idea)
			assert.Error(t, err)
		})
	}
}

func (s *GiftIdeaServiceSuite) TestAddGiftIdea_PersonNotFound() {
	idea := &models.GiftIdea{
		PersonID: 123,
		Idea:     "Плойка",
		Notes:    "Давно хотела",
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(123)).Return(nil, errors.New("person not found"))

	err := s.service.AddGiftIdea(s.ctx, idea)

	assert.Error(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestAddGiftIdea_RepositoryError() {
	idea := &models.GiftIdea{
		PersonID: 123,
		Idea:     "Плойка",
		Notes:    "Давно хотела",
	}

	person := &models.Person{
		ID:       123,
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(123)).Return(person, nil)
	s.giftIdeaRepo.On("AddGiftIdea", s.ctx, idea).Return(errors.New("database error"))

	err := s.service.AddGiftIdea(s.ctx, idea)

	assert.Error(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestGetGiftIdeasByPersonID_Success() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Слепов Павел",
		Birthday: "2004-09-24",
		Category: models.CategoryFamily,
	}

	expectedIdeas := []*models.GiftIdea{
		{
			ID:       1,
			PersonID: int64(personID),
			Idea:     "3 тома книг Real-RPG",
			Notes:    "Хотел как только выйдут у автора",
		},
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftIdeaRepo.On("GetGiftIdeasByPersonID", s.ctx, int64(personID)).Return(expectedIdeas, nil)

	result, err := s.service.GetGiftIdeasByPersonID(s.ctx, personID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedIdeas, result)
}

func (s *GiftIdeaServiceSuite) TestGetGiftIdeasByPersonID_ZeroID() {
	personID := uint64(0)

	result, err := s.service.GetGiftIdeasByPersonID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *GiftIdeaServiceSuite) TestGetGiftIdeasByPersonID_PersonNotFound() {
	personID := uint64(999)
	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(nil, errors.New("person not found"))

	result, err := s.service.GetGiftIdeasByPersonID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *GiftIdeaServiceSuite) TestGetGiftIdeasByPersonID_RepositoryError() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Слепов Павел",
		Birthday: "2004-09-24",
		Category: models.CategoryFamily,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftIdeaRepo.On("GetGiftIdeasByPersonID", s.ctx, int64(personID)).Return(nil, errors.New("database error"))

	result, err := s.service.GetGiftIdeasByPersonID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *GiftIdeaServiceSuite) TestGetGiftIdeasByPersonID_EmptyList() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Слепов Павел",
		Birthday: "2004-09-24",
		Category: models.CategoryFamily,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftIdeaRepo.On("GetGiftIdeasByPersonID", s.ctx, int64(personID)).Return([]*models.GiftIdea{}, nil)

	result, err := s.service.GetGiftIdeasByPersonID(s.ctx, personID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
}

func (s *GiftIdeaServiceSuite) TestDeleteGiftIdea_Success() {
	personID := uint64(123)
	ideaID := uint64(456)

	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftIdeaRepo.On("DeleteGiftIdea", s.ctx, int64(personID), int64(ideaID)).Return(nil)

	err := s.service.DeleteGiftIdea(s.ctx, personID, ideaID)

	assert.NoError(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestDeleteGiftIdea_ZeroPersonID() {
	personID := uint64(0)
	ideaID := uint64(456)

	err := s.service.DeleteGiftIdea(s.ctx, personID, ideaID)

	assert.Error(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestDeleteGiftIdea_ZeroIdeaID() {
	personID := uint64(123)
	ideaID := uint64(0)

	err := s.service.DeleteGiftIdea(s.ctx, personID, ideaID)

	assert.Error(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestDeleteGiftIdea_PersonNotFound() {
	personID := uint64(123)
	ideaID := uint64(456)

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(nil, errors.New("person not found"))

	err := s.service.DeleteGiftIdea(s.ctx, personID, ideaID)

	assert.Error(s.T(), err)
}

func (s *GiftIdeaServiceSuite) TestDeleteGiftIdea_RepositoryError() {
	personID := uint64(123)
	ideaID := uint64(456)

	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.giftIdeaRepo.On("DeleteGiftIdea", s.ctx, int64(personID), int64(ideaID)).Return(errors.New("database error"))

	err := s.service.DeleteGiftIdea(s.ctx, personID, ideaID)

	assert.Error(s.T(), err)
}

func TestGiftIdeaServiceSuite(t *testing.T) {
	suite.Run(t, new(GiftIdeaServiceSuite))
}
