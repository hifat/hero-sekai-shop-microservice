package paymentHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment/paymentUsecase"
)

type (
	httpPaymentQueue struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentQueue(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *httpPaymentQueue {
	return &httpPaymentQueue{cfg, paymentUsecase}
}
