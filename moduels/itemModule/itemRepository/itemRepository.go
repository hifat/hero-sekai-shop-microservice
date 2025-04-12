package itemRepository

import (
	"context"
	"errors"
	"time"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IItemRepository interface {
		IsUnique(pctx context.Context, title string, exceptId string) (bool, error)
		Create(pctx context.Context, req *itemModule.CreateItemReq) (string, error)
		FindById(pctx context.Context, itemId string) (*itemModule.Item, error)
		Find(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*itemModule.ItemShowCase, error)
		Count(pctx context.Context, filter primitive.D) (int64, error)
		Update(pctx context.Context, itemId string, req *itemModule.ItemUpdateReq) error
		UpdateUsageStatus(pctx context.Context, itemId string, isActive bool) error
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

func (r *itemRepository) IsUnique(pctx context.Context, title string, exceptId string) (bool, error) {
	db := r.dbConn()
	col := db.Collection("items")

	cond := bson.M{
		"title": title,
	}

	if exceptId != "" {
		cond["_id"] = bson.M{
			"$ne": utils.ConvertToObjectId(exceptId),
		}
	}

	result := new(itemModule.Item)
	if err := col.FindOne(pctx, cond).Decode(result); err != nil {
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

func (r *itemRepository) Find(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*itemModule.ItemShowCase, error) {
	db := r.dbConn()
	col := db.Collection("items")

	cursors, err := col.Find(pctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	results := make([]*itemModule.ItemShowCase, 0)
	for cursors.Next(pctx) {
		result := new(itemModule.Item)
		if err := cursors.Decode(result); err != nil {
			return nil, err
		}

		results = append(results, &itemModule.ItemShowCase{
			ItemId:   result.Id.Hex(),
			Title:    result.Title,
			Price:    result.Price,
			Damage:   result.Damage,
			ImageUrl: result.ImageUrl,
		})
	}

	return results, nil
}

func (r *itemRepository) Count(pctx context.Context, filter primitive.D) (int64, error) {
	db := r.dbConn()
	col := db.Collection("items")

	return col.CountDocuments(pctx, filter)
}

func (r *itemRepository) Update(pctx context.Context, itemId string, req *itemModule.ItemUpdateReq) error {
	db := r.dbConn()
	col := db.Collection("items")

	_, err := col.UpdateOne(pctx, bson.M{
		"_id": utils.ConvertToObjectId(itemId),
	}, bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"price":      req.Price,
			"damage":     req.Damage,
			"image_url":  req.ImageUrl,
			"updated_at": time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *itemRepository) UpdateUsageStatus(pctx context.Context, itemId string, isActive bool) error {
	db := r.dbConn()
	col := db.Collection("items")

	_, err := col.UpdateOne(pctx, bson.M{
		"_id": utils.ConvertToObjectId(itemId),
	}, bson.M{
		"$set": bson.M{
			"usage_status": isActive,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
