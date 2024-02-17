package database

import (
	"context"
	"fmt"
	"time"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DbConnect(pctx context.Context, cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.Url))
	if err != nil {
		logger.Error(fmt.Sprintf("Connect to database error: %s", err.Error()))
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error(fmt.Sprintf("Error: Pinging to databse error: %s", err.Error()))
	}

	return client
}
