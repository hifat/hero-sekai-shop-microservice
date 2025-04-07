package paymentHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
)

type (
	paymentHttp struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentHttp(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *paymentHttp {
	return &paymentHttp{cfg, paymentUsecase}
}
