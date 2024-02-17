package playerHandler

type Handler struct {
	PlayerHttp  *playerHttp
	PlayerGrpc  *playerGrpc
	PlayerQueue *playerQueue
}

func NewHandler(playerHttp *playerHttp, playerGrpc *playerGrpc, playerQueue *playerQueue) Handler {
	return Handler{
		playerHttp,
		playerGrpc,
		playerQueue,
	}
}
