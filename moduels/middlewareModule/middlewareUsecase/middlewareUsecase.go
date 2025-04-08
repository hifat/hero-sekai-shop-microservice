package middlewareUsecase

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middlewareModule/middlewareRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/rbac"
)

type (
	IMiddlewareUsecase interface {
		JwtAuth(pctx context.Context, accessToken string) (*jwtauth.AuthMapClaims, error)
		RbacAuth(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error)
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

func (u *middlewareUsecase) RbacAuth(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error) {
	ctx := c.Request().Context()
	claim := c.Get("credential").(*jwtauth.AuthMapClaims)

	roleCount, err := u.middlewareRepo.RoleCount(ctx, cfg.Grpc.AuthUrl)
	if err != nil {
		return nil, err
	}

	playerRoleBinary := rbac.IntToBinary(int(claim.RoleCode), int(roleCount))

	for i := 0; i < int(roleCount); i++ {
		if playerRoleBinary[i]&expected[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("permission denied")
}
