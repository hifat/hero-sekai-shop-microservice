package itemHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemUsecase"
)

type (
	itemGrpc struct {
		cfg         *config.Config
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemGrpc(cfg *config.Config, itemUsecase itemUsecase.IItemUsecase) *itemGrpc {
	return &itemGrpc{cfg, itemUsecase}
}
