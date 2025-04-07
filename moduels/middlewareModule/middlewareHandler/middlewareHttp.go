package middlewareHandler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middlewareModule/middlewareUsecase"
)

type (
	middlewareHttp struct {
		middlewareUsecase middlewareUsecase.IMiddlewareUsecase
	}
)

func NewHttp(middlewareUsecase middlewareUsecase.IMiddlewareUsecase) *middlewareHttp {
	return &middlewareHttp{middlewareUsecase}
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
