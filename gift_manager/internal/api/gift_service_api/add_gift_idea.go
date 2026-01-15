package giftserviceapi

import (
    "context"
    mymodels "gift_manager/internal/models"
    giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

func (s *GiftServiceAPI) AddGiftIdea(ctx context.Context, req *giftmanagerpb.AddGiftIdeaRequest) (*giftmanagerpb.AddGiftIdeaResponse, error) {
    idea := &mymodels.GiftIdea{
        PersonID: int64(req.PersonId),
        Idea:     req.Idea,
        Notes:    req.Notes,
    }
    
    err := s.giftService.AddGiftIdea(ctx, idea)
    if err != nil {
        return nil, err
    }
    
    return &giftmanagerpb.AddGiftIdeaResponse{}, nil
}