//go:build wireinject
// +build wireinject

package paymentDI

import (
	"github.com/google/wire"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	paymentandler "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentHandler"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

var RepoSet = wire.NewSet(
	paymentRepository.NewPayment,
)

var UsecaseSet = wire.NewSet(
	paymentUsecase.NewPayment,
)

var HandlerSet = wire.NewSet(
	paymentandler.NewHandler,
	paymentandler.Newpaymentttp,
	paymentandler.NewPaymentGrpc,
	paymentandler.NewPaymentQueue,
)

func InitPayment(cfg *config.Config, db *mongo.Client) paymentandler.Handler {
	wire.Build(
		RepoSet,
		UsecaseSet,
		HandlerSet,
	)

	return paymentandler.Handler{}
}
