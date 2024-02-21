package authHandler

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	authGrpc struct {
		authProto.UnimplementedAuthGrpcServiceServer
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewAuthGrpc(authUsecase authUsecase.IAuthUsecase) *authGrpc {
	return &authGrpc{
		authUsecase: authUsecase,
	}
}

func (g *authGrpc) AccessTokenSearch(context.Context, *authProto.AccessTokenSearchReq) (*authProto.AccessTokenSearchRes, error) {
	return nil, nil
}

func (g *authGrpc) RolesCount(context.Context, *authProto.RolesCountReq) (*authProto.RolesCountRes, error) {
	return nil, nil
}
