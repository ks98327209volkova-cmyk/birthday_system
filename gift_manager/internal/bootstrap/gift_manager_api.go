package bootstrap

import (
	giftserviceapi "gift_manager/internal/api/gift_service_api"
	giftservice "gift_manager/internal/services/giftService"
)

func InitGiftServiceAPI(giftService *giftservice.Service) *giftserviceapi.GiftServiceAPI {
	return giftserviceapi.NewGiftServiceAPI(giftService)
}
