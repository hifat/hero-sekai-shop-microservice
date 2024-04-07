package itemHandler

type Handler struct {
	ItemGrpc *itemGrpc
	ItemHttp *itemHttp
}

func NewHandler(ItemGrpc *itemGrpc, ItemHttp *itemHttp) Handler {
	return Handler{
		ItemGrpc,
		ItemHttp,
	}
}
