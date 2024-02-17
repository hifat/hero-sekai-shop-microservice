package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	authGrpcHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewAuthGrpcHandler(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) authGrpcHandler {
	return authGrpcHandler{cfg, authUsecase}
}
