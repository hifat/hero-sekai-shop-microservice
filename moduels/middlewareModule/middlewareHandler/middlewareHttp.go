package middlewareHandler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middlewareModule/middlewareUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
)

type (
	middlewareHttp struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecase.IMiddlewareUsecase
	}
)

func NewHttp(cfg *config.Config, middlewareUsecase middlewareUsecase.IMiddlewareUsecase) *middlewareHttp {
	return &middlewareHttp{
		cfg,
		middlewareUsecase,
	}
}

func (h *middlewareHttp) JwtAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		claims, err := h.middlewareUsecase.JwtAuth(c.Request().Context(), accessToken)
		if err != nil {
			return err
		}

		c.Set("credential", claims)

		return next(c)
	}
}

func (h *middlewareHttp) RbcaAuth(next echo.HandlerFunc, expected []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.RbacAuth(c, h.cfg, expected)
		if err != nil {
			return err
		}

		return next(newCtx)
	}
}

func (h *middlewareHttp) PlayerIdParamValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		playerIdReq := c.Param(("player_id"))
		claim := c.Get("credential").(*jwtauth.AuthMapClaims)

		newCtx, err := h.middlewareUsecase.PlayerIdParamValidation(c, playerIdReq, claim.PlayerId)
		if err != nil {
			return err
		}

		return next(newCtx)
	}
}
