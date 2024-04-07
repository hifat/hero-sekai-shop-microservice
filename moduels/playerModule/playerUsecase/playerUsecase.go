package playerUsecase

import (
	"context"

	"github.com/jinzhu/copier"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerError"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type (
	IPlayerUsecase interface {
		Create(pctx context.Context, req playerModule.CreatePlayerReq) (*playerModule.PlayerProfile, error)
		GetProfile(pctx context.Context, playerId string) (*playerModule.PlayerProfile, error)
	}

	playerUsecase struct {
		playerRepo playerRepository.IPlayerRepository
	}
)

func NewPlayer(playerRepo playerRepository.IPlayerRepository) IPlayerUsecase {
	return &playerUsecase{playerRepo}
}

func (u *playerUsecase) Create(pctx context.Context, req playerModule.CreatePlayerReq) (*playerModule.PlayerProfile, error) {
	isUsernameExists, err := u.playerRepo.ExistsByField(pctx, "username", req.Username)
	if err != nil {
		return nil, err
	}

	if isUsernameExists {
		return nil, playerError.ErrDuplicateUsername
	}

	isEmailExists, err := u.playerRepo.ExistsByField(pctx, "email", req.Email)
	if err != nil {
		return nil, err
	}

	if isEmailExists {
		return nil, playerError.ErrDuplicateEmail
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	req.Password = string(hashPassword)

	var newPlayer playerModule.Player
	if err := copier.Copy(&newPlayer, &req); err != nil {
		logger.Error(err)
		return nil, err
	}

	newPlayer.PlayerRoles = []playerModule.PlayerRole{
		{
			RoleTitle: "Player",
			RoleCode:  0,
		},
	}

	playerId, err := u.playerRepo.Create(pctx, newPlayer)
	if err != nil {
		return nil, nil
	}

	return u.GetProfile(pctx, playerId)
}

func (u *playerUsecase) GetProfile(pctx context.Context, playerId string) (*playerModule.PlayerProfile, error) {
	getPlayer, err := u.playerRepo.FirstByField(pctx, "_id", playerId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var res playerModule.PlayerProfile
	if err := copier.Copy(&res, &getPlayer); err != nil {
		logger.Error(err)
		return nil, err
	}

	res.Id = getPlayer.Id.Hex()

	return &res, nil
}
