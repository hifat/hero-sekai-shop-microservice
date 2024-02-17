package server

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory"

func (s *server) inventoryService() {
	inventoryGroup := s.app.Group("inventory_v1")

	inventoryHandler := inventory.InitInventory(s.cfg, s.db)
	_ = inventoryGroup
	_ = inventoryHandler
}
