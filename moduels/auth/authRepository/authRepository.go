package authRepository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IAuthRepository interface{}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuth(db *mongo.Client) IAuthRepository {
	return &authRepository{db}
}

func (r *authRepository) dbConn() *mongo.Database {
	return r.db.Database("auth_db")
}
