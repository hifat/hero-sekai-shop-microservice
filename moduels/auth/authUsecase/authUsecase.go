package authUsecase

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/auth/authRepository"

type (
	IAuthUsecase interface{}

	authUsecase struct {
		authRepo authRepository.IAuthRepository
	}
)

func NewAuthUsecase(authRepo authRepository.IAuthRepository) IAuthUsecase {
	return &authUsecase{authRepo}
}
