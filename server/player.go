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

	playerGroup := s.app.Group("player_v1")

	_ = playerHandler

	playerGroup.GET("", s.healthCheckService)
}
