package server

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item"

func (s *server) itemService() {
	itemGroup := s.app.Group("item_v1")

	itemHandler := item.InitItem(s.cfg, s.db)
	_ = itemGroup
	_ = itemHandler
}
