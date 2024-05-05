package playerUsecase

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
)

type (
	IPlayerTransactionUsecase interface {
		AddMoney(pctx context.Context, req playerModule.CreatePlayerTransactionReq) error
	}

	playerTransactionUsecase struct {
		playerTransactionRepo playerRepository.IPlayerTransactionRepository
	}
)

func NewPlayerTransaction(playerTransactionRepo playerRepository.IPlayerTransactionRepository) IPlayerTransactionUsecase {
	return &playerTransactionUsecase{playerTransactionRepo}
}

func (u *playerTransactionUsecase) AddMoney(pctx context.Context, req playerModule.CreatePlayerTransactionReq) error {
	if _, err := u.playerTransactionRepo.Create(pctx, playerModule.PlayerTransaction{
		PlayerId: req.PlayerId,
		Amount:   req.Amount,
	}); err != nil {
		logger.Error(err)
		return err
	}

	// Get player saving account
	return nil
}
