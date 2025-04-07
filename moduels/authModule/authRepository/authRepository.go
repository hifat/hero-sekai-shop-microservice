package authRepository

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	IAuthRepository interface {
		CredentialSearch(pctx context.Context, grpcUrl string, req *playerProto.CredentialSearchReq) (*playerProto.PlayerProfile, error)
		InsertOne(pctx context.Context, req *authModule.Credential) (string, error)
		FindByCredentialId(pctx context.Context, credentialId string) (*authModule.Credential, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerProto.FindOnePlayerProfileToRefreshReq) (*playerProto.PlayerProfile, error)
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

func (r *authRepository) InsertOne(pctx context.Context, req *authModule.Credential) (string, error) {
	db := r.dbConn()
	col := db.Collection("auth")

	req.CreatedAt = utils.TimeNow()
	req.UpdatedAt = utils.TimeNow()

	result, err := col.InsertOne(pctx, req)
	if err != nil {
		return "", err
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

	return insertedId, nil
}

func (r *authRepository) FindByCredentialId(pctx context.Context, credentialId string) (*authModule.Credential, error) {
	db := r.dbConn()
	col := db.Collection("auth")

	result := new(authModule.Credential)
	if err := col.FindOne(pctx, bson.M{
		"_id": utils.ConvertToObjectId(credentialId),
	}).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
