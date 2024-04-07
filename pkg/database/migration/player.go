package migration

import (
	"context"
	"log"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/playerModule"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func playerDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("player_db")
}

func PlayerMigrate(pctx context.Context, cfg *config.Config) {
	db := playerDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	/* ---------------------------------- Index --------------------------------- */

	col := db.Collection("player_transactions")

	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "player_id", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	col = db.Collection("players")

	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "email", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	/* ---------------------------------- Seeder --------------------------------- */

	documents := func() []any {
		roles := []*playerModule.Player{
			{
				Email:    "player001@sekai.com",
				Password: "123456",
				Username: "player001",
				PlayerRoles: []playerModule.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.TimeNow(),
				UpdatedAt: utils.TimeNow(),
			},
			{
				Email:    "player002@sekai.com",
				Password: "123456",
				Username: "player002",
				PlayerRoles: []playerModule.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.TimeNow(),
				UpdatedAt: utils.TimeNow(),
			},
			{
				Email:    "player003@sekai.com",
				Password: "123456",
				Username: "player003",
				PlayerRoles: []playerModule.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.TimeNow(),
				UpdatedAt: utils.TimeNow(),
			},
			{
				Email:    "admin001@sekai.com",
				Password: "123456",
				Username: "admin001",
				PlayerRoles: []playerModule.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.TimeNow(),
				UpdatedAt: utils.TimeNow(),
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

	playerTransactions := make([]any, 0)
	for _, p := range results.InsertedIDs {
		playerTransactions = append(playerTransactions, &playerModule.PlayerTransaction{
			PlayerId:  "player:" + p.(primitive.ObjectID).Hex(),
			Amount:    1000,
			CreatedAt: utils.TimeNow(),
			UpdatedAt: utils.TimeNow(),
		})
	}

	col = db.Collection("player_transactions")
	results, err = col.InsertMany(pctx, playerTransactions)
	if err != nil {
		panic(err)
	}

	log.Println("Migrate player_transactions complete: ", results)

	col = db.Collection("player_transactions_queue")
	result, err := col.InsertOne(pctx, bson.M{"offset": -1})
	if err != nil {
		panic(err)
	}

	log.Println("Migrate player_transactions_queue complete: ", result)
}
