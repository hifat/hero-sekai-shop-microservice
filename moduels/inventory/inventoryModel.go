package inventory

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"
)

type (
	UpdateInventoryReq struct {
		PlayerId string `json:"player_id" validate:"required,max=64"`
		ItemId   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInInventory struct {
		InventoryId string `json:"inventory_id"`
		*item.ItemShowCase
	}

	PlayerInventory struct {
		PlayerId string `json:"player_id"`
		*model.PaginateRes
	}
)
