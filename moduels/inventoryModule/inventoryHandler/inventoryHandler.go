package inventoryHandler

type Handler struct {
	InventoryHttp  *inventoryHttp
	InventoryQueue *inventoryQueue
}

func NewHandler(InventoryHttp *inventoryHttp, InventoryQueue *inventoryQueue) Handler {
	return Handler{
		InventoryHttp,
		InventoryQueue,
	}
}
