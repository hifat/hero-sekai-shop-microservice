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
	IPlayerTransactionRepository interface {
		FirstByField(pctx context.Context, field string, expected any) (*playerModule.PlayerTransaction, error)
		ExistsByField(pctx context.Context, field string, expected any) (bool, error)
		Create(pctx context.Context, req playerModule.PlayerTransaction) (string, error)
		GetSavingAccount(pctx context.Context, playerId string) (*playerModule.PlayerSavingAccount, error)
	}

	playerTransactionRepository struct {
		db *mongo.Client
	}
)

func NewPlayerTransaction(db *mongo.Client) IPlayerTransactionRepository {
	return &playerTransactionRepository{db}
}

func (r *playerTransactionRepository) dbConn() *mongo.Database {
	return r.db.Database("player_db")
}

func (r *playerTransactionRepository) FirstByField(pctx context.Context, field string, expected any) (*playerModule.PlayerTransaction, error) {
	db := r.dbConn()
	col := db.Collection("player_transactions")

	if field == "_id" {
		expected = utils.ConvertToObjectId(fmt.Sprintf("%v", expected))
	}

	var result playerModule.PlayerTransaction
	if err := col.FindOne(
		pctx,
		bson.M{field: expected},
		options.FindOne().SetProjection(bson.M{
			"_id":        1,
			"player_id":  1,
			"amount":     1,
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

func (r *playerTransactionRepository) ExistsByField(pctx context.Context, field string, expected any) (bool, error) {
	_, err := r.FirstByField(pctx, field, expected)
	if err != nil {
		if errors.Is(err, appError.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *playerTransactionRepository) Create(pctx context.Context, req playerModule.PlayerTransaction) (string, error) {
	col := r.dbConn().Collection("player_transactions")

	req.CreatedAt = utils.TimeNow()
	req.UpdatedAt = utils.TimeNow()

	playerTransactionId, err := col.InsertOne(pctx, req)
	if err != nil {
		return "", err
	}

	return playerTransactionId.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *playerTransactionRepository) GetSavingAccount(pctx context.Context, playerId string) (*playerModule.PlayerSavingAccount, error) {
	col := r.dbConn().Collection("player_transactions")
	fmt.Println(playerId)
	filter := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerId}}}},
		bson.D{
			{Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$player_id"},
					{Key: "balance", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
				},
			},
		},
		bson.D{{Key: "$project",
			Value: bson.D{
				{Key: "player_id", Value: "$_id"},
				{Key: "_id", Value: 0},
				{Key: "balance", Value: 1},
			},
		}},
	}

	cursors, err := col.Aggregate(pctx, filter)
	if err != nil {
		return nil, err
	}

	var result playerModule.PlayerSavingAccount
	if cursors.Next(pctx) {
		if err := cursors.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}
