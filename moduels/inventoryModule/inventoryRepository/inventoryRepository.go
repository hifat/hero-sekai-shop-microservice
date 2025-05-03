package inventoryRepository

import (
	"context"
	"errors"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IInventoryRepository interface {
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		FindItemInIds(pctx context.Context, grpcUrl string, req *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error)
		FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventoryModule.Inventory, error)
		CountByPlayerId(pctx context.Context, playerId string) (int64, error)
		AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *paymentModule.PaymentTransferRes) error
		RemovePlayerItemRes(pctx context.Context, cfg *config.Config, req *paymentModule.PaymentTransferRes) error
		CreatePlayerItem(pctx context.Context, req *inventoryModule.Inventory) (string, error)
		DeleteById(pctx context.Context, id string) error
		IsExistedByPlayerItem(pctx context.Context, playerId string, itemId string) (bool, error)
		DeletePlayerItem(pctx context.Context, playerId string, itemId string) error
	}

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

func (r *inventoryRepository) GetOffset(pctx context.Context) (int64, error) {
	db := r.dbConn()
	col := db.Collection("player_inventories_queue")

	result := new(model.KafkaOffset)
	if err := col.FindOne(pctx, bson.M{}).
		Decode(result); err != nil {
		return -1, err
	}

	return result.Offset, nil
}

func (r *inventoryRepository) UpsertOffset(pctx context.Context, offset int64) error {
	db := r.dbConn()
	col := db.Collection("player_inventories_queue")

	_, err := col.UpdateOne(
		pctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"offset": offset,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *inventoryRepository) FindItemInIds(pctx context.Context, grpcUrl string, req *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error) {
	jwtauth.SetApiKeyInContext(&pctx)

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}

	result, err := conn.Item().FindItemInIds(pctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *inventoryRepository) FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventoryModule.Inventory, error) {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	cursors, err := col.Find(pctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	results := make([]*inventoryModule.Inventory, 0)
	for cursors.Next(pctx) {
		result := new(inventoryModule.Inventory)
		if err := cursors.Decode(result); err != nil {
			logger.Error(err)
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *inventoryRepository) CountByPlayerId(pctx context.Context, playerId string) (int64, error) {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	return col.CountDocuments(pctx, bson.M{
		"player_id": utils.ConvertToObjectId(playerId),
	})
}

func (r *inventoryRepository) CreatePlayerItem(pctx context.Context, req *inventoryModule.Inventory) (string, error) {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	result, err := col.InsertOne(pctx, bson.M{
		"player_id": utils.ConvertToObjectId(req.PlayerId),
		"item_id":   utils.ConvertToObjectId(req.ItemId),
	})
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *inventoryRepository) DeleteById(pctx context.Context, id string) error {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	_, err := col.DeleteOne(pctx, bson.M{
		"_id": utils.ConvertToObjectId(id),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *inventoryRepository) IsExistedByPlayerItem(pctx context.Context, playerId string, itemId string) (bool, error) {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	result := new(inventoryModule.Inventory)
	if err := col.FindOne(pctx, bson.M{
		"player_id": utils.ConvertToObjectId(playerId),
		"item_id":   utils.ConvertToObjectId(itemId),
	}).Decode(result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *inventoryRepository) DeletePlayerItem(pctx context.Context, playerId string, itemId string) error {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	_, err := col.DeleteOne(pctx, bson.M{
		"player_id": utils.ConvertToObjectId(playerId),
		"item_id":   utils.ConvertToObjectId(itemId),
	})
	if err != nil {
		return err
	}

	return nil
}
