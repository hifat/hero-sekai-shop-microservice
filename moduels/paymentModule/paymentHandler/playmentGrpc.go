package paymentandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
)

type (
	paymentGrpc struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentGrpc(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *paymentGrpc {
	return &paymentGrpc{cfg, paymentUsecase}
}
