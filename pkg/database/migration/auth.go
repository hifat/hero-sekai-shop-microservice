package migration

import (
	"context"
	"log"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/authModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func authDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("auth_db")
}

func AuthMigrate(pctx context.Context, cfg *config.Config) {
	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	/* ---------------------------------- Index --------------------------------- */

	col := db.Collection("auth")

	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "player_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "refresh_token", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	col = db.Collection("roles")

	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "code", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	/* --------------------------------- Seeder --------------------------------- */

	documents := func() []any {
		roles := []*authModule.Role{
			{
				Title: "player",
				Code:  0,
			},
			{
				Title: "admin",
				Code:  1,
			},
		}

		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}

		return docs
	}

	results, err := col.InsertMany(pctx, documents())
	if err != nil {
		panic(err)
	}

	log.Printf("Migrate auth complete: %+v", results)
}
