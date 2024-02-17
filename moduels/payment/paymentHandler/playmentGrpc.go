package paymentHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment/paymentUsecase"
)

type (
	httpPaymentGrpc struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentGrpc(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *httpPaymentGrpc {
	return &httpPaymentGrpc{cfg, paymentUsecase}
}
