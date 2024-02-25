package migration

import (
	"context"
	"log"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func inventoryDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("inventory_db")
}

func InventoryMigrate(pctx context.Context, cfg *config.Config) {
	db := inventoryDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	/* ---------------------------------- Index --------------------------------- */

	col := db.Collection("player_inventories")

	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "player_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "item_id", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	/* ---------------------------------- Seeder --------------------------------- */

	col = db.Collection("player_inventories_queue")

	results, err := col.InsertOne(pctx, bson.M{"offset": -1})
	if err != nil {
		panic(err)
	}

	log.Printf("Migrate inventories complete: %+v", results)
}
