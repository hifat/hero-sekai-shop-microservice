package migration

import (
	"context"
	"log"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/moduels/item"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func itemDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("item_db")
}

func ItemMigrate(pctx context.Context, cfg *config.Config) {
	db := itemDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	/* ---------------------------------- Index --------------------------------- */

	col := db.Collection("items")

	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{primitive.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{primitive.E{Key: "title", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	/* ---------------------------------- Seeder --------------------------------- */

	documents := func() []any {
		items := []*item.Item{
			{
				Title:       "Diamon Sword",
				Price:       1000,
				Damage:      100,
				ImageUrl:    "https://shorturl.at/pKPZ9",
				UsageStatus: true,
				CreatedAt:   utils.TimeNow(),
				UpdatedAt:   utils.TimeNow(),
			},
			{
				Title:       "Iron Sword",
				Price:       500,
				Damage:      50,
				ImageUrl:    "https://shorturl.at/pKPZ9",
				UsageStatus: true,
				CreatedAt:   utils.TimeNow(),
				UpdatedAt:   utils.TimeNow(),
			},
			{
				Title:       "Wooden Sword",
				Price:       100,
				Damage:      20,
				ImageUrl:    "https://shorturl.at/pKPZ9",
				UsageStatus: true,
				CreatedAt:   utils.TimeNow(),
				UpdatedAt:   utils.TimeNow(),
			},
		}

		docs := make([]any, 0)
		for _, r := range items {
			docs = append(docs, r)
		}

		return docs
	}

	results, err := col.InsertMany(pctx, documents())
	if err != nil {
		panic(err)
	}

	log.Printf("Migrate items complete: %+v", results)
}
