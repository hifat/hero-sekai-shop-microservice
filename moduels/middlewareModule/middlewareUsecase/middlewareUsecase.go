package middlewareUsecase

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middlewareModule/middlewareRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
)

type (
	IMiddlewareUsecase interface {
		JwtAuth(pctx context.Context, accessToken string) (*jwtauth.AuthMapClaims, error)
	}

	middlewareUsecase struct {
		cfg            *config.Config
		middlewareRepo middlewareRepository.IMiddlewareRepository
	}
)

func NewMiddleware(cfg *config.Config, middlewareRepo middlewareRepository.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
		cfg,
		middlewareRepo,
	}
}

func (u *middlewareUsecase) JwtAuth(pctx context.Context, accessToken string) (*jwtauth.AuthMapClaims, error) {
	claims, err := jwtauth.ParseToken(u.cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := u.middlewareRepo.AccessTokenSearch(pctx, u.cfg.Grpc.AuthUrl, accessToken); err != nil {
		logger.Error(err)
		return nil, err
	}

	return claims, nil
}
