package itemRepository

import (
	"context"
	"errors"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IItemRepository interface {
		IsUnique(pctx context.Context, title string) (bool, error)
		Create(pctx context.Context, req *itemModule.CreateItemReq) (string, error)
		FindById(pctx context.Context, itemId string) (*itemModule.Item, error)
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItem(db *mongo.Client) IItemRepository {
	return &itemRepository{db}
}

func (r *itemRepository) dbConn() *mongo.Database {
	return r.db.Database("item_db")
}

func (r *itemRepository) IsUnique(pctx context.Context, title string) (bool, error) {
	db := r.dbConn()
	col := db.Collection("items")

	result := new(itemModule.Item)
	if err := col.FindOne(pctx, bson.M{
		"title": title,
	}).Decode(result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}

		return true, err
	}

	return false, nil
}

func (r *itemRepository) Create(pctx context.Context, req *itemModule.CreateItemReq) (string, error) {
	db := r.dbConn()
	col := db.Collection("items")

	result, err := col.InsertOne(pctx, itemModule.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		ImageUrl:    req.ImageUrl,
		UsageStatus: true,
		CreatedAt:   utils.TimeNow(),
		UpdatedAt:   utils.TimeNow(),
	})
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *itemRepository) FindById(pctx context.Context, itemId string) (*itemModule.Item, error) {
	db := r.dbConn()
	col := db.Collection("items")

	result := new(itemModule.Item)
	if err := col.FindOne(pctx, bson.M{
		"_id": utils.ConvertToObjectId(itemId),
	}).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
