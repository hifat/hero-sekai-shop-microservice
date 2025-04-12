package inventoryRepository

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
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
		FindItemInIds(pctx context.Context, grpcUrl string, req *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error)
		FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventoryModule.Inventory, error)
		CountByPlayerId(pctx context.Context, playerId string) (int64, error)
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

	utils.Debug(results)

	return results, nil
}

func (r *inventoryRepository) CountByPlayerId(pctx context.Context, playerId string) (int64, error) {
	db := r.dbConn()
	col := db.Collection("player_inventories")

	return col.CountDocuments(pctx, bson.M{
		"player_id": utils.ConvertToObjectId(playerId),
	})
}
