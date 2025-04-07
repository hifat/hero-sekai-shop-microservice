package authUsecase

import (
	"context"
	"time"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authProto"
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
		RefreshToken(pctx context.Context, req *authModule.RefreshTokenReq) (*authModule.ProfileIntercepter, error)
		Logout(pctx context.Context, credentialId string) (int64, error)
		AccessTokenSearch(pctx context.Context, accessToken string) (*authProto.AccessTokenSearchRes, error)
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

func (u *authUsecase) RefreshToken(pctx context.Context, req *authModule.RefreshTokenReq) (*authModule.ProfileIntercepter, error) {
	claims, err := jwtauth.ParseToken(u.cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	profile, err := u.authRepo.FindOnePlayerProfileToRefresh(pctx, u.cfg.Grpc.PlayerUrl, &playerProto.FindOnePlayerProfileToRefreshReq{
		PlayerId: claims.PlayerId,
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	accessToken := jwtauth.NewAccessToken(u.cfg.Jwt.RefreshSecretKey, u.cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: profile.RoleCode,
	})

	refreshToken := jwtauth.ReloadToken(u.cfg.Jwt.RefreshSecretKey, time.Duration(claims.ExpiresAt.Unix()), &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: profile.RoleCode,
	})

	if err := u.authRepo.UpdateRefreshToken(pctx, req.CredentialId, &authModule.UpdateRefreshTokenReq{
		PlayerId:     profile.Id,
		AccessToken:  accessToken.SignToken(),
		RefreshToken: refreshToken,
		UpdatedAt:    time.Now(),
	}); err != nil {
		logger.Error(err)
		return nil, err
	}

	credential, err := u.authRepo.FindByCredentialId(pctx, req.CredentialId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	createdAt := utils.MustStrToTime(profile.CreatedAt)
	updatedAt := utils.MustStrToTime(profile.UpdatedAt)

	return &authModule.ProfileIntercepter{
		PlayerProfile: &playerModule.PlayerProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		},
		Credential: &authModule.CredentialRes{
			Id:           req.CredentialId,
			PlayerId:     credential.PlayerId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt,
			UpdatedAt:    credential.UpdatedAt,
		},
	}, nil
}

func (u *authUsecase) Logout(pctx context.Context, credentialId string) (int64, error) {
	amount, err := u.authRepo.DeleteByCredentialId(pctx, credentialId)
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	return amount, err
}

func (u *authUsecase) AccessTokenSearch(pctx context.Context, accessToken string) (*authProto.AccessTokenSearchRes, error) {
	credential, err := u.authRepo.FindByAccessToken(pctx, accessToken)
	if err != nil {
		logger.Error(err)
		return &authProto.AccessTokenSearchRes{
			IsValid: false,
		}, err
	}

	if credential == nil {
		return &authProto.AccessTokenSearchRes{
			IsValid: false,
		}, nil
	}

	return &authProto.AccessTokenSearchRes{
		IsValid: true,
	}, nil
}
