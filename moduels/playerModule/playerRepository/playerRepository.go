package playerRepository

import (
	"context"
	"errors"
	"fmt"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/appError"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IPlayerRepository interface {
		FirstByField(pctx context.Context, field string, expected any) (*playerModule.Player, error)
		ExistsByField(pctx context.Context, field string, expected any) (bool, error)
		Create(pctx context.Context, req playerModule.Player) (string, error)
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

func (r *playerRepository) FirstByField(pctx context.Context, field string, expected any) (*playerModule.Player, error) {
	db := r.dbConn()
	col := db.Collection("players")

	if field == "_id" {
		expected = utils.ConvertToObjectId(fmt.Sprintf("%v", expected))
	}

	var result playerModule.Player
	if err := col.FindOne(
		pctx,
		bson.M{field: expected},
		options.FindOne().SetProjection(bson.M{
			"_id":        1,
			"email":      1,
			"username":   1,
			"created_at": 1,
			"updated_at": 1,
		}),
	).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appError.ErrRecordNotFound
		}
	}

	return &result, nil
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

func (r *playerRepository) Create(pctx context.Context, req playerModule.Player) (string, error) {
	col := r.dbConn().Collection("players")

	req.CreatedAt = utils.TimeNow()
	req.UpdatedAt = utils.TimeNow()

	playerId, err := col.InsertOne(pctx, req)
	if err != nil {
		return "", err
	}

	return playerId.InsertedID.(primitive.ObjectID).Hex(), nil
}
