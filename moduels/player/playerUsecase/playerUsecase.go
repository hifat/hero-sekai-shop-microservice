package playerUsecase

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerRepository"
)

type (
	IPlayerUsecase interface{}

	playerUsecase struct {
		playerRepo playerRepository.IPlayerRepository
	}
)

func NewPlayer(playerRepo playerRepository.IPlayerRepository) IPlayerUsecase {
	return &playerUsecase{playerRepo}
}
