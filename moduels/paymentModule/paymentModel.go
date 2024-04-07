package payment

type (
	ItemServiceReqDatum struct {
		ItemId string  `json:"item_id" validate:"required,max=64"`
		Price  float64 `json:"price"`
	}

	ItemServiceReq struct {
		Items []*ItemServiceReqDatum `json:"items" validate:"required"`
	}
)
