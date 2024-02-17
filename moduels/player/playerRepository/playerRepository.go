package playerRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IPlayerRepository interface{}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayer(db *mongo.Client) IPlayerRepository {
	return &playerRepository{db}
}

func (r *playerRepository) dbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("player_db")
}
