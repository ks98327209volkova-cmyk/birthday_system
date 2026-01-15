package giftserviceapi

import (
	"context"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
	pbmodels "gift_manager/internal/pb/models"
)

func (s *GiftServiceAPI) GetAllPersons(ctx context.Context, req *giftmanagerpb.GetAllPersonsRequest) (*giftmanagerpb.GetAllPersonsResponse, error) {
	persons, err := s.giftService.GetAllPersons(ctx)
	if err != nil {
		return nil, err
	}

	var protoPersons []*pbmodels.Person
	for _, p := range persons {
		protoPersons = append(protoPersons, &pbmodels.Person{
			Id:       uint64(p.ID),
			Name:     p.Name,
			Birthday: p.Birthday,
			Category: string(p.Category),
		})
	}

	return &giftmanagerpb.GetAllPersonsResponse{
		Persons: protoPersons,
	}, nil
}
