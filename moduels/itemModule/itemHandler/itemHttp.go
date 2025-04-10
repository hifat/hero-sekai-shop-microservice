package itemHandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/request"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/response"
)

type (
	itemHttp struct {
		cfg         *config.Config
		itemUsecase itemUsecase.IItemUsecase
	}
)

func NewItemHttp(cfg *config.Config, itemUsecase itemUsecase.IItemUsecase) *itemHttp {
	return &itemHttp{cfg, itemUsecase}
}

func (h *itemHttp) Create(c echo.Context) error {
	ctx := c.Request().Context()
	wrapper := request.NewHttpContext(c)
	req := new(itemModule.CreateItemReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.Create(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *itemHttp) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	itemId := c.Param("item_id")

	res, err := h.itemUsecase.FindById(ctx, itemId)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *itemHttp) Find(c echo.Context) error {
	ctx := c.Request().Context()
	wrapper := request.NewHttpContext(c)
	req := new(itemModule.ItemSearchReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.Find(ctx, h.cfg.Paginate.ItemNextPageBasedUrl, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}
