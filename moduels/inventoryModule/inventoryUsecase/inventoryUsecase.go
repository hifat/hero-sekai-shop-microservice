package inventoryUsecase

import (
	"context"
	"fmt"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IInventoryUsecase interface {
		FindPlayerItems(pctx context.Context, basePaginationUrl, playerId string, req *inventoryModule.InventorySearchReq) (*model.PaginateRes, error)
	}

	inventoryUsecase struct {
		inventoryRepo inventoryRepository.IInventoryRepository
	}
)

func NewInventory(inventoryRepo inventoryRepository.IInventoryRepository) IInventoryUsecase {
	return &inventoryUsecase{inventoryRepo}
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, basePaginateUrl, playerId string, req *inventoryModule.InventorySearchReq) (*model.PaginateRes, error) {
	// Filter
	filter := bson.D{}

	if req.Start != "" {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{
			{Key: "$gt", Value: utils.ConvertToObjectId(req.Start)},
		}})
	}
	filter = append(filter, bson.E{Key: "player_id", Value: utils.ConvertToObjectId(playerId)})

	// Option
	opts := make([]*options.FindOptions, 0)

	opts = append(opts, options.Find().SetSort(bson.D{
		{Key: "_id", Value: 1},
	}))
	opts = append(opts, options.Find().SetLimit(req.Limit))

	// Find
	inventories, err := u.inventoryRepo.FindPlayerItems(pctx, filter, opts)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	results := make([]*inventoryModule.ItemInInventory, 0)
	for _, inventory := range inventories {
		results = append(results, &inventoryModule.ItemInInventory{
			InventoryId: inventory.Id,
			PlayerId:    inventory.PlayerId,
			ItemShowCase: &itemModule.ItemShowCase{
				ItemId: inventory.ItemId,
			},
		})
	}

	// Count
	total, err := u.inventoryRepo.CountByPlayerId(pctx, playerId)
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
				Href: fmt.Sprintf("%s?limit=%d", basePaginateUrl, req.Limit),
			},
			Next: model.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	start := results[len(results)-1].InventoryId

	return &model.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: model.FirstPaginate{
			Href: fmt.Sprintf("%s/my-item?limit=%d", basePaginateUrl, req.Limit),
		},
		Next: model.NextPaginate{
			Start: start,
			Href:  fmt.Sprintf("%s/my-item?limit=%d&start=%s", basePaginateUrl, req.Limit, start),
		},
	}, nil
}
