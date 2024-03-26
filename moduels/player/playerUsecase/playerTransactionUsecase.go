package playerUsecase

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/player/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
)

type (
	IPlayerTransactionUsecase interface {
		AddMoney(pctx context.Context, req player.CreatePlayerTransactionReq) (*player.PlayerTransaction, error)
	}

	playerTransactionUsecase struct {
		playerTransactionRepo playerRepository.IPlayerTransactionRepository
	}
)

func NewPlayerTransaction(playerTransactionRepo playerRepository.IPlayerTransactionRepository) IPlayerTransactionUsecase {
	return &playerTransactionUsecase{playerTransactionRepo}
}

func (u *playerTransactionUsecase) AddMoney(pctx context.Context, req player.CreatePlayerTransactionReq) (*player.PlayerTransaction, error) {
	if _, err := u.playerTransactionRepo.Create(pctx, player.PlayerTransaction{
		PlayerId: req.PlayerId,
		Amount:   req.Amount,
	}); err != nil {
		logger.Error(err)
		return nil, err
	}

	// Get player saving account
	return nil, nil
}
