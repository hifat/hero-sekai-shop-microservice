package server

import (
	"fmt"
	"log/slog"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemDI"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
)

func (s *server) itemService() {
	itemHandler := itemDI.InitItem(s.cfg, s.db)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemProto.RegisterItemGrpcServiceServer(grpcServer, itemHandler.ItemGrpc)

		slog.Info(fmt.Sprintf("Item gRPC server listening on: %s", s.cfg.Grpc.ItemUrl))
		grpcServer.Serve(lis)
	}()

	itemGroup := s.app.Group("item_v1")

	_ = itemHandler

	itemGroup.GET("", s.healthCheckService)
	itemGroup.POST("/item", s.middleware.MiddlewareHttp.JwtAuth(
		s.middleware.MiddlewareHttp.RbcaAuth(
			itemHandler.ItemHttp.Create,
			[]int{1, 0},
		),
	))
	itemGroup.GET("/item", itemHandler.ItemHttp.Find)
	itemGroup.GET("/item/:item_id", itemHandler.ItemHttp.FindById)
}
