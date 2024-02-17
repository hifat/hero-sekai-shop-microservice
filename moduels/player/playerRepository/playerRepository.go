package playerRepository

import "go.mongodb.org/mongo-driver/mongo"

type (
	IPlayerRepository interface{}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayer(db *mongo.Client) IPlayerRepository {
	return &playerRepository{db}
}
