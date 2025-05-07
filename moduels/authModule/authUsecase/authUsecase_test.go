package authUsecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/whydoweneedtest"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testLogin struct {
	ctx      context.Context
	req      *authModule.PlayerLoginReq
	expected *authModule.ProfileIntercepter
	err      error
	isErr    bool
}

func TestLogin(t *testing.T) {
	cfg := whydoweneedtest.NewTestConfig()
	repoMock := new(authRepository.AuthRepositoryMock)
	usecase := NewAuth(cfg, repoMock)
	zeroTime := time.Time{}.String()

	ctx := context.Background()
	grpcUrl := cfg.Grpc.PlayerUrl

	successEmail := "success@sekai.com"
	errEmail := "err@sekai.com"
	errInsertEmail := "err_insert@sekai.com"
	errFindByCredEmail := "err_find_by_cred_search@sekai.com"
	username := "player001"
	password := "123456"
	mockToken := "xxx"
	mockCredentialId := primitive.NewObjectID()
	mockCredentialIdErr := primitive.NewObjectID()
	mockPlayerId := primitive.NewObjectID().Hex()
	mockPlayerIdErr := "player_id_err"
	mockPlayerIdFindByCredErr := "player_id_find_by_cred_err"

	repoMock.On("CredentialSearch", mock.Anything, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    successEmail,
		Password: password,
	}).Return(&playerProto.PlayerProfile{
		Id:        mockPlayerId,
		Email:     successEmail,
		Username:  username,
		RoleCode:  0,
		CreatedAt: zeroTime,
		UpdatedAt: zeroTime,
	}, nil)

	repoMock.On("CredentialSearch", mock.Anything, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    errEmail,
		Password: password,
	}).Return(&playerProto.PlayerProfile{}, errors.New("mock_error"))

	repoMock.On("NewAccessToken", cfg.Jwt, mock.AnythingOfType("*playerProto.PlayerProfile")).
		Return(mockToken)

	repoMock.On("NewRefreshToken", cfg.Jwt, mock.AnythingOfType("*playerProto.PlayerProfile")).
		Return(mockToken)

	repoMock.On("InsertOne", ctx, &authModule.Credential{
		PlayerId:     mockPlayerId,
		RoleCode:     0,
		AccessToken:  mockToken,
		RefreshToken: mockToken,
	}).Return(mockCredentialId.Hex(), nil)

	repoMock.On("CredentialSearch", mock.Anything, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    errInsertEmail,
		Password: password,
	}).Return(&playerProto.PlayerProfile{
		Id:        mockPlayerIdErr,
		Email:     errInsertEmail,
		Username:  username,
		RoleCode:  0,
		CreatedAt: zeroTime,
		UpdatedAt: zeroTime,
	}, nil)

	repoMock.On("InsertOne", ctx, &authModule.Credential{
		PlayerId:     mockPlayerIdErr,
		RoleCode:     0,
		AccessToken:  mockToken,
		RefreshToken: mockToken,
	}).Return("", errors.New("mock_err"))

	repoMock.On("FindByCredentialId", ctx, mockCredentialId.Hex()).
		Return(&authModule.Credential{
			Id:           mockCredentialId,
			PlayerId:     mockPlayerId,
			RoleCode:     0,
			RefreshToken: mockToken,
			AccessToken:  mockToken,
			CreatedAt:    &time.Time{},
			UpdatedAt:    &time.Time{},
		}, nil)

	/* ----------------------- FindByCredentialId err case ---------------------- */

	repoMock.On("CredentialSearch", mock.Anything, grpcUrl, &playerProto.CredentialSearchReq{
		Email:    errFindByCredEmail,
		Password: password,
	}).Return(&playerProto.PlayerProfile{
		Id:        mockPlayerIdFindByCredErr,
		Email:     successEmail,
		Username:  username,
		RoleCode:  0,
		CreatedAt: zeroTime,
		UpdatedAt: zeroTime,
	}, nil)

	repoMock.On("InsertOne", ctx, &authModule.Credential{
		PlayerId:     mockPlayerIdFindByCredErr,
		RoleCode:     0,
		AccessToken:  mockToken,
		RefreshToken: mockToken,
	}).Return(mockCredentialIdErr.Hex(), nil)

	repoMock.On("FindByCredentialId", ctx, mockCredentialIdErr.Hex()).
		Return(&authModule.Credential{}, errors.New("mock_err"))

	tests := []testLogin{
		{
			ctx: ctx,
			req: &authModule.PlayerLoginReq{
				Email:    successEmail,
				Password: password,
			},
			expected: &authModule.ProfileIntercepter{
				PlayerProfile: &playerModule.PlayerProfile{
					Id:        mockPlayerId,
					Email:     successEmail,
					Username:  username,
					CreatedAt: &time.Time{},
					UpdatedAt: &time.Time{},
				},
				Credential: &authModule.CredentialRes{
					Id:           mockCredentialId.Hex(),
					PlayerId:     mockPlayerId,
					RoleCode:     0,
					AccessToken:  mockToken,
					RefreshToken: mockToken,
					CreatedAt:    &time.Time{},
					UpdatedAt:    &time.Time{},
				},
			},
			isErr: false,
		},
		{
			ctx: ctx,
			req: &authModule.PlayerLoginReq{
				Email:    errEmail,
				Password: password,
			},
			expected: nil,
			isErr:    true,
		},
		{
			ctx: ctx,
			req: &authModule.PlayerLoginReq{
				Email:    errInsertEmail,
				Password: password,
			},
			expected: nil,
			isErr:    true,
		},
		{
			ctx: ctx,
			req: &authModule.PlayerLoginReq{
				Email:    errFindByCredEmail,
				Password: password,
			},
			expected: nil,
			isErr:    true,
		},
	}

	for i, _test := range tests {
		fmt.Printf("case -> %d\n", i+1)

		result, err := usecase.Login(_test.ctx, _test.req)
		if _test.isErr {
			assert.NotEmpty(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, _test.expected, result)
	}
}
