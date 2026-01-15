package giftserviceapi

import (
	"context"
	mymodels "gift_manager/internal/models"
	giftmanagerpb "gift_manager/internal/pb/gift_manager_api"
)

type giftService interface {
	AddPerson(ctx context.Context, person *mymodels.Person) (*mymodels.Person, error)
	GetPersonByID(ctx context.Context, id uint64) (*mymodels.Person, error)
	UpdatePerson(ctx context.Context, person *mymodels.Person) error
	DeletePerson(ctx context.Context, id uint64) error
	GetAllPersons(ctx context.Context) ([]*mymodels.Person, error)
	AddGiftIdea(ctx context.Context, idea *mymodels.GiftIdea) error
	GetGiftIdeasByPersonID(ctx context.Context, personID uint64) ([]*mymodels.GiftIdea, error)
	DeleteGiftIdea(ctx context.Context, personID, ideaID uint64) error
	AddGiftHistory(ctx context.Context, history *mymodels.GiftHistory) error
	GetGiftHistoryByPersonID(ctx context.Context, personID uint64) ([]*mymodels.GiftHistory, error)
}

type GiftServiceAPI struct {
	giftmanagerpb.UnimplementedGiftManagerServiceServer
	giftService giftService
}

func NewGiftServiceAPI(giftService giftService) *GiftServiceAPI {
	return &GiftServiceAPI{
		giftService: giftService,
	}
}
