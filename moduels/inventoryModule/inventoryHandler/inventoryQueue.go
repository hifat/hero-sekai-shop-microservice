package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryUsecase"
)

type (
	inventoryQueue struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryQueue(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *inventoryQueue {
	return &inventoryQueue{cfg, inventoryUsecase}
}
