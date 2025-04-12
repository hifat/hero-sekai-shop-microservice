package itemHandler

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemUsecase"
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

func (g *itemGrpc) FindItemInIds(ctx context.Context, req *itemProto.FindItemsInIdsReq) (*itemProto.FindItemsInIdsRes, error) {
	return g.itemUsecase.FindItemInIds(ctx, req)
}
