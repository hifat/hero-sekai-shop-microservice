package whydoweneedtest

import (
	"context"
	"errors"
	"testing"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
)

type testLogin struct {
	ctx      context.Context
	cfg      *config.Config
	req      *authModule.PlayerLoginReq
	expected *authModule.ProfileIntercepter
	err      error
	isErr    bool
}

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
	}).Return(nil, errors.New("mock_error"))

}
