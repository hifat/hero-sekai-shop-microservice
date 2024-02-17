package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) authHttpHandler {
	return authHttpHandler{cfg, authUsecase}
}
