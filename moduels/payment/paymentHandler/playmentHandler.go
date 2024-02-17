package paymentHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment/paymentUsecase"
)

type (
	httpPaymentHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentHttp(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) httpPaymentHandler {
	return httpPaymentHandler{cfg, paymentUsecase}
}
