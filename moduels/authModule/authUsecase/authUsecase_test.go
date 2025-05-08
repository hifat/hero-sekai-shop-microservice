package authUsecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/whydoweneedtest"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthTestSuite struct {
	suite.Suite
	cfg      *config.Config
	repoMock *authRepository.AuthRepositoryMock
	usecase  IAuthUsecase
	testData struct {
		grpcURL     string
		emails      map[string]string
		credentials map[string]string
		tokens      map[string]string
		playerInfo  map[string]string
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupTest() {
	s.cfg = whydoweneedtest.NewTestConfig()
	s.repoMock = new(authRepository.AuthRepositoryMock)
	s.usecase = NewAuth(s.cfg, s.repoMock)

	// Initialize test data
	s.testData = struct {
		grpcURL     string
		emails      map[string]string
		credentials map[string]string
		tokens      map[string]string
		playerInfo  map[string]string
	}{
		grpcURL: s.cfg.Grpc.PlayerUrl,
		emails: map[string]string{
			"success":       gofakeit.Email(),
			"error":         gofakeit.Email(),
			"insertError":   gofakeit.Email(),
			"findByCredErr": gofakeit.Email(),
		},
		credentials: map[string]string{
			"valid":   primitive.NewObjectID().Hex(),
			"invalid": primitive.NewObjectID().Hex(),
		},
		tokens: map[string]string{
			"mock": gofakeit.UUID(),
		},
		playerInfo: map[string]string{
			"username":        gofakeit.Username(),
			"password":        gofakeit.Password(true, true, true, true, false, 10),
			"validId":         primitive.NewObjectID().Hex(),
			"invalidId":       "player_id_err",
			"findByCredErrId": "player_id_find_by_cred_err",
		},
	}
}

func (s *AuthTestSuite) TearDownTest() {
	s.repoMock.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestLogin() {
	tests := []struct {
		name     string
		ctx      context.Context
		req      *authModule.PlayerLoginReq
		setup    func()
		expected *authModule.ProfileIntercepter
		wantErr  bool
	}{
		{
			name: "success - valid login credentials",
			ctx:  context.Background(),
			req: &authModule.PlayerLoginReq{
				Email:    s.testData.emails["success"],
				Password: s.testData.playerInfo["password"],
			},
			setup: func() {
				s.setupSuccessfulLoginMocks()
			},
			expected: s.getExpectedSuccessProfile(),
			wantErr:  false,
		},
		{
			name: "error - invalid credentials",
			ctx:  context.Background(),
			req: &authModule.PlayerLoginReq{
				Email:    s.testData.emails["error"],
				Password: s.testData.playerInfo["password"],
			},
			setup: func() {
				s.repoMock.On("CredentialSearch", mock.Anything, s.testData.grpcURL, &playerProto.CredentialSearchReq{
					Email:    s.testData.emails["error"],
					Password: s.testData.playerInfo["password"],
				}).Return(&playerProto.PlayerProfile{}, errors.New("mock_error"))
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "error - failed to insert credential",
			ctx:  context.Background(),
			req: &authModule.PlayerLoginReq{
				Email:    s.testData.emails["insertError"],
				Password: s.testData.playerInfo["password"],
			},
			setup: func() {
				s.repoMock.On("CredentialSearch", mock.Anything, s.testData.grpcURL, &playerProto.CredentialSearchReq{
					Email:    s.testData.emails["insertError"],
					Password: s.testData.playerInfo["password"],
				}).Return(&playerProto.PlayerProfile{
					Id:        s.testData.playerInfo["invalidId"],
					Email:     s.testData.emails["insertError"],
					Username:  s.testData.playerInfo["username"],
					RoleCode:  0,
					CreatedAt: time.Time{}.String(),
					UpdatedAt: time.Time{}.String(),
				}, nil)

				s.repoMock.On("InsertOne", context.Background(), &authModule.Credential{
					PlayerId:     s.testData.playerInfo["invalidId"],
					RoleCode:     0,
					AccessToken:  s.testData.tokens["mock"],
					RefreshToken: s.testData.tokens["mock"],
				}).Return("", errors.New("mock_err"))
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "error - failed to find credential by ID",
			ctx:  context.Background(),
			req: &authModule.PlayerLoginReq{
				Email:    s.testData.emails["findByCredErr"],
				Password: s.testData.playerInfo["password"],
			},
			setup: func() {
				s.repoMock.On("CredentialSearch", mock.Anything, s.testData.grpcURL, &playerProto.CredentialSearchReq{
					Email:    s.testData.emails["findByCredErr"],
					Password: s.testData.playerInfo["password"],
				}).Return(&playerProto.PlayerProfile{
					Id:        s.testData.playerInfo["findByCredErrId"],
					Email:     s.testData.emails["findByCredErr"],
					Username:  s.testData.playerInfo["username"],
					RoleCode:  0,
					CreatedAt: time.Time{}.String(),
					UpdatedAt: time.Time{}.String(),
				}, nil)

				s.repoMock.On("InsertOne", context.Background(), &authModule.Credential{
					PlayerId:     s.testData.playerInfo["findByCredErrId"],
					RoleCode:     0,
					AccessToken:  s.testData.tokens["mock"],
					RefreshToken: s.testData.tokens["mock"],
				}).Return(s.testData.credentials["invalid"], nil)

				s.repoMock.On("FindByCredentialId", context.Background(), s.testData.credentials["invalid"]).
					Return(&authModule.Credential{}, errors.New("mock_err"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setup != nil {
				tt.setup()
			}

			result, err := s.usecase.Login(tt.ctx, tt.req)
			if tt.wantErr {
				s.Error(err)
				s.Nil(result)
				return
			}

			s.NoError(err)
			s.Equal(tt.expected, result)
		})
	}
}

func (s *AuthTestSuite) setupSuccessfulLoginMocks() {
	s.repoMock.On("CredentialSearch", mock.Anything, s.testData.grpcURL, &playerProto.CredentialSearchReq{
		Email:    s.testData.emails["success"],
		Password: s.testData.playerInfo["password"],
	}).Return(&playerProto.PlayerProfile{
		Id:        s.testData.playerInfo["validId"],
		Email:     s.testData.emails["success"],
		Username:  s.testData.playerInfo["username"],
		RoleCode:  0,
		CreatedAt: time.Time{}.String(),
		UpdatedAt: time.Time{}.String(),
	}, nil)

	s.repoMock.On("NewAccessToken", s.cfg.Jwt, mock.AnythingOfType("*playerProto.PlayerProfile")).
		Return(s.testData.tokens["mock"])

	s.repoMock.On("NewRefreshToken", s.cfg.Jwt, mock.AnythingOfType("*playerProto.PlayerProfile")).
		Return(s.testData.tokens["mock"])

	s.repoMock.On("InsertOne", context.Background(), &authModule.Credential{
		PlayerId:     s.testData.playerInfo["validId"],
		RoleCode:     0,
		AccessToken:  s.testData.tokens["mock"],
		RefreshToken: s.testData.tokens["mock"],
	}).Return(s.testData.credentials["valid"], nil)

	s.repoMock.On("FindByCredentialId", context.Background(), s.testData.credentials["valid"]).
		Return(&authModule.Credential{
			Id:           primitive.NewObjectID(),
			PlayerId:     s.testData.playerInfo["validId"],
			RoleCode:     0,
			RefreshToken: s.testData.tokens["mock"],
			AccessToken:  s.testData.tokens["mock"],
			CreatedAt:    &time.Time{},
			UpdatedAt:    &time.Time{},
		}, nil)
}

func (s *AuthTestSuite) getExpectedSuccessProfile() *authModule.ProfileIntercepter {
	return &authModule.ProfileIntercepter{
		PlayerProfile: &playerModule.PlayerProfile{
			Id:        s.testData.playerInfo["validId"],
			Email:     s.testData.emails["success"],
			Username:  s.testData.playerInfo["username"],
			CreatedAt: &time.Time{},
			UpdatedAt: &time.Time{},
		},
		Credential: &authModule.CredentialRes{
			Id:           s.testData.credentials["valid"],
			PlayerId:     s.testData.playerInfo["validId"],
			RoleCode:     0,
			AccessToken:  s.testData.tokens["mock"],
			RefreshToken: s.testData.tokens["mock"],
			CreatedAt:    &time.Time{},
			UpdatedAt:    &time.Time{},
		},
	}
}
