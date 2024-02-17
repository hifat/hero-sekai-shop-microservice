package middlewareHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	middlewareHandler struct {
		middlewareUsecase middlewareUsecase.IUsecase
	}
)

func NewMiddlewareHandler(db *mongo.Client) middlewareHandler {
	return middlewareHandler{db}
}
