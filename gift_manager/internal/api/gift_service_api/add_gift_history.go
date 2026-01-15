package giftserviceapi

import (
    "context"
    mymodels "gift_manager/internal/models"
    giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

func (s *GiftServiceAPI) AddGiftHistory(ctx context.Context, req *giftmanagerpb.AddGiftHistoryRequest) (*giftmanagerpb.AddGiftHistoryResponse, error) {
    history := &mymodels.GiftHistory{
        PersonID: int64(req.PersonId),
        Gift:     req.Gift,
        Year:     int32(req.Year),
    }
    
    err := s.giftService.AddGiftHistory(ctx, history)
    if err != nil {
        return nil, err
    }
    
    return &giftmanagerpb.AddGiftHistoryResponse{}, nil
}