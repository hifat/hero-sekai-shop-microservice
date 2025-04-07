package authUsecase

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
)

type (
	IAuthUsecase interface {
		Login(pctx context.Context, req *authModule.PlayerLoginReq) (*authModule.ProfileIntercepter, error)
	}

	authUsecase struct {
		cfg      *config.Config
		authRepo authRepository.IAuthRepository
	}
)

func NewAuth(cfg *config.Config, authRepo authRepository.IAuthRepository) IAuthUsecase {
	return &authUsecase{
		cfg,
		authRepo,
	}
}

func (u *authUsecase) Login(pctx context.Context, req *authModule.PlayerLoginReq) (*authModule.ProfileIntercepter, error) {
	profile, err := u.authRepo.CredentialSearch(pctx, u.cfg.Grpc.PlayerUrl, &playerProto.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	accessToken := jwtauth.NewAccessToken(u.cfg.Jwt.RefreshSecretKey, u.cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: profile.RoleCode,
	})

	refreshToken := jwtauth.NewRefreshToken(u.cfg.Jwt.RefreshSecretKey, u.cfg.Jwt.RefreshDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: profile.RoleCode,
	})

	credentialId, err := u.authRepo.InsertOne(pctx, &authModule.Credential{
		PlayerId:     profile.Id,
		RoleCode:     profile.RoleCode,
		AccessToken:  accessToken.SignToken(),
		RefreshToken: refreshToken.SignToken(),
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	createdAt := utils.MustStrToTime(profile.CreatedAt)
	updatedAt := utils.MustStrToTime(profile.UpdatedAt)

	credential, err := u.authRepo.FindByCredentialId(pctx, credentialId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &authModule.ProfileIntercepter{
		PlayerProfile: &playerModule.PlayerProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		},
		Credential: &authModule.CredentialRes{
			Id:           credentialId,
			PlayerId:     credential.PlayerId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt,
			UpdatedAt:    credential.UpdatedAt,
		},
	}, nil
}
