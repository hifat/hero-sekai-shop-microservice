package itemHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemUsecase"
)

type (
	itemHttp struct {
		cfg         *config.Config
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemHttp(cfg *config.Config, itemUsecase itemUsecase.IItemUsecase) *itemHttp {
	return &itemHttp{cfg, itemUsecase}
}
