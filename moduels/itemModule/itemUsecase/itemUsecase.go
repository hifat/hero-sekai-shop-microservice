package itemUsecase

import (
	"context"
	"errors"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
)

type (
	IItemUsecase interface {
		Create(pctx context.Context, req *itemModule.CreateItemReq) (*itemModule.ItemShowCase, error)
		FindById(pctx context.Context, itemId string) (*itemModule.ItemShowCase, error)
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
