package middlewareUsecase

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middlewareModule/middlewareRepository"
)

type (
	IMiddlewareUsecase interface{}

	middlewareUsecase struct {
		middlewareRepo middlewareRepository.IMiddlewareRepository
	}
)

func NewMiddleware(middlewareRepo middlewareRepository.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{middlewareRepo}
}
