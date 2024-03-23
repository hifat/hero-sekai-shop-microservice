package playerRepository

import (
	"context"
	"errors"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IPlayerRepository interface {
		FirstByField(pctx context.Context, field string, expected any) (bool, error)
		Create(pctx context.Context, req player.CreatePlayerReq) error
	}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayer(db *mongo.Client) IPlayerRepository {
	return &playerRepository{db}
}

func (r *playerRepository) dbConn() *mongo.Database {
	return r.db.Database("player_db")
}

func (r *playerRepository) FirstByField(pctx context.Context, field string, expected any) (bool, error) {
	db := r.dbConn()
	col := db.Collection("players")

	player := player.Player{}
	if err := col.FindOne(
		pctx,
		bson.M{
			field: expected,
		},
	).Decode(player); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, errors.New("record not found")
		}

		logger.Error(err.Error())
	}

	return false, nil
}

func (r *playerRepository) Create(pctx context.Context, req player.CreatePlayerReq) error {
	return nil
}
