package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	httpAuthHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewAuthHttp(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) *httpAuthHandler {
	return &httpAuthHandler{cfg, authUsecase}
}
