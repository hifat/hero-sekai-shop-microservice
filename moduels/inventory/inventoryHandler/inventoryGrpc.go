package inventoryHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryUsecase"
)

type (
	inventoryGrpc struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryGrpc(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *inventoryGrpc {
	return &inventoryGrpc{cfg, inventoryUsecase}
}

func (h *inventoryGrpc) Foo() {

}
