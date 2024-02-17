package playerHandler

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerUsecase"

type (
	playerQueue struct {
		playerUsecase playerUsecase.IPlayerUsecase
	}
)

func NewPlayerQueue(playerUsecase playerUsecase.IPlayerUsecase) *playerQueue {
	return &playerQueue{playerUsecase}
}
