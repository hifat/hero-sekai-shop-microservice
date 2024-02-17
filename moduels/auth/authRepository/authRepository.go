package authRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IAuthRepository interface{}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) IAuthRepository {
	return &authRepository{db}
}

func (r *authRepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}
