package model

type (
	PaginateReq struct {
		Start string `query:"start" validate:"max=64"`
		Limit int64  `query:"limit" validate:"min=2,max=10"`
	}

	FirstPaginate struct {
		Href string `json:"href"`
	}

	NextPaginate struct {
		Href  string `json:"href"`
		Start string `json:"start"`
	}

	PaginateRes struct {
		Data  any           `json:"data"`
		Limit int64         `json:"limit"`
		Total int64         `json:"total"`
		First FirstPaginate `json:"first"`
		Next  NextPaginate  `json:"next"`
	}

	KafkaOffset struct {
		Offset int64 `json:"offset" bson:"offset"`
	}
)
