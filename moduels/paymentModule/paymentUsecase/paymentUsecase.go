package paymentUsecase

import (
	"context"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	paymentRepository "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/queue"
)

type (
	IPaymentUsecase interface {
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		FindItemInIds(pctx context.Context, req []*paymentModule.ItemServiceReqDatum) error
		PaymentConsumer(pctx context.Context) (sarama.PartitionConsumer, error)
		TradingItemConsumer(pctx context.Context, key string, resCh chan<- *paymentModule.PaymentTransferRes)
		BuyItem(pctx context.Context, playerId string, req *paymentModule.ItemServiceReq) ([]*paymentModule.PaymentTransferRes, error)
		SellItem(pctx context.Context, playerId string, req *paymentModule.ItemServiceReq) ([]*paymentModule.PaymentTransferRes, error)
	}

	paymentUsecase struct {
		cfg         *config.Config
		paymentRepo paymentRepository.IPaymentRepository
	}
)

func NewPayment(cfg *config.Config, paymentRepo paymentRepository.IPaymentRepository) IPaymentUsecase {
	return &paymentUsecase{
		cfg,
		paymentRepo,
	}
}

func (u *paymentUsecase) GetOffset(pctx context.Context) (int64, error) {
	offset, err := u.paymentRepo.GetOffset(pctx)
	if err != nil {
		return -1, err
	}

	return offset, nil
}
func (u *paymentUsecase) UpsertOffset(pctx context.Context, offset int64) error {
	err := u.paymentRepo.UpsertOffset(pctx, offset)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (u *paymentUsecase) PaymentConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{u.cfg.Kafka.Url}, u.cfg.Kafka.ApiKey, u.cfg.Kafka.Secret)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	offset, err := u.paymentRepo.GetOffset(pctx)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		logger.Warn("try to set offset as 0")
		consumer, err = worker.ConsumePartition("payment", 0, 0)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	return consumer, nil
}

func (u *paymentUsecase) TradingItemConsumer(pctx context.Context, key string, resCh chan<- *paymentModule.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(pctx)
	if err != nil {
		resCh <- nil
		return
	}
	defer consumer.Close()

	logger.Info("start BuyItemConsumer...")

	select {
	case err := <-consumer.Errors():
		slog.Error("BuyItemConsumer Failed: ", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			u.UpsertOffset(pctx, msg.Offset+1)

			req := new(paymentModule.PaymentTransferRes)
			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				resCh <- nil
				return
			}

			resCh <- req
			log.Printf("BuyItemConsumer | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
		}
	}

	resCh <- nil
}

func (u *paymentUsecase) BuyItem(pctx context.Context, playerId string, req *paymentModule.ItemServiceReq) ([]*paymentModule.PaymentTransferRes, error) {
	if err := u.FindItemInIds(pctx, req.Items); err != nil {
		logger.Error(err)
		return nil, err
	}

	stage1 := make([]*paymentModule.PaymentTransferRes, 0, len(req.Items))
	for _, item := range req.Items {
		if err := u.paymentRepo.DockedPlayerMoney(pctx, u.cfg, &playerModule.CreatePlayerTransactionReq{
			PlayerId: playerId,
			Amount:   -item.Price,
		}); err != nil {
			logger.Error(err)
			return nil, err
		}

		resCh := make(chan *paymentModule.PaymentTransferRes)

		go u.TradingItemConsumer(pctx, "buy", resCh)

		res := <-resCh
		if res != nil {
			stage1 = append(stage1, &paymentModule.PaymentTransferRes{
				InventoryId:   res.InventoryId,
				TransactionId: res.TransactionId,
				PlayerId:      res.PlayerId,
				ItemId:        res.ItemId,
				Amount:        res.Amount,
				Error:         res.Error,
			})
		}
	}

	for _, s1 := range stage1 {
		if s1.Error != "" {
			for _, ss1 := range stage1 {
				logger.Error(s1.Error)
				u.paymentRepo.RollbackDockedPlayerMoney(pctx, u.cfg, &playerModule.RollbackPlayerTransactionReq{
					TransactionId: ss1.TransactionId,
				})
			}
		}
	}

	stage2 := make([]*paymentModule.PaymentTransferRes, 0, len(stage1))
	for _, s1 := range stage1 {
		if err := u.paymentRepo.AddPlayerItem(pctx, u.cfg, &inventoryModule.UpdateInventoryReq{
			PlayerId: playerId,
			ItemId:   s1.ItemId,
		}); err != nil {
			logger.Error(err)
			return nil, err
		}

		resCh := make(chan *paymentModule.PaymentTransferRes)

		go u.TradingItemConsumer(pctx, "buy", resCh)

		res := <-resCh
		if res != nil {
			stage2 = append(stage2, &paymentModule.PaymentTransferRes{
				InventoryId:   res.InventoryId,
				TransactionId: res.TransactionId,
				PlayerId:      res.PlayerId,
				ItemId:        res.ItemId,
				Amount:        res.Amount,
				Error:         res.Error,
			})
		}
	}

	for _, s2 := range stage2 {
		if s2.Error != "" {
			logger.Error(s2.Error)
			for _, ss2 := range stage2 {
				u.paymentRepo.RollbackAddPlayerItem(pctx, u.cfg, &inventoryModule.RollbackPlayerInventoryReq{
					InventoryId: ss2.TransactionId,
				})
			}

			for _, ss2 := range stage2 {
				u.paymentRepo.RollbackDockedPlayerMoney(pctx, u.cfg, &playerModule.RollbackPlayerTransactionReq{
					TransactionId: ss2.TransactionId,
				})
			}
		}
	}

	return stage1, nil
}

func (u *paymentUsecase) SellItem(pctx context.Context, playerId string, req *paymentModule.ItemServiceReq) ([]*paymentModule.PaymentTransferRes, error) {
	if err := u.FindItemInIds(pctx, req.Items); err != nil {
		logger.Error(err)
		return nil, err
	}

	return nil, nil
}

func (u *paymentUsecase) FindItemInIds(pctx context.Context, req []*paymentModule.ItemServiceReqDatum) error {
	itemData, err := u.paymentRepo.FindItemInIds(pctx, u.cfg.Grpc.ItemUrl, &itemProto.FindItemsInIdsReq{
		Ids: func() []string {
			setId := map[string]struct{}{}
			itemIds := make([]string, 0, len(req))
			for _, v := range req {
				if _, ok := setId[v.ItemId]; !ok {
					itemIds = append(itemIds, v.ItemId)
				}
			}

			return itemIds
		}(),
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	itemMap := make(map[string]*itemModule.ItemShowCase)
	for _, v := range itemData.Items {
		itemMap[v.Id] = &itemModule.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			Damage:   v.Damage,
			ImageUrl: v.ImageUrl,
		}
	}

	for i := range req {
		if _, ok := itemMap[req[i].ItemId]; ok {
			req[i].Price = itemMap[req[i].ItemId].Price
		}
	}

	return nil
}
