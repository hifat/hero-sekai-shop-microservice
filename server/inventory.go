package server

import (
	"fmt"
	"log/slog"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventory/inventoryProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
)

func (s *server) inventoryService() {
	inventoryHandler := inventory.InitInventory(s.cfg, s.db)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)

		inventoryProto.RegisterInventoryGrpcServiceServer(grpcServer, inventoryHandler.InventoryGrpc)

		slog.Info(fmt.Sprintf("Inventory gRPC server listening on: %s", s.cfg.Grpc.InventoryUrl))
		grpcServer.Serve(lis)
	}()

	inventoryGroup := s.app.Group("inventory_v1")

	_ = inventoryHandler

	inventoryGroup.GET("", s.healthCheckService)
}
