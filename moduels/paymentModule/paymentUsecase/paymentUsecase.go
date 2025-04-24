package paymentUsecase

import (
	"context"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	paymentRepository "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/queue"
)

type (
	IPaymentUsecase interface {
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		FindItemInIds(pctx context.Context, req []*paymentModule.ItemServiceReqDatum) error
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
		return err
	}

	return nil
}

func (u *paymentUsecase) PaymentConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{u.cfg.Kafka.Url}, u.cfg.Kafka.ApiKey, u.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := u.paymentRepo.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		slog.Warn("try to set offset as 0")
		consumer, err = worker.ConsumePartition("payment", 0, 0)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return consumer, nil
}

func (u *paymentUsecase) TradingItemConsumer(pctx context.Context, key string, resCh chan<- *paymentModule.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(pctx)
	if err != nil {
		resCh <- nil
		return
	}

	slog.Info("start BuyItemConsumer")

	select {
	case err := <-consumer.Errors():
		slog.Error("BuyItemConsumer Failed: ", err.Error())
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
}

func (u *paymentUsecase) BuyItem(pctx context.Context, playerId string, req *paymentModule.ItemServiceReq) error {
	if err := u.FindItemInIds(pctx, req.Items); err != nil {
		return err
	}

	return nil
}

func (u *paymentUsecase) SellItem(pctx context.Context, playerId string, req *paymentModule.ItemServiceReq) error {
	if err := u.FindItemInIds(pctx, req.Items); err != nil {
		return err
	}

	return nil
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
