package playerHandler

type Handler struct {
	PlayerHttp            *playerHttp
	PlayerTransactionHttp *playerTransactionHttp
	PlayerGrpc            *playerGrpc
	PlayerQueue           *playerQueue
}

func NewHandler(playerHttp *playerHttp, playerTransactionHttp *playerTransactionHttp, playerGrpc *playerGrpc, playerQueue *playerQueue) Handler {
	return Handler{
		playerHttp,
		playerTransactionHttp,
		playerGrpc,
		playerQueue,
	}
}
