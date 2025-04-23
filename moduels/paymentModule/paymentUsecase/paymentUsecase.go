package paymentUsecase

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/itemModule/itemProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	paymentRepository "gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule/paymentRepository"
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
