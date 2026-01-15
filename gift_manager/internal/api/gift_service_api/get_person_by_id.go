package giftserviceapi

import (
	"context"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
	pbmodels "gift_manager/internal/pb/models"
)

func (s *GiftServiceAPI) GetPersonByID(ctx context.Context, req *giftmanagerpb.GetPersonByIDRequest) (*giftmanagerpb.GetPersonByIDResponse, error) {
	person, err := s.giftService.GetPersonByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &giftmanagerpb.GetPersonByIDResponse{
		Person: &pbmodels.Person{
			Id:       uint64(person.ID),
			Name:     person.Name,
			Birthday: person.Birthday,
			Category: string(person.Category),
		},
	}, nil
}
