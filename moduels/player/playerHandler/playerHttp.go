package playerHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/request"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/response"
)

type (
	playerHttp struct {
		playerUsecase playerUsecase.IPlayerUsecase
	}
)

func NewPlayerHttp(playerUsecase playerUsecase.IPlayerUsecase) *playerHttp {
	return &playerHttp{playerUsecase}
}

func (h *playerHttp) Create(c echo.Context) error {
	httpCtx := request.NewHttpContext(c)

	var req player.CreatePlayerReq
	if err := httpCtx.Bind(&req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUsecase.Create(c.Request().Context(), req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessReponse(c, http.StatusOK, res)
}
