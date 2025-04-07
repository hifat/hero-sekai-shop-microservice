package authModule

import (
	"time"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
)

type (
	PlayerLoginReq struct {
		Email    string `json:"email" form:"email" validate:"required,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=32"`
	}

	RefreshTokenReq struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required,max=500"`
	}

	InsertPlayerRole struct {
		PlayerId string `json:"player_id" validate:"required"`
		RoleCode []int  `json:"role_id" validate:"required"`
	}

	ProfileIntercepter struct {
		*playerModule.PlayerProfile
		Credential *CredentialRes `json:"credential"`
	}

	CredentialRes struct {
		Id           string     `json:"_id"`
		PlayerId     string     `json:"player_id"`
		RoleCode     int32      `json:"role_code"`
		AccessToken  string     `json:"access_token"`
		RefreshToken string     `json:"refresh_token"`
		CreatedAt    *time.Time `json:"created_at"`
		UpdatedAt    *time.Time `json:"updated_at"`
	}
)
