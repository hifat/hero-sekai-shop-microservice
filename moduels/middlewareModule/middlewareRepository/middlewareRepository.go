package middlewareRepository

import (
	"context"
	"errors"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule/authProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IMiddlewareRepository interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
		RoleCount(pctx context.Context, grpcUrl string) (int64, error)
	}

	middlewareRepository struct {
		db *mongo.Client
	}
)

func NewMiddleware(db *mongo.Client) IMiddlewareRepository {
	return &middlewareRepository{db}
}

func (r *middlewareRepository) dbConn() *mongo.Database {
	return r.db.Database("auth_db")
}

func (r *middlewareRepository) AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return err
	}

	result, err := conn.Auth().AccessTokenSearch(pctx, &authProto.AccessTokenSearchReq{
		AccessToken: accessToken,
	})
	if err != nil {
		return err
	}

	if result == nil || !result.IsValid {
		return errors.New("invalid access token")
	}

	return nil
}

func (r *middlewareRepository) RoleCount(pctx context.Context, grpcUrl string) (int64, error) {
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return -1, err
	}

	result, err := conn.Auth().RolesCount(pctx, &authProto.RolesCountReq{})
	if err != nil {
		return -1, err
	}

	return result.Count, nil
}
