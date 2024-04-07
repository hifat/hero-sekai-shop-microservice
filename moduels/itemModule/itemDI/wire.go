//go:build wireinject
// +build wireinject

package itemDI

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	itemRepository.NewItem,
)

var UsecaseSet = wire.NewSet(
	itemUsecase.NewItem,
)

var HandlerSet = wire.NewSet(
	itemHandler.NewHandler,
	itemHandler.NewItemHttp,
	itemHandler.NewItemGrpc,
)

func InitItem(cfg *config.Config, db *mongo.Client) itemHandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return itemHandler.Handler{}
}
