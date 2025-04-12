package inventoryHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/request"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/response"
)

type (
	inventoryHttp struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryHttp(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *inventoryHttp {
	return &inventoryHttp{cfg, inventoryUsecase}
}

func (h *inventoryHttp) FindPlayerItems(c echo.Context) error {
	ctx := c.Request().Context()
	wrapper := request.NewHttpContext(c)

	credential := c.Get("credential").(*jwtauth.AuthMapClaims)

	req := new(inventoryModule.InventorySearchReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.inventoryUsecase.FindPlayerItems(ctx, h.cfg.Paginate.InventoryNextPageBasedUrl, credential.PlayerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
