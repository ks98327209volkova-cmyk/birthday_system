package giftserviceapi

import (
	"context"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
	pbmodels "gift_manager/internal/pb/models"
)

func (s *GiftServiceAPI) GetGiftHistoryByPersonID(ctx context.Context, req *giftmanagerpb.GetGiftHistoryByPersonIDRequest) (*giftmanagerpb.GetGiftHistoryByPersonIDResponse, error) {
	history, err := s.giftService.GetGiftHistoryByPersonID(ctx, req.PersonId)
	if err != nil {
		return nil, err
	}

	var protoHistory []*pbmodels.GiftHistory
	for _, h := range history {
		protoHistory = append(protoHistory, &pbmodels.GiftHistory{
			Id:       uint64(h.ID),
			PersonId: uint64(h.PersonID),
			Gift:     h.Gift,
			Year:     uint32(h.Year),
		})
	}

	return &giftmanagerpb.GetGiftHistoryByPersonIDResponse{
		GiftHistory: protoHistory,
	}, nil
}
