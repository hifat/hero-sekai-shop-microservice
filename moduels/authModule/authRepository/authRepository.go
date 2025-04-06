package authRepository

import (
	"context"
	"fmt"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/grpccon"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IAuthRepository interface {
		CredentialSearch(pctx context.Context, grpcUrl string, req *playerProto.CredentialSearchReq) (*playerProto.PlayerProfile, error)
		InsertOne(pctx context.Context, req *authModule.Credential) (string, error)
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuth(db *mongo.Client) IAuthRepository {
	return &authRepository{db}
}

func (r *authRepository) dbConn() *mongo.Database {
	return r.db.Database("auth_db")
}

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

func (r *authRepository) InsertOne(pctx context.Context, req *authModule.Credential) (string, error) {
	db := r.dbConn()
	col := db.Collection("auth")

	result, err := col.InsertOne(pctx, req)
	if err != nil {
		return "", err
	}

	fmt.Println(result.InsertedID.(primitive.ObjectID))
	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

	return insertedId, nil
}
