package inventoryUsecase

import (
	"context"
	"fmt"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/model"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IInventoryUsecase interface {
		FindPlayerItems(pctx context.Context, playerId string, req *inventoryModule.InventorySearchReq) (*model.PaginateRes, error)
	}

	inventoryUsecase struct {
		cfg           *config.Config
		inventoryRepo inventoryRepository.IInventoryRepository
	}
)

func NewInventory(cfg *config.Config, inventoryRepo inventoryRepository.IInventoryRepository) IInventoryUsecase {
	return &inventoryUsecase{cfg, inventoryRepo}
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, playerId string, req *inventoryModule.InventorySearchReq) (*model.PaginateRes, error) {
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

	itemData, err := u.inventoryRepo.FindItemInIds(pctx, u.cfg.Grpc.ItemUrl, &itemProto.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0, len(inventories))
			for _, v := range inventories {
				itemIds = append(itemIds, v.ItemId)
			}

			return itemIds
		}(),
	})
	if err != nil {
		return nil, err
	}

	itemMap := make(map[string]*itemModule.ItemShowCase, len(itemData.Items))
	for _, v := range itemData.Items {
		itemMap[v.Id] = &itemModule.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			Damage:   v.Damage,
			ImageUrl: v.ImageUrl,
		}
	}

	results := make([]*inventoryModule.ItemInInventory, 0)
	for _, inventory := range inventories {
		inventoryRes := &inventoryModule.ItemInInventory{
			InventoryId: inventory.Id,
			PlayerId:    inventory.PlayerId,
		}

		item, ok := itemMap[inventory.ItemId]
		if ok {
			inventoryRes.ItemShowCase = &itemModule.ItemShowCase{
				ItemId:   inventory.ItemId,
				Title:    item.Title,
				Price:    item.Price,
				Damage:   item.Damage,
				ImageUrl: item.ImageUrl,
			}
		}

		results = append(results, inventoryRes)
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
				Href: fmt.Sprintf("%s?limit=%d", u.cfg.Paginate.InventoryNextPageBasedUrl, req.Limit),
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
			Href: fmt.Sprintf("%s/my-item?limit=%d", u.cfg.Paginate.InventoryNextPageBasedUrl, req.Limit),
		},
		Next: model.NextPaginate{
			Start: start,
			Href:  fmt.Sprintf("%s/my-item?limit=%d&start=%s", u.cfg.Paginate.InventoryNextPageBasedUrl, req.Limit, start),
		},
	}, nil
}
