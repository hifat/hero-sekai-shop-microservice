//go:build wireinject
// +build wireinject

package auth

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	authRepository.NewAuth,
)

var UsecaseSet = wire.NewSet(
	authUsecase.NewAuth,
)

var HandlerSet = wire.NewSet(
	authHandler.NewHandler,
	authHandler.NewAuthHttp,
	authHandler.NewAuthGrpc,
)

func InitAuth(cfg *config.Config, db *mongo.Client) authHandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return authHandler.Handler{}
}
