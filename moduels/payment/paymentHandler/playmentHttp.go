package paymentHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment/paymentUsecase"
)

type (
	httpPaymentHttp struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentHttp(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) httpPaymentHttp {
	return httpPaymentHttp{cfg, paymentUsecase}
}
