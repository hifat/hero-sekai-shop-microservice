package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryUsecase"
)

type (
	httpInventoryQueue struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryQueue(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *httpInventoryQueue {
	return &httpInventoryQueue{cfg, inventoryUsecase}
}
