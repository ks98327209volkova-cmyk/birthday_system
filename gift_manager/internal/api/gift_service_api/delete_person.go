package giftserviceapi

import (
	"context"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

func (s *GiftServiceAPI) DeletePerson(ctx context.Context, req *giftmanagerpb.DeletePersonRequest) (*giftmanagerpb.DeletePersonResponse, error) {
	err := s.giftService.DeletePerson(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &giftmanagerpb.DeletePersonResponse{}, nil
}
