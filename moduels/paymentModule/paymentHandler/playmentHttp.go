package paymentandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
)

type (
	paymentttp struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func Newpaymentttp(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *paymentttp {
	return &paymentttp{cfg, paymentUsecase}
}
