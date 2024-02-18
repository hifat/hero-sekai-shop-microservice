package inventory

type (
	Inventory struct {
		Id       string `json:"id" bson:"id"`
		PlayerId string `json:"player_id" bson:"player_id"`
		ItemId   string `json:"item_id" bson:"item_id"`
	}
)
