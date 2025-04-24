package playerUsecase

import (
	"context"
	"log/slog"
	"math"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/paymentModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerRepository"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
)

type (
	IPlayerTransactionUsecase interface {
		AddMoney(pctx context.Context, req playerModule.CreatePlayerTransactionReq) (*playerModule.PlayerSavingAccount, error)
		GetSavingAccount(pctx context.Context, playerId string) (*playerModule.PlayerSavingAccount, error)
	}

	playerTransactionUsecase struct {
		playerTransactionRepo playerRepository.IPlayerTransactionRepository
	}
)

func NewPlayerTransaction(playerTransactionRepo playerRepository.IPlayerTransactionRepository) IPlayerTransactionUsecase {
	return &playerTransactionUsecase{playerTransactionRepo}
}

func (u *playerTransactionUsecase) AddMoney(pctx context.Context, req playerModule.CreatePlayerTransactionReq) (*playerModule.PlayerSavingAccount, error) {
	if _, err := u.playerTransactionRepo.Create(pctx, &playerModule.PlayerTransaction{
		PlayerId: req.PlayerId,
		Amount:   req.Amount,
	}); err != nil {
		logger.Error(err)
		return nil, err
	}

	// Get player saving account
	return u.GetSavingAccount(pctx, req.PlayerId)
}

func (u *playerTransactionUsecase) GetSavingAccount(pctx context.Context, playerId string) (*playerModule.PlayerSavingAccount, error) {
	player, err := u.playerTransactionRepo.GetSavingAccount(pctx, playerId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return player, nil
}

func (u *playerTransactionUsecase) DockedPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq) {
	savingAccount, err := u.playerTransactionRepo.GetSavingAccount(pctx, req.PlayerId)
	if err != nil {
		u.playerTransactionRepo.DockedPlayerMoneyRes(pctx, cfg, &paymentModule.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        "",
			Amount:        req.Amount,
			Error:         err.Error(),
		})

		return
	}

	if savingAccount.Balance < math.Abs(req.Amount) {
		slog.Error("Err: DockedPlayerRes failed: not enough money")
		u.playerTransactionRepo.DockedPlayerMoneyRes(pctx, cfg, &paymentModule.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        "",
			Amount:        req.Amount,
			Error:         "not enough money",
		})

		return
	}

	transactionId, err := u.playerTransactionRepo.Create(pctx, &playerModule.PlayerTransaction{
		PlayerId:  req.PlayerId,
		Amount:    req.Amount,
		CreatedAt: utils.TimeNow(),
	})
	if err != nil {
		u.playerTransactionRepo.DockedPlayerMoneyRes(pctx, cfg, &paymentModule.PaymentTransferRes{
			InventoryId:   "",
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        "",
			Amount:        req.Amount,
			Error:         err.Error(),
		})

		return
	}

	u.playerTransactionRepo.DockedPlayerMoneyRes(pctx, cfg, &paymentModule.PaymentTransferRes{
		InventoryId:   "",
		TransactionId: transactionId,
		PlayerId:      req.PlayerId,
		ItemId:        "",
		Amount:        req.Amount,
		Error:         "",
	})
}

func (u *playerTransactionUsecase) RollbackPlayerTransaction(pctx context.Context, req playerModule.RollbackPlayerTransactionReq) {
	if err := u.playerTransactionRepo.DeleteById(pctx, req.TransactionId); err != nil {
		slog.Error(err.Error())
	}
}
