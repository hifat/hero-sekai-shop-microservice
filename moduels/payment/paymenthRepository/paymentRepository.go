package paymentRepository

import "go.mongodb.org/mongo-driver/mongo"

type (
	IPaymentRepository interface{}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPayment(db *mongo.Client) IPaymentRepository {
	return &paymentRepository{db}
}
