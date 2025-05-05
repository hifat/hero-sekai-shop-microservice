package whydoweneedtest

import (
	"context"
	"testing"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
)

func TestLogin(t *testing.T) {
	cfg := NewTestConfig()
	repo := authRepository.NewAuthRepositoryMock()
	usecase := authUsecase.NewAuth(cfg, repo)
	_ = usecase

	mock := new(authRepository.AuthRepositoryMock)

	ctx := context.Background()
	grpcUrl := "http://localhost:3000"

	mock.On("CredentialSearch", ctx, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    "test@sekai.com",
		Password: "123456",
	}).Return(nil, nil)
}
