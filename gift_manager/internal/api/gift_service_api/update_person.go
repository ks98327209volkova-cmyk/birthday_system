package giftserviceapi

import (
	"context"
	mymodels "gift_manager/internal/models"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

func (s *GiftServiceAPI) UpdatePerson(ctx context.Context, req *giftmanagerpb.UpdatePersonRequest) (*giftmanagerpb.UpdatePersonResponse, error) {
	person := &mymodels.Person{
		ID:       int64(req.Person.Id),
		Name:     req.Person.Name,
		Birthday: req.Person.Birthday,
		Category: mymodels.Category(req.Person.Category),
	}

	err := s.giftService.UpdatePerson(ctx, person)
	if err != nil {
		return nil, err
	}

	return &giftmanagerpb.UpdatePersonResponse{}, nil
}
