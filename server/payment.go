package server

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentDI"
)

func (s *server) paymentService() {
	paymentGroup := s.app.Group("payment_v1")

	paymentHandler := paymentDI.InitPayment(s.cfg, s.db)

	paymentGroup.GET("", s.healthCheckService)

	paymentGroup.POST("/buy", paymentHandler.PaymentHttp.BuyItem, s.middleware.MiddlewareHttp.JwtAuth)
	paymentGroup.POST("/sell", paymentHandler.PaymentHttp.SellItem, s.middleware.MiddlewareHttp.JwtAuth)
}
