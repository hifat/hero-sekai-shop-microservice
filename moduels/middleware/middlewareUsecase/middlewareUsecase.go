package middlewareUsecase

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareRepository"
)

type (
	IUsecase interface{}

	usecase struct {
		middlewareRepo middlewareRepository.IRepository
	}
)

func New(middlewareRepo middlewareRepository.IRepository) IUsecase {
	return &usecase{middlewareRepo}
}
