//go:build wireinject
// +build wireinject

package middleware

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	middlewareRepository.NewMiddleware,
)

var UsecaseSet = wire.NewSet(
	middlewareUsecase.NewMiddleware,
)

var HandlerSet = wire.NewSet(
	middlewareHandler.NewHandler,
	middlewareHandler.NewHttp,
)

func InitMiddleware(cfg *config.Config, db *mongo.Client) middlewareHandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return middlewareHandler.Handler{}
}
