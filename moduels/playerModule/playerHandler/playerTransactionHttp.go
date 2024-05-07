package playerHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/request"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/response"
)

type (
	playerTransactionHttp struct {
		playerUsecase playerUsecase.IPlayerTransactionUsecase
	}
)

func NewPlayerTransactionHttp(playerUsecase playerUsecase.IPlayerTransactionUsecase) *playerTransactionHttp {
	return &playerTransactionHttp{playerUsecase}
}

func (h *playerTransactionHttp) AddMoney(c echo.Context) error {
	httpCtx := request.NewHttpContext(c)

	var req playerModule.CreatePlayerTransactionReq
	if err := httpCtx.Bind(&req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUsecase.AddMoney(c.Request().Context(), req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *playerTransactionHttp) GetSavingAccount(c echo.Context) error {
	playerId := c.Param("player_id")

	res, err := h.playerUsecase.GetSavingAccount(c.Request().Context(), playerId)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
