package server

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth"

func (s *server) authService() {
	authGroup := s.app.Group("auth_v1")

	authHandler := auth.InitAuth(s.cfg, s.db)
	_ = authHandler

	authGroup.GET("", s.healthCheckService)
}
