package main

import (
	"context"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/database/migration"
	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/utils"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadAppConfig(utils.GetEnvPath())

	switch cfg.App.Name {
	case "player":
		migration.PlayerMigrate(ctx, cfg)
	case "auth":
		migration.AuthMigrate(ctx, cfg)
	case "item":
		migration.ItemMigrate(ctx, cfg)
	case "inventory":
		migration.InventoryMigrate(ctx, cfg)
	case "payment":
		migration.PaymentMigrate(ctx, cfg)
	}
}
