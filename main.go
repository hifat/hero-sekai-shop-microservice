package main

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
	"gitnub.com/hifat/hero-sekai-shop-microservice/server"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadAppConfig(utils.GetEnvPath())

	db := database.DbConnect(ctx, cfg)
	defer db.Disconnect(ctx)

	server.Start(ctx, cfg, db)
}
