package paymentUsecase

import paymentRepository "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/payment/paymenthRepository"

type (
	IPaymentUsecase interface{}

	paymentUsecase struct {
		paymentRepo paymentRepository.IPaymentRepository
	}
)

func NewPayment(paymentRepo paymentRepository.IPaymentRepository) IPaymentUsecase {
	return &paymentUsecase{paymentRepo}
}
