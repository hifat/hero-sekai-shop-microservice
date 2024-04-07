package middlewareHandler

import (
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
