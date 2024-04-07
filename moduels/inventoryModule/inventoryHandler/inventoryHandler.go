package inventoryHandler

type Handler struct {
	InventoryHttp  *inventoryHttp
	InventoryGrpc  *inventoryGrpc
	InventoryQueue *inventoryQueue
}

func NewHandler(InventoryHttp *inventoryHttp, InventoryGrpc *inventoryGrpc, InventoryQueue *inventoryQueue) Handler {
	return Handler{
		InventoryHttp,
		InventoryGrpc,
		InventoryQueue,
	}
}
