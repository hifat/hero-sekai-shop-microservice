package authHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authUsecase"
)

type (
	authHttp struct {
		cfg         *config.Config
		authUsecase authUsecase.IAuthUsecase
	}
)

func NewAuthHttp(cfg *config.Config, authUsecase authUsecase.IAuthUsecase) *authHttp {
	return &authHttp{cfg, authUsecase}
}
