package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	grpcHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewGrpc(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) grpcHandler {
	return grpcHandler{cfg, authUsecase}
}
