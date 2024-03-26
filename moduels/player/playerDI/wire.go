//go:build wireinject
// +build wireinject

package playerDI

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	playerRepository.NewPlayer,
	playerRepository.NewPlayerTransaction,
)

var UsecaseSet = wire.NewSet(
	playerUsecase.NewPlayer,
	playerUsecase.NewPlayerTransaction,
)

var HandlerSet = wire.NewSet(
	playerHandler.NewHandler,
	playerHandler.NewPlayerHttp,
	playerHandler.NewPlayerTransactionHttp,
	playerHandler.NewPlayerGrpc,
	playerHandler.NewPlayerQueue,
)

func InitPlayer(cfg *config.Config, db *mongo.Client) playerHandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return playerHandler.Handler{}
}
