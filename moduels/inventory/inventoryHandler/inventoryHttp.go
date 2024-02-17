package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryUsecase"
)

type (
	httpInventoryHttp struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryHttp(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) httpInventoryHttp {
	return httpInventoryHttp{cfg, inventoryUsecase}
}
