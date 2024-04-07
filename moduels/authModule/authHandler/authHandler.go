package authHandler

type Handler struct {
	AuthHttp *authHttp
	AuthGrpc *authGrpc
}

func NewHandler(AuthHttp *authHttp, AuthGrpc *authGrpc) Handler {
	return Handler{
		AuthHttp,
		AuthGrpc,
	}
}
