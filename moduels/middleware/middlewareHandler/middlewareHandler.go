package middlewareHandler

import (
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/middleware/middlewareUsecase"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	middlewareHttpHandler struct {
		middlewareUsecase middlewareUsecase.IMiddlewareUsecase
	}
)

func NewHttp(db *mongo.Client) middlewareHttpHandler {
	return middlewareHttpHandler{db}
}
