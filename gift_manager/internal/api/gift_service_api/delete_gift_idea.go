package giftserviceapi

import (
    "context"
    giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

func (s *GiftServiceAPI) DeleteGiftIdea(ctx context.Context, req *giftmanagerpb.DeleteGiftIdeaRequest) (*giftmanagerpb.DeleteGiftIdeaResponse, error) {
    err := s.giftService.DeleteGiftIdea(ctx, req.PersonId, req.IdeaId)
    if err != nil {
        return nil, err
    }
    
    return &giftmanagerpb.DeleteGiftIdeaResponse{}, nil
}