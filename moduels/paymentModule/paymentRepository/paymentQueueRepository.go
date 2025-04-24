package paymentRepository

import (
	"context"
	"encoding/json"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/queue"
)

func (r *paymentRepository) DockedPlayerMoney(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"buy",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) RollbackDockedPlayerMoney(pctx context.Context, cfg *config.Config, req *playerModule.RollbackPlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"rtransaction",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}
