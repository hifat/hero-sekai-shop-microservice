package paymentRepository

import (
	"context"
	"encoding/json"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
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

func (r *paymentRepository) AddPlayerItem(pctx context.Context, cfg *config.Config, req *inventoryModule.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"buy",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventoryModule.RollbackPlayerInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"radd",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) RemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventoryModule.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"sell",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventoryModule.RollbackPlayerInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"rremove",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) AddPlayerMoney(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"sell",
		reqInBytes,
	); err != nil {
		return err
	}

	return nil
}
