package authRepository

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
)

func (r *authRepository) CredentialSearch(pctx context.Context, grpcUrl string, req *playerProto.CredentialSearchReq) (*playerProto.PlayerProfile, error) {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}

	result, err := conn.Player().CredentialSearch(pctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *authRepository) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerProto.FindOnePlayerProfileToRefreshReq) (*playerProto.PlayerProfile, error) {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, err
	}

	result, err := conn.Player().FindOnePlayerProfileToRefresh(pctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
