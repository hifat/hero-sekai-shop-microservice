package server

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryDI"
)

func (s *server) inventoryService() {
	inventoryHandler := inventoryDI.InitInventory(s.cfg, s.db)

	go inventoryHandler.InventoryQueue.AddPlayerItem()
	go inventoryHandler.InventoryQueue.RollbackAddPlayerItem()
	go inventoryHandler.InventoryQueue.RemovePlayerItem()
	go inventoryHandler.InventoryQueue.RollbackRemovePlayerItem()

	inventoryGroup := s.app.Group("inventory_v1")

	inventoryGroup.GET("", s.healthCheckService)
	inventoryGroup.GET("/inventory/my-item", inventoryHandler.InventoryHttp.FindPlayerItems, s.middleware.MiddlewareHttp.JwtAuth)
}
