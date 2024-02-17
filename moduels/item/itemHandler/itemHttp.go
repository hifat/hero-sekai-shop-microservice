package itemHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemUsecase"
)

type (
	httpItemHttp struct {
		cfg         *config.Config
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemHttp(cfg *config.Config, itemUsecase itemUsecase.IItemUsecase) *httpItemHttp {
	return &httpItemHttp{cfg, itemUsecase}
}
