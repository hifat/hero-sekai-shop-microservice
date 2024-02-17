package itemHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemUsecase"
)

type (
	httpItemGrpc struct {
		cfg         *config.Config
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemGrpc(cfg *config.Config, itemUsecase itemUsecase.IItemUsecase) httpItemGrpc {
	return httpItemGrpc{cfg, itemUsecase}
}
