package inventoryHandler

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryUsecase"
)

type (
	inventoryGrpc struct {
		inventoryProto.UnimplementedInventoryGrpcServiceServer
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryGrpc(inventoryUsecase inventoryUsecase.IInventoryUsecase) *inventoryGrpc {
	return &inventoryGrpc{
		inventoryUsecase: inventoryUsecase,
	}
}

func (g *inventoryGrpc) IsAvailableToSell(context.Context, *inventoryProto.IsAvailableToSellReq) (*inventoryProto.IsAvailableToSellRes, error) {
	return nil, nil
}
