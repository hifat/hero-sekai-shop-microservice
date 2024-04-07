package inventoryRepository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IInventoryRepository interface{}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventory(db *mongo.Client) IInventoryRepository {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) dbConn() *mongo.Database {
	return r.db.Database("inventory_db")
}
