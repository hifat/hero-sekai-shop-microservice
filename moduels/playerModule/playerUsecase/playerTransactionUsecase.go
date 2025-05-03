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
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		AddMoney(pctx context.Context, req playerModule.CreatePlayerTransactionReq) (*playerModule.PlayerSavingAccount, error)
		GetSavingAccount(pctx context.Context, playerId string) (*playerModule.PlayerSavingAccount, error)

		// Queue
		DockedPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq)
		RollbackPlayerTransaction(pctx context.Context, req *playerModule.RollbackPlayerTransactionReq)
		AddPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq)
	}

	playerTransactionUsecase struct {
		playerTransactionRepo playerRepository.IPlayerTransactionRepository
	}
)

func NewPlayerTransaction(playerTransactionRepo playerRepository.IPlayerTransactionRepository) IPlayerTransactionUsecase {
	return &playerTransactionUsecase{playerTransactionRepo}
}

func (u *playerTransactionUsecase) GetOffset(pctx context.Context) (int64, error) {
	offset, err := u.playerTransactionRepo.GetOffset(pctx)
	if err != nil {
		return -1, err
	}

	return offset, nil
}

func (u *playerTransactionUsecase) UpsertOffset(pctx context.Context, offset int64) error {
	err := u.playerTransactionRepo.UpsertOffset(pctx, offset)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
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
	res := &paymentModule.PaymentTransferRes{
		InventoryId:   "",
		TransactionId: "",
		PlayerId:      req.PlayerId,
		ItemId:        "",
		Amount:        req.Amount,
		Error:         "",
	}

	isValidateFailed := false
	if err != nil {
		isValidateFailed = true
		res.Error = err.Error()
	}

	if savingAccount.Balance < math.Abs(req.Amount) {
		isValidateFailed = true
		slog.Error("Err: DockedPlayerRes failed: not enough money")
		res.Error = "not enough money"
	}

	if !isValidateFailed {
		transactionId, err := u.playerTransactionRepo.Create(pctx, &playerModule.PlayerTransaction{
			PlayerId:  req.PlayerId,
			Amount:    req.Amount,
			CreatedAt: utils.TimeNow(),
		})
		if err != nil {
			res.Error = err.Error()
		}

		res.TransactionId = transactionId
	}

	if err := u.playerTransactionRepo.DockedPlayerMoneyRes(pctx, cfg, res); err != nil {
		logger.Error(err)
	}
}

func (u *playerTransactionUsecase) RollbackPlayerTransaction(pctx context.Context, req *playerModule.RollbackPlayerTransactionReq) {
	if err := u.playerTransactionRepo.DeleteById(pctx, req.TransactionId); err != nil {
		slog.Error(err.Error())
	}
}

func (u *playerTransactionUsecase) AddPlayerMoneyRes(pctx context.Context, cfg *config.Config, req *playerModule.CreatePlayerTransactionReq) {
	savingAccount, err := u.playerTransactionRepo.GetSavingAccount(pctx, req.PlayerId)
	res := &paymentModule.PaymentTransferRes{
		InventoryId:   "",
		TransactionId: "",
		PlayerId:      req.PlayerId,
		ItemId:        "",
		Amount:        req.Amount,
		Error:         "",
	}

	isValidateFailed := false
	if err != nil {
		isValidateFailed = true
		res.Error = err.Error()
	}

	if savingAccount.Balance < math.Abs(req.Amount) {
		isValidateFailed = true
		slog.Error("Err: AddPlayerRes failed: not enough money")
		res.Error = "not enough money"
	}

	if !isValidateFailed {
		transactionId, err := u.playerTransactionRepo.Create(pctx, &playerModule.PlayerTransaction{
			PlayerId:  req.PlayerId,
			Amount:    req.Amount,
			CreatedAt: utils.TimeNow(),
		})
		if err != nil {
			res.Error = err.Error()
		}

		res.TransactionId = transactionId
	}

	if err := u.playerTransactionRepo.AddPlayerMoneyRes(pctx, cfg, res); err != nil {
		logger.Error(err)
	}
}
