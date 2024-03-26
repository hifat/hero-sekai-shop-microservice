package server

import (
	"fmt"
	"log/slog"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerDI"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
)

func (s *server) playerService() {
	playerHandler := playerDI.InitPlayer(s.cfg, s.db)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PlayerUrl)

		playerProto.RegisterPlayerGrpcServiceServer(grpcServer, playerHandler.PlayerGrpc)

		slog.Info(fmt.Sprintf("Player gRPC server listening on: %s", s.cfg.Grpc.PlayerUrl))
		grpcServer.Serve(lis)
	}()

	playerV1Group := s.app.Group("player_v1")

	playerV1Group.GET("/healtz", s.healthCheckService)

	playerGroup := playerV1Group.Group("/players")

	playerGroup.POST("/register", playerHandler.PlayerHttp.Create)
	playerGroup.GET("/:player_id", playerHandler.PlayerHttp.GetProfile)

	playerTransaction := playerV1Group.Group("/transactions")
	playerTransaction.POST("", playerHandler.PlayerTransactionHttp.AddMoney)
}
