package inventoryRepository

import (
	"context"
	"encoding/json"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/queue"
)

func (r *inventoryRepository) AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *paymentModule.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"buy",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}

func (r *inventoryRepository) RemovePlayerItemRes(pctx context.Context, cfg *config.Config, req *paymentModule.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"payment",
		"sell",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}
