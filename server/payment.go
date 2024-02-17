package server

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment"

func (s *server) paymentService() {
	paymentGroup := s.app.Group("payment_v1")

	paymentandler := payment.InitPayment(s.cfg, s.db)
	_ = paymentGroup
	_ = paymentandler
}
