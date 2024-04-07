package playerHandler

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerUsecase"
)

type (
	playerGrpc struct {
		playerProto.UnimplementedPlayerGrpcServiceServer
		playerUsecase playerUsecase.IPlayerUsecase
	}
)

func NewPlayerGrpc(playerUsecase playerUsecase.IPlayerUsecase) *playerGrpc {
	return &playerGrpc{
		playerUsecase: playerUsecase,
	}
}

func (g *playerGrpc) CredentialSearch(context.Context, *playerProto.CredentialSearchReq) (*playerProto.PlayerProfile, error) {
	return nil, nil
}

func (g *playerGrpc) FindOnePlayerProfileToRefresh(context.Context, *playerProto.FindOnePlayerProfileToRefreshReq) (*playerProto.PlayerProfile, error) {
	return nil, nil
}

func (g *playerGrpc) GetPlayerSavingAccount(context.Context, *playerProto.GetPlayerSavingAccountReq) (*playerProto.GetPlayerSavingAccountRes, error) {
	return nil, nil
}
