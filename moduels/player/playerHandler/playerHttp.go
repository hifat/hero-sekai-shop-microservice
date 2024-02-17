package playerHandler

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerUsecase"

type (
	playerHttp struct {
		playerUsecase playerUsecase.IPlayerUsecase
	}
)

func NewPlayerHttp(playerUsecase playerUsecase.IPlayerUsecase) *playerHttp {
	return &playerHttp{playerUsecase}
}
