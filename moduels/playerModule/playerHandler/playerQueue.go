package playerHandler

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerUsecase"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/queue"
)

type (
	playerQueue struct {
		cfg                      *config.Config
		playerUsecase            playerUsecase.IPlayerUsecase
		playerTransactionUsecase playerUsecase.IPlayerTransactionUsecase
	}
)

func NewPlayerQueue(cfg *config.Config, playerUsecase playerUsecase.IPlayerUsecase, playerTransactionUsecase playerUsecase.IPlayerTransactionUsecase) *playerQueue {
	return &playerQueue{
		cfg,
		playerUsecase,
		playerTransactionUsecase,
	}
}

func (h *playerQueue) PlayerConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.playerTransactionUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("player", 0, offset)
	if err != nil {
		logger.Warn("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("player", 0, 0)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	return consumer, nil
}

func (h *playerQueue) DockedPlayerMoney() {
	pctx := context.Background()

	consumer, err := h.PlayerConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	slog.Info("start DockedPlayerMoney...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Handle consumer is nil
	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("DockedPlayerMoney Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.playerTransactionUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(playerModule.CreatePlayerTransactionReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					logger.Error(err)
					continue
				}

				h.playerTransactionUsecase.DockedPlayerMoneyRes(pctx, h.cfg, req)

				log.Printf("DockedPlayerMoney | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			logger.Info("Stop DockedPlayerMoney...")
			return
		}
	}
}

func (h *playerQueue) RollbackPlayerTransaction() {
	pctx := context.Background()

	consumer, err := h.PlayerConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("start RollbackPlayerTransaction...")

	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("RollbackPlayerTransaction Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "transaction" {
				h.playerTransactionUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(playerModule.RollbackPlayerTransactionReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerTransactionUsecase.RollbackPlayerTransaction(pctx, req)

				log.Printf("RollbackPlayerTransaction | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		}
	}
}

func (h *playerQueue) AddPlayerMoney() {
	pctx := context.Background()

	consumer, err := h.PlayerConsumer(pctx)
	if err != nil {
		logger.Error(err)
		return
	}

	slog.Info("start AddPlayerMoney...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Handle consumer is nil
	for {
		select {
		case err := <-consumer.Errors():
			logger.Error("AddPlayerMoney Failed: " + err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.playerTransactionUsecase.UpsertOffset(pctx, msg.Offset+1)

				req := new(playerModule.CreatePlayerTransactionReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					logger.Error(err)
					continue
				}

				h.playerTransactionUsecase.AddPlayerMoneyRes(pctx, h.cfg, req)

				log.Printf("AddPlayerMoney | topic(%s) | Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			logger.Info("Stop AddPlayerMoney...")
			return
		}
	}
}
