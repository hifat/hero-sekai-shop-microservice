package inventoryHandler

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/inventoryModule/inventoryUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/queue"
)

type (
	inventoryQueue struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.IInventoryUsecase
	}
)

func NewInventoryQueue(cfg *config.Config, inventoryUsecase inventoryUsecase.IInventoryUsecase) *inventoryQueue {
	return &inventoryQueue{cfg, inventoryUsecase}
}

func (h *inventoryQueue) inventoryConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.inventoryUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("inventory", 0, offset)
	if err != nil {
		logger.Warn("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("inventory", 0, 0)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	return consumer, nil
}

func (h *inventoryQueue) AddPlayerItem() {
	pctx := context.Background()

	consumer, err := h.inventoryConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	slog.Info("start AddPlayerItem...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Handle consumer is nil
	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("AddPlayerItem Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.inventoryUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(inventoryModule.UpdateInventoryReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					logger.Error(err)
					continue
				}

				h.inventoryUsecase.AddPlayerItemRes(pctx, h.cfg, req)

				log.Printf("AddPlayerItem | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			logger.Info("Stop AddPlayerItem...")
			return
		}
	}
}

func (h *inventoryQueue) RollbackAddPlayerItem() {
	pctx := context.Background()

	consumer, err := h.inventoryConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	slog.Info("start RollbackAddPlayerItem...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Handle consumer is nil
	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("RollbackAddPlayerItem Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "radd" {
				h.inventoryUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(inventoryModule.RollbackPlayerInventoryReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					logger.Error(err)
					continue
				}

				h.inventoryUsecase.RollbackAddPlayerItem(pctx, h.cfg, req)

				log.Printf("RollbackAddPlayerItem | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			logger.Info("Stop RollbackAddPlayerItem...")
			return
		}
	}
}

func (h *inventoryQueue) RemovePlayerItem() {
	pctx := context.Background()

	consumer, err := h.inventoryConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	slog.Info("start RemovePlayerItem...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Handle consumer is nil
	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("RemovePlayerItem Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.inventoryUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(inventoryModule.UpdateInventoryReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					logger.Error(err)
					continue
				}

				h.inventoryUsecase.RemovePlayerItemRes(pctx, h.cfg, req)

				log.Printf("RemovePlayerItem | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			logger.Info("Stop RemovePlayerItem...")
			return
		}
	}
}

func (h *inventoryQueue) RollbackRemovePlayerItem() {
	pctx := context.Background()

	consumer, err := h.inventoryConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	slog.Info("start RollbackRemovePlayerItem...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Handle consumer is nil
	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("RollbackRemovePlayerItem Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rremove" {
				h.inventoryUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(inventoryModule.RollbackPlayerInventoryReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					logger.Error(err)
					continue
				}

				h.inventoryUsecase.RollbackRemovePlayerItem(pctx, h.cfg, req)

				log.Printf("RollbackRemovePlayerItem | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			logger.Info("Stop RollbackRemovePlayerItem...")
			return
		}
	}
}
