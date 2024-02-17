//go:build wireinject
// +build wireinject

package inventory

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	inventoryRepository.NewInventory,
)

var UsecaseSet = wire.NewSet(
	inventoryUsecase.NewInventory,
)

var HandlerSet = wire.NewSet(
	inventoryHandler.NewHandler,
	inventoryHandler.NewInventoryHttp,
	inventoryHandler.NewInventoryGrpc,
	inventoryHandler.NewInventoryQueue,
)

func InitInventory(cfg *config.Config, db *mongo.Client) inventoryHandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return inventoryHandler.Handler{}
}
