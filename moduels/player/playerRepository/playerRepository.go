package playerRepository

import (
	"context"
	"errors"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/appError"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IPlayerRepository interface {
		FirstByField(pctx context.Context, field string, expected any) (*player.Player, error)
		ExistsByField(pctx context.Context, field string, expected any) (bool, error)
		Create(pctx context.Context, req player.Player) (string, error)
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

func (r *playerRepository) FirstByField(pctx context.Context, field string, expected any) (*player.Player, error) {
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
			return nil, appError.ErrRecordNotFound
		}
	}

	return &player, nil
}

func (r *playerRepository) ExistsByField(pctx context.Context, field string, expected any) (bool, error) {
	_, err := r.FirstByField(pctx, field, expected)
	if err != nil {
		if errors.Is(err, appError.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *playerRepository) Create(pctx context.Context, req player.Player) (string, error) {
	col := r.dbConn().Collection("player")

	req.CreatedAt = utils.TimeNow()
	req.UpdatedAt = utils.TimeNow()

	playerId, err := col.InsertOne(pctx, req)
	if err != nil {
		return "", err
	}

	return playerId.InsertedID.(primitive.ObjectID).Hex(), nil
}
