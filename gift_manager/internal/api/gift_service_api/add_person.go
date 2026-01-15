package giftserviceapi

import (
	"context"
	mymodels "gift_manager/internal/models"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
	pbmodels "gift_manager/internal/pb/models"
)

func (s *GiftServiceAPI) AddPerson(ctx context.Context, req *giftmanagerpb.AddPersonRequest) (*giftmanagerpb.AddPersonResponse, error) {
	person := &mymodels.Person{
		Name:     req.Person.Name,
		Birthday: req.Person.Birthday,
		Category: mymodels.Category(req.Person.Category),
	}

	result, err := s.giftService.AddPerson(ctx, person)
	if err != nil {
		return nil, err
	}

	return &giftmanagerpb.AddPersonResponse{
		Person: &pbmodels.Person{
			Id:       uint64(result.ID),
			Name:     result.Name,
			Birthday: result.Birthday,
			Category: string(result.Category),
		},
	}, nil
}
