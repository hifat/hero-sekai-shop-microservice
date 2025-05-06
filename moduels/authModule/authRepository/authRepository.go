package authRepository

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule/playerProto"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/jwtauth"
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
		UpdateRefreshToken(pctx context.Context, credentialId string, req *authModule.UpdateRefreshTokenReq) error
		DeleteByCredentialId(pctx context.Context, credentialId string) (int64, error)
		FindByAccessToken(pctx context.Context, accessToken string) (*authModule.Credential, error)
		RoleCount(pctx context.Context) (int64, error)
		NewAccessToken(cfg config.Jwt, profile *playerProto.PlayerProfile) jwtauth.AuthFactory
		NewRefreshToken(cfg config.Jwt, profile *playerProto.PlayerProfile) jwtauth.AuthFactory
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

func (r *authRepository) UpdateRefreshToken(pctx context.Context, credentialId string, req *authModule.UpdateRefreshTokenReq) error {
	db := r.dbConn()
	col := db.Collection("auth")

	_, err := col.UpdateOne(
		pctx,
		bson.M{
			"_id": utils.ConvertToObjectId(credentialId),
		},
		bson.M{
			"$set": bson.M{
				"player_id":     req.PlayerId,
				"access_token":  req.AccessToken,
				"refresh_token": req.RefreshToken,
				"updated_at":    req.UpdatedAt,
			},
		})
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) DeleteByCredentialId(pctx context.Context, credentialId string) (int64, error) {
	db := r.dbConn()
	col := db.Collection("auth")

	result, err := col.DeleteOne(pctx, bson.M{
		"_id": utils.ConvertToObjectId(credentialId),
	})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *authRepository) FindByAccessToken(pctx context.Context, accessToken string) (*authModule.Credential, error) {
	db := r.dbConn()
	col := db.Collection("auth")

	credential := new(authModule.Credential)
	if err := col.FindOne(pctx, bson.M{
		"access_token": accessToken,
	}).Decode(credential); err != nil {
		return nil, err
	}

	return credential, nil
}

func (r *authRepository) RoleCount(pctx context.Context) (int64, error) {
	db := r.dbConn()
	col := db.Collection("roles")

	count, err := col.CountDocuments(pctx, bson.M{})
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *authRepository) NewAccessToken(cfg config.Jwt, profile *playerProto.PlayerProfile) jwtauth.AuthFactory {
	return jwtauth.NewAccessToken(cfg.RefreshSecretKey, cfg.AccessDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: profile.RoleCode,
	})

}

func (r *authRepository) NewRefreshToken(cfg config.Jwt, profile *playerProto.PlayerProfile) jwtauth.AuthFactory {
	return jwtauth.NewRefreshToken(cfg.RefreshSecretKey, cfg.RefreshDuration, &jwtauth.Claims{
		PlayerId: profile.Id,
		RoleCode: profile.RoleCode,
	})
}
