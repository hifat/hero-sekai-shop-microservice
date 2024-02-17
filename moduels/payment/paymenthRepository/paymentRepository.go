package paymentRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IPaymentRepository interface{}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPayment(db *mongo.Client) IPaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) dbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}
