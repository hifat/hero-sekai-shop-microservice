package server

import (
	"fmt"
	"log/slog"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
)

func (s *server) authService() {
	authHandler := auth.InitAuth(s.cfg, s.db)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authProto.RegisterAuthGrpcServiceServer(grpcServer, authHandler.AuthGrpc)

		slog.Info(fmt.Sprintf("Auth gRPC server listening on: %s", s.cfg.Grpc.AuthUrl))
		grpcServer.Serve(lis)
	}()

	authGroup := s.app.Group("auth_v1")

	_ = authHandler

	authGroup.GET("", s.healthCheckService)
}
