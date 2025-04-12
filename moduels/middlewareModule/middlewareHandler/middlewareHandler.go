package middlewareHandler

type Handler struct {
	MiddlewareHttp *middlewareHttp
}

func NewHandler(MiddlewareHttp *middlewareHttp) Handler {
	return Handler{
		MiddlewareHttp,
	}
}
