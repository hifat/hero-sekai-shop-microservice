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

func (g *playerGrpc) CredentialSearch(ctx context.Context, req *playerProto.CredentialSearchReq) (*playerProto.PlayerProfile, error) {
	return g.playerUsecase.FindByCredential(ctx, req)
}

func (g *playerGrpc) FindOnePlayerProfileToRefresh(ctx context.Context, req *playerProto.FindOnePlayerProfileToRefreshReq) (*playerProto.PlayerProfile, error) {
	return g.playerUsecase.FindOnePlayerToRefreshToken(ctx, req.PlayerId)
}

func (g *playerGrpc) GetPlayerSavingAccount(context.Context, *playerProto.GetPlayerSavingAccountReq) (*playerProto.GetPlayerSavingAccountRes, error) {
	return nil, nil
}
