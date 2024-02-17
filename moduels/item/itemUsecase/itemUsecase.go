package itemUsecase

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemRepository"
)

type (
	IItemUsecase interface{}

	itemUsecase struct {
		itemRepo itemRepository.IItemRepository
	}
)

func NewItem(itemRepo itemRepository.IItemRepository) IItemUsecase {
	return &itemUsecase{itemRepo}
}
