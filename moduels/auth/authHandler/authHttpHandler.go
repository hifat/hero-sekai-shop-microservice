package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	httpHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewHttp(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) httpHandler {
	return httpHandler{cfg, authUsecase}
}
