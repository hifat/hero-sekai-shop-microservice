package paymentandler

type Handler struct {
	PaymentGrpc  *paymentGrpc
	paymentttp   *paymentttp
	PaymentQueue *paymentQueue
}

func NewHandler(PaymentGrpc *paymentGrpc, paymentttp *paymentttp, PaymentQueue *paymentQueue) Handler {
	return Handler{
		PaymentGrpc,
		paymentttp,
		PaymentQueue,
	}
}
