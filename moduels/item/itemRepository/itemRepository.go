package itemRepository

import "go.mongodb.org/mongo-driver/mongo"

type (
	IItemRepository interface{}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItem(db *mongo.Client) IItemRepository {
	return &itemRepository{db}
}
