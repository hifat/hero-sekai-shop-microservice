package inventoryUsecase

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryRepository"

type (
	IInventoryUsecase interface{}

	inventoryUsecase struct {
		inventoryRepo inventoryRepository.IInventoryRepository
	}
)

func NewInventory(inventoryRepo inventoryRepository.IInventoryRepository) IInventoryUsecase {
	return &inventoryUsecase{inventoryRepo}
}
