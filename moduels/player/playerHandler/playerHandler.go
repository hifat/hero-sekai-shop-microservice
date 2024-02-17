package playerHandler

import "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerUsecase"

type (
	playerHandler struct {
		playerUsecase playerUsecase.IPlayerUsecase
	}
)

func NewPlayerHttp(playerUsecase playerUsecase.IPlayerUsecase) playerHandler {
	return playerHandler{playerUsecase}
}
