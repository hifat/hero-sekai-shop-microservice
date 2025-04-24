package paymentHandler

type Handler struct {
	PaymentGrpc *paymentGrpc
	PaymentHttp *paymentHttp
}

func NewHandler(PaymentGrpc *paymentGrpc, PaymentHttp *paymentHttp) Handler {
	return Handler{
		PaymentGrpc,
		PaymentHttp,
	}
}
