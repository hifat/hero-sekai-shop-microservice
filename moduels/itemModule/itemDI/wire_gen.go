// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package itemDI

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitItem(cfg *config.Config, db *mongo.Client) itemHandler.Handler {
	iItemRepository := itemRepository.NewItem(db)
	iItemUsecase := itemUsecase.NewItem(iItemRepository)
	itemGrpc := itemHandler.NewItemGrpc(iItemUsecase)
	itemHttp := itemHandler.NewItemHttp(cfg, iItemUsecase)
	handler := itemHandler.NewHandler(itemGrpc, itemHttp)
	return handler
}

// wire.go:

var RepoSet = wire.NewSet(itemRepository.NewItem)

var UsecaseSet = wire.NewSet(itemUsecase.NewItem)

var HandlerSet = wire.NewSet(itemHandler.NewHandler, itemHandler.NewItemHttp, itemHandler.NewItemGrpc)
