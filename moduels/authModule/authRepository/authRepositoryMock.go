package authRepository

import (
	"context"

	"github.com/stretchr/testify/mock"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func NewAuthRepositoryMock() IAuthRepository {
	return &AuthRepositoryMock{}
}

func (m *AuthRepositoryMock) CredentialSearch(pctx context.Context, grpcUrl string, req *playerProto.CredentialSearchReq) (*playerProto.PlayerProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	return args.Get(0).(*playerProto.PlayerProfile), args.Error(1)
}

func (m *AuthRepositoryMock) InsertOne(pctx context.Context, req *authModule.Credential) (string, error) {
	args := m.Called(pctx, req)
	return args.String(0), args.Error(1)
}

func (m *AuthRepositoryMock) FindByCredentialId(pctx context.Context, credentialId string) (*authModule.Credential, error) {
	args := m.Called(pctx, credentialId)
	return args.Get(0).(*authModule.Credential), args.Error(1)
}

func (m *AuthRepositoryMock) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerProto.FindOnePlayerProfileToRefreshReq) (*playerProto.PlayerProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	return args.Get(0).(*playerProto.PlayerProfile), args.Error(1)
}

func (m *AuthRepositoryMock) UpdateRefreshToken(pctx context.Context, credentialId string, req *authModule.UpdateRefreshTokenReq) error {
	args := m.Called(pctx, credentialId, req)
	return args.Error(1)
}

func (m *AuthRepositoryMock) DeleteByCredentialId(pctx context.Context, credentialId string) (int64, error) {
	args := m.Called(pctx, credentialId)
	return int64(args.Int(0)), args.Error(1)
}

func (m *AuthRepositoryMock) FindByAccessToken(pctx context.Context, accessToken string) (*authModule.Credential, error) {
	args := m.Called(pctx, accessToken)
	return args.Get(0).(*authModule.Credential), args.Error(1)
}

func (m *AuthRepositoryMock) RoleCount(pctx context.Context) (int64, error) {
	args := m.Called(pctx)
	return int64(args.Int(0)), args.Error(1)
}

func (m *AuthRepositoryMock) NewAccessToken(cfg config.Jwt, profile *playerProto.PlayerProfile) string {
	args := m.Called(cfg, profile)
	return args.String(0)
}

func (m *AuthRepositoryMock) NewRefreshToken(cfg config.Jwt, profile *playerProto.PlayerProfile) string {
	args := m.Called(cfg, profile)
	return args.String(0)
}

func (m *AuthRepositoryMock) ParseToken(secret string, tokenString string) (*jwtauth.AuthMapClaims, error) {
	args := m.Called(secret, tokenString)
	return args.Get(0).(*jwtauth.AuthMapClaims), args.Error(1)
}
