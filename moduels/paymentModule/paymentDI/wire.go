//go:build wireinject
// +build wireinject

package paymentDI

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	paymentRepository.NewPayment,
	playerRepository.NewPlayer,
	playerRepository.NewPlayerTransaction,
)

var UsecaseSet = wire.NewSet(
	paymentUsecase.NewPayment,
	playerUsecase.NewPlayerTransaction,
)

var HandlerSet = wire.NewSet(
	paymentHandler.NewHandler,
	paymentHandler.NewPaymentHttp,
	paymentHandler.NewPaymentGrpc,
)

func InitPayment(cfg *config.Config, db *mongo.Client) paymentHandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return paymentHandler.Handler{}
}
