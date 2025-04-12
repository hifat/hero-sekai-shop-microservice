package itemModule

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"

type (
	CreateItemReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		Damage   int32   `json:"damage" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
	}

	ItemShowCase struct {
		ItemId   string  `json:"item_id"`
		Title    string  `json:"title"`
		Price    float64 `json:"price"`
		Damage   int32   `json:"damage"`
		ImageUrl string  `json:"image_url"`
	}

	ItemSearchReq struct {
		Title string `query:"title" validate:"max=64"`
		model.PaginateReq
	}

	ItemUpdateReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
		Damage   int     `json:"damage" validate:"required"`
	}

	EnableOrDisableItemReq struct {
		UsageStatus bool `json:"usage_status"`
	}
)
