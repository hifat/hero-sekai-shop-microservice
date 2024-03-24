package playerUsecase

import (
	"context"

	"github.com/jinzhu/copier"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerError"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type (
	IPlayerUsecase interface {
		Create(pctx context.Context, req player.CreatePlayerReq) (string, error)
	}

	playerUsecase struct {
		playerRepo playerRepository.IPlayerRepository
	}
)

func NewPlayer(playerRepo playerRepository.IPlayerRepository) IPlayerUsecase {
	return &playerUsecase{playerRepo}
}

func (u *playerUsecase) Create(pctx context.Context, req player.CreatePlayerReq) (string, error) {
	isUsernameExists, err := u.playerRepo.ExistsByField(pctx, "username", req.Username)
	if err != nil {
		return "", err
	}

	if isUsernameExists {
		return "", playerError.ErrDuplicateUsername
	}

	isEmailExists, err := u.playerRepo.ExistsByField(pctx, "email", req.Email)
	if err != nil {
		return "", err
	}

	if isEmailExists {
		return "", playerError.ErrDuplicateEmail
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	req.Password = string(hashPassword)

	var newPlayer player.Player
	if err := copier.Copy(&newPlayer, &req); err != nil {
		logger.Error(err)
		return "", err
	}

	newPlayer.PlayerRoles = []player.PlayerRole{
		{
			RoleTitle: "Player",
			RoleCode:  0,
		},
	}

	playerId, err := u.playerRepo.Create(pctx, newPlayer)
	if err != nil {
		return "", nil
	}

	return playerId, nil
}
