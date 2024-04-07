package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryUsecase"
)

type (
	inventoryHttp struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryHttp(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *inventoryHttp {
	return &inventoryHttp{cfg, inventoryUsecase}
}
