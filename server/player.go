package server

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player"

func (s *server) playerService() {
	playerGroup := s.app.Group("player_v1")

	playerHandler := player.InitPlayer(s.cfg, s.db)
	_ = playerHandler

	playerGroup.GET("", s.healthCheckService)
}
