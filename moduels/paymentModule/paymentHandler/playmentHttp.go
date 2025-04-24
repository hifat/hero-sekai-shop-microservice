package paymentHandler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/request"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/response"
)

type (
	paymentHttp struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.IPaymentUsecase
	}
)

func NewPaymentHttp(cfg *config.Config, paymentUsecase paymentUsecase.IPaymentUsecase) *paymentHttp {
	return &paymentHttp{cfg, paymentUsecase}
}

func (h *paymentHttp) BuyItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewHttpContext(c)

	credential := c.Get("credential").(*jwtauth.AuthMapClaims)

	req := &paymentModule.ItemServiceReq{
		Items: make([]*paymentModule.ItemServiceReqDatum, 0),
	}

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.paymentUsecase.BuyItem(ctx, credential.PlayerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *paymentHttp) SellItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.NewHttpContext(c)

	credential := c.Get("credential").(*jwtauth.AuthMapClaims)

	req := &paymentModule.ItemServiceReq{
		Items: make([]*paymentModule.ItemServiceReqDatum, 0),
	}

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.paymentUsecase.SellItem(ctx, credential.PlayerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
