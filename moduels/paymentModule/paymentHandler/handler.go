package paymentHandler

type Handler struct {
	PaymentGrpc  *paymentGrpc
	paymentHttp  *paymentHttp
	PaymentQueue *paymentQueue
}

func NewHandler(PaymentGrpc *paymentGrpc, paymentHttp *paymentHttp, PaymentQueue *paymentQueue) Handler {
	return Handler{
		PaymentGrpc,
		paymentHttp,
		PaymentQueue,
	}
}
