package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryUsecase"
)

type (
	httpInventoryGrpc struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryGrpc(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *httpInventoryGrpc {
	return &httpInventoryGrpc{cfg, inventoryUsecase}
}
