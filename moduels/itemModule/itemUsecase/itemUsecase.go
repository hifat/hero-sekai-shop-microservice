package itemUsecase

import (
	"context"
	"errors"
	"fmt"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IItemUsecase interface {
		Create(pctx context.Context, req *itemModule.CreateItemReq) (*itemModule.ItemShowCase, error)
		FindById(pctx context.Context, itemId string) (*itemModule.ItemShowCase, error)
		Find(pctx context.Context, basePaginateUrl string, req *itemModule.ItemSearchReq) (*model.PaginateRes, error)
	}

	itemUsecase struct {
		itemRepo itemRepository.IItemRepository
	}
)

func NewItem(itemRepo itemRepository.IItemRepository) IItemUsecase {
	return &itemUsecase{itemRepo}
}

func (u *itemUsecase) Create(pctx context.Context, req *itemModule.CreateItemReq) (*itemModule.ItemShowCase, error) {
	isUnique, err := u.itemRepo.IsUnique(pctx, req.Title)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !isUnique {
		return nil, errors.New("duplicated item")
	}

	itemId, err := u.itemRepo.Create(pctx, req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return u.FindById(pctx, itemId)
}

func (u *itemUsecase) FindById(pctx context.Context, itemId string) (*itemModule.ItemShowCase, error) {

	result, err := u.itemRepo.FindById(pctx, itemId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &itemModule.ItemShowCase{
		ItemId:   result.Id.Hex(),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}

func (u *itemUsecase) Find(pctx context.Context, basePaginateUrl string, req *itemModule.ItemSearchReq) (*model.PaginateRes, error) {
	findItemsFilter := bson.D{}
	findItemsOpts := make([]*options.FindOptions, 0)

	countItemsFilter := bson.D{}

	// Filter
	if req.Start != "" {
		findItemsFilter = append(findItemsFilter, bson.E{Key: "_id", Value: bson.D{
			{Key: "$gt", Value: utils.ConvertToObjectId(req.Start)},
		}})
	}

	if req.Title != "" {
		findItemsFilter = append(findItemsFilter, bson.E{Key: "title", Value: primitive.Regex{
			Pattern: req.Title,
			Options: "i",
		}})
		countItemsFilter = append(countItemsFilter, bson.E{Key: "title", Value: primitive.Regex{
			Pattern: req.Title,
			Options: "i",
		}})
	}

	findItemsFilter = append(findItemsFilter, bson.E{Key: "usage_status", Value: true})
	countItemsFilter = append(countItemsFilter, bson.E{Key: "usage_status", Value: true})

	// Options
	findItemsOpts = append(findItemsOpts, options.Find().SetSort(bson.D{
		{Key: "_id", Value: 1},
	}))
	findItemsOpts = append(findItemsOpts, options.Find().SetLimit(req.Limit))

	// Find
	fmt.Printf("%+v", findItemsFilter)
	results, err := u.itemRepo.Find(pctx, findItemsFilter, findItemsOpts)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// Count
	total, err := u.itemRepo.Count(pctx, countItemsFilter)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(results) == 0 {
		return &model.PaginateRes{
			Data:  results,
			Total: total,
			Limit: req.Limit,
			First: model.FirstPaginate{
				Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
			},
			Next: model.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	start := results[len(results)-1].ItemId

	return &model.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: model.FirstPaginate{
			Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
		},
		Next: model.NextPaginate{
			Start: start,
			Href:  fmt.Sprintf("%s?limit=%d&title=%s&start=%s", basePaginateUrl, req.Limit, req.Title, start),
		},
	}, nil
}
