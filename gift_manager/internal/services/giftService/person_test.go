package giftservice

import (
	"context"
	"testing"

	"gift_manager/internal/models"
	"gift_manager/internal/services/giftService/mocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PersonServiceSuite struct {
	suite.Suite
	ctx             context.Context
	personRepo      *mocks.PersonRepository
	giftIdeaRepo    *mocks.GiftIdeaRepository
	giftHistoryRepo *mocks.GiftHistoryRepository
	eventPublisher  *mocks.EventPublisher
	service         *Service
}

func (s *PersonServiceSuite) SetupTest() {
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

func (s *PersonServiceSuite) TestAddPerson_Success() {
	person := &models.Person{
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("AddPerson", s.ctx, person).Return(int64(123), nil)
	s.eventPublisher.On("PublishPersonEvent", s.ctx, mock.Anything).Return(nil)

	result, err := s.service.AddPerson(s.ctx, person)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(123), result.ID)
	assert.Equal(s.T(), "Волкова Ксения", result.Name)
}

func (s *PersonServiceSuite) TestAddPerson_ValidationError() {
	testCases := []struct {
		name   string
		person *models.Person
	}{
		{
			name: "empty name",
			person: &models.Person{
				Name:     "",
				Birthday: "2000-01-01",
				Category: models.CategoryFriends,
			},
		},
		{
			name: "future birthday",
			person: &models.Person{
				Name:     "Волкова Ксения",
				Birthday: "2100-01-01",
				Category: models.CategoryFriends,
			},
		},
		{
			name: "invalid category",
			person: &models.Person{
				Name:     "Волкова Ксения",
				Birthday: "2001-06-27",
				Category: "invalid-category",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			result, err := s.service.AddPerson(s.ctx, tc.person)

			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func (s *PersonServiceSuite) TestAddPerson_RepositoryError() {
	person := &models.Person{
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("AddPerson", s.ctx, person).Return(int64(0), errors.New("database error"))

	result, err := s.service.AddPerson(s.ctx, person)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *PersonServiceSuite) TestGetPersonByID_Success() {
	personID := uint64(123)
	expectedPerson := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(expectedPerson, nil)

	result, err := s.service.GetPersonByID(s.ctx, personID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedPerson, result)
}

func (s *PersonServiceSuite) TestGetPersonByID_ZeroID() {
	personID := uint64(0)

	result, err := s.service.GetPersonByID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *PersonServiceSuite) TestGetPersonByID_RepositoryError() {
	personID := uint64(123)
	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(nil, errors.New("person not found"))

	result, err := s.service.GetPersonByID(s.ctx, personID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *PersonServiceSuite) TestUpdatePerson_Success() {
	person := &models.Person{
		ID:       123,
		Name:     "Волкова Ксения (изменено)",
		Birthday: "2001-06-27",
		Category: models.CategoryFamily,
	}

	s.personRepo.On("UpdatePerson", s.ctx, person).Return(nil)
	s.eventPublisher.On("PublishPersonEvent", s.ctx, mock.Anything).Return(nil)

	err := s.service.UpdatePerson(s.ctx, person)

	assert.NoError(s.T(), err)
}

func (s *PersonServiceSuite) TestUpdatePerson_InvalidID() {
	testCases := []struct {
		name   string
		person *models.Person
	}{
		{
			name: "zero ID",
			person: &models.Person{
				ID:       0,
				Name:     "Волкова Ксения",
				Birthday: "2001-06-27",
				Category: models.CategoryFriends,
			},
		},
		{
			name: "negative ID",
			person: &models.Person{
				ID:       -1,
				Name:     "Волкова Ксения",
				Birthday: "2001-06-27",
				Category: models.CategoryFriends,
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.service.UpdatePerson(s.ctx, tc.person)

			assert.Error(t, err)
		})
	}
}

func (s *PersonServiceSuite) TestUpdatePerson_ValidationError() {
	person := &models.Person{
		ID:       123,
		Name:     "",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	err := s.service.UpdatePerson(s.ctx, person)

	assert.Error(s.T(), err)
}

func (s *PersonServiceSuite) TestUpdatePerson_RepositoryError() {
	person := &models.Person{
		ID:       123,
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("UpdatePerson", s.ctx, person).Return(errors.New("update failed"))

	err := s.service.UpdatePerson(s.ctx, person)

	assert.Error(s.T(), err)
}

func (s *PersonServiceSuite) TestDeletePerson_Success() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.personRepo.On("DeletePerson", s.ctx, int64(personID)).Return(nil)
	s.eventPublisher.On("PublishPersonEvent", s.ctx, mock.Anything).Return(nil)

	err := s.service.DeletePerson(s.ctx, personID)

	assert.NoError(s.T(), err)
}

func (s *PersonServiceSuite) TestDeletePerson_ZeroID() {
	personID := uint64(0)

	err := s.service.DeletePerson(s.ctx, personID)

	assert.Error(s.T(), err)
}

func (s *PersonServiceSuite) TestDeletePerson_PersonNotFound() {
	personID := uint64(999)
	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(nil, errors.New("not found"))

	err := s.service.DeletePerson(s.ctx, personID)

	assert.Error(s.T(), err)
}

func (s *PersonServiceSuite) TestDeletePerson_RepositoryDeleteError() {
	personID := uint64(123)
	person := &models.Person{
		ID:       int64(personID),
		Name:     "Волкова Ксения",
		Birthday: "2001-06-27",
		Category: models.CategoryFriends,
	}

	s.personRepo.On("GetPersonByID", s.ctx, int64(personID)).Return(person, nil)
	s.personRepo.On("DeletePerson", s.ctx, int64(personID)).Return(errors.New("delete failed"))

	err := s.service.DeletePerson(s.ctx, personID)

	assert.Error(s.T(), err)
}

func (s *PersonServiceSuite) TestGetAllPersons_Success() {
	expectedPersons := []*models.Person{
		{
			ID:       1,
			Name:     "Волкова Ксения",
			Birthday: "2001-06-27",
			Category: models.CategoryFriends,
		},
		{
			ID:       2,
			Name:     "Слепов Павел",
			Birthday: "2004-09-24",
			Category: models.CategoryFamily,
		},
	}

	s.personRepo.On("GetAllPersons", s.ctx).Return(expectedPersons, nil)

	result, err := s.service.GetAllPersons(s.ctx)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedPersons, result)
}

func (s *PersonServiceSuite) TestGetAllPersons_RepositoryError() {
	s.personRepo.On("GetAllPersons", s.ctx).Return(nil, errors.New("database error"))

	result, err := s.service.GetAllPersons(s.ctx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *PersonServiceSuite) TestGetAllPersons_EmptyList() {
	s.personRepo.On("GetAllPersons", s.ctx).Return([]*models.Person{}, nil)

	result, err := s.service.GetAllPersons(s.ctx)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
}

func (s *PersonServiceSuite) TestGetAllPersons_NilList() {
	s.personRepo.On("GetAllPersons", s.ctx).Return(nil, nil)

	result, err := s.service.GetAllPersons(s.ctx)

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

func TestPersonServiceSuite(t *testing.T) {
	suite.Run(t, new(PersonServiceSuite))
}
