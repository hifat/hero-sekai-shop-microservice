package server

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment/paymentDI"
)

func (s *server) paymentService() {
	paymentGroup := s.app.Group("payment_v1")

	paymentandler := paymentDI.InitPayment(s.cfg, s.db)
	_ = paymentandler

	paymentGroup.GET("", s.healthCheckService)
}
