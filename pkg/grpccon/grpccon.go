package grpccon

import (
	"net"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
)

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

	clientConn, err := grpc.Dial(host, opts...)
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

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return grpcServer, lis
}
