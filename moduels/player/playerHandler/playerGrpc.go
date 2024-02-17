package playerHandler

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerUsecase"

type (
	playerGrpc struct {
		playerUsecase playerUsecase.IPlayerUsecase
	}
)

func NewPlayerGrpc(playerUsecase playerUsecase.IPlayerUsecase) *playerGrpc {
	return &playerGrpc{playerUsecase}
}
