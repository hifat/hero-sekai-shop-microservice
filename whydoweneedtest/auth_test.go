package whydoweneedtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
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
	zeroTime := time.Time{}.String()

	_ = usecase

	repoMock := new(authRepository.AuthRepositoryMock)

	ctx := context.Background()
	grpcUrl := "http://localhost:3000"

	repoMock.On("CredentialSearch", ctx, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    "credential_pass@sekai.com",
		Password: "123456",
	}).Return(&playerProto.PlayerProfile{
		Id:        "001",
		Email:     "success@sekai.com",
		Username:  "player001",
		RoleCode:  0,
		CreatedAt: zeroTime,
		UpdatedAt: zeroTime,
	}, errors.New("mock_error"))

	repoMock.On("CredentialSearch", ctx, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    "err@sekai.com",
		Password: "123456",
	}).Return(nil, errors.New("mock_error"))

	repoMock.On("NewAccessToken", cfg, mock.AnythingOfType("*playerProto.PlayerProfile")).Return("xxx")

	repoMock.On("NewRefreshToken", ctx, mock.AnythingOfType("*playerProto.PlayerProfile")).Return("xxx")

	repoMock.On("InsertOnePlayerCredential", ctx, &authModule.Credential{
		PlayerId:     "player:001",
		RoleCode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    &time.Time{},
		UpdatedAt:    &time.Time{},
	}).Return("player_id", nil)

	repoMock.On("InsertOnePlayerCredential", ctx, &authModule.Credential{
		PlayerId:     "player:001",
		RoleCode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    &time.Time{},
		UpdatedAt:    &time.Time{},
	}).Return("player_id", nil)
}
