package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	authGrpc struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewAuthGrpc(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) *authGrpc {
	return &authGrpc{cfg, authUsecase}
}
