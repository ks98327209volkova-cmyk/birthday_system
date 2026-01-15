package giftserviceapi

import (
    "context"
    pbmodels "gift_manager/internal/pb/models"
    giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

func (s *GiftServiceAPI) GetGiftIdeasByPersonID(ctx context.Context, req *giftmanagerpb.GetGiftIdeasByPersonIDRequest) (*giftmanagerpb.GetGiftIdeasByPersonIDResponse, error) {
    ideas, err := s.giftService.GetGiftIdeasByPersonID(ctx, req.PersonId)
    if err != nil {
        return nil, err
    }
    
    var protoIdeas []*pbmodels.GiftIdea
    for _, idea := range ideas {
        protoIdeas = append(protoIdeas, &pbmodels.GiftIdea{
            Id:       uint64(idea.ID),
            PersonId: uint64(idea.PersonID),
            Idea:     idea.Idea,
            Notes:    idea.Notes,
        })
    }
    
    return &giftmanagerpb.GetGiftIdeasByPersonIDResponse{
        GiftIdeas: protoIdeas,
    }, nil
}