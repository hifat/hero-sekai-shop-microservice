package grpccon

import (
	"context"
	"errors"
	"fmt"
	"net"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type (
	GrpcClientFactoryHandler interface {
		Auth() authProto.AuthGrpcServiceClient
		Player() playerProto.PlayerGrpcServiceClient
		Item() itemProto.ItemGrpcServiceClient
		Inventory() inventoryProto.InventoryGrpcServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
	}
)

func (g *grpcAuth) unaryAuth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata not found")
	}

	authHandler, ok := md["x-api-key"]
	if !ok {
		return nil, errors.New("x-api-key metadata not found")
	}

	if len(authHandler) == 0 {
		return nil, errors.New("x-api-key is empty value")
	}

	claim, err := jwtauth.ParseToken(g.secretKey, string(authHandler[0]))
	if err != nil {
		return nil, errors.New("token is invalid")
	}
	logger.Info(fmt.Sprintf("claims: %+v", claim))

	return handler(ctx, req)
}

func (g *grpcClientFactory) Auth() authProto.AuthGrpcServiceClient {
	return authProto.NewAuthGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Player() playerProto.PlayerGrpcServiceClient {
	return playerProto.NewPlayerGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Item() itemProto.ItemGrpcServiceClient {
	return itemProto.NewItemGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Inventory() inventoryProto.InventoryGrpcServiceClient {
	return inventoryProto.NewInventoryGrpcServiceClient(g.client)
}

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
	opts := make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.NewClient(host, opts...)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &grpcClientFactory{
		client: clientConn,
	}, nil
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{
		secretKey: cfg.ApiSecretKey,
	}

	opts = append(opts, grpc.UnaryInterceptor(grpcAuth.unaryAuth))

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return grpcServer, lis
}
