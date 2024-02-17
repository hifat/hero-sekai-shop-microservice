package middlewareRepository

import "go.mongodb.org/mongo-driver/mongo"

type (
	IRepository interface{}

	repository struct {
		db *mongo.Client
	}
)

func New(db *mongo.Client) IRepository {
	return &repository{db}
}
