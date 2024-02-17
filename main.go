package main

import (
	"context"
	"log"
	"os"
	"strings"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadAppConfig(func() (string, string) {
		if len(os.Args) < 2 {
			log.Fatal()
			logger.Error("Err: env path is required")
		}

		splitPaths := strings.Split(os.Args[1], "/")
		path := strings.Join(splitPaths[:len(splitPaths)-1], "/") + "/"
		filename := splitPaths[len(splitPaths)-1]

		return path, filename
	}())

	db := database.DbConnect(ctx, cfg)
	defer db.Disconnect(ctx)
}
