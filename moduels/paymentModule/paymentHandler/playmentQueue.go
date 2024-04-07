package paymentandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
)

type (
	paymentQueue struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentQueue(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *paymentQueue {
	return &paymentQueue{cfg, paymentUsecase}
}
