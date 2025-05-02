package paymentRepository

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IPaymentRepository interface {
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		FindItemInIds(pctx context.Context, grpcUrl string, req *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error)
		DockedPlayerMoney(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq) error
		RollbackDockedPlayerMoney(pctx context.Context, cfg *config.Config, req *playerModule.RollbackPlayerTransactionReq) error
		AddPlayerItem(pctx context.Context, cfg *config.Config, req *inventoryModule.UpdateInventoryReq) error
		RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventoryModule.RollbackPlayerInventoryReq) error
	}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPayment(db *mongo.Client) IPaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) dbConn() *mongo.Database {
	return r.db.Database("payment_db")
}

func (r *paymentRepository) GetOffset(pctx context.Context) (int64, error) {
	db := r.dbConn()
	col := db.Collection("payment_queue")

	result := new(model.KafkaOffset)
	if err := col.FindOne(pctx, bson.M{}).
		Decode(result); err != nil {
		return -1, err
	}

	return result.Offset, nil
}

func (r *paymentRepository) UpsertOffset(pctx context.Context, offset int64) error {
	db := r.dbConn()
	col := db.Collection("payment_queue")

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

func (r *paymentRepository) FindItemInIds(pctx context.Context, grpcUrl string, req *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error) {
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
