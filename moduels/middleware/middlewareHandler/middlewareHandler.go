package middlewareHandler

type Handler struct {
	middlewareHttp *middlewareHttp
}

func NewHandler(middlewareHttp *middlewareHttp) Handler {
	return Handler{
		middlewareHttp,
	}
}
