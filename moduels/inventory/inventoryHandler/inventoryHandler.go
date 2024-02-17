package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryUsecase"
)

type (
	httpInventoryHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryHttp(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) httpInventoryHandler {
	return httpInventoryHandler{cfg, inventoryUsecase}
}
