package itemHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemUsecase"
)

type (
	httpItemHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemHttp(cfg *config.Config, itemUsecase itemUsecase.IItemUsecase) httpItemHandler {
	return httpItemHandler{cfg, itemUsecase}
}
