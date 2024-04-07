package middlewareRepository

import "go.mongodb.org/mongo-driver/mongo"

type (
	IMiddlewareRepository interface{}

	middlewareRepository struct {
		db *mongo.Client
	}
)

func NewMiddleware(db *mongo.Client) IMiddlewareRepository {
	return &middlewareRepository{db}
}
