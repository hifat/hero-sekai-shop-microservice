package itemHandler

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemUsecase"
)

type (
	itemGrpc struct {
		itemProto.UnimplementedItemGrpcServiceServer
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemGrpc(itemUsecase itemUsecase.IItemUsecase) *itemGrpc {
	return &itemGrpc{
		itemUsecase: itemUsecase,
	}
}

func (g *itemGrpc) FindItemInIds(context.Context, *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error) {
	return nil, nil
}
