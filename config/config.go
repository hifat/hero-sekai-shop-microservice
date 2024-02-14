package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      App
		Db       Db
		Jwt      Jwt
		Kafka    Kafka
		Grpc     Grpc
		Paginate Paginate
	}

	App struct {
		Name  string `mapstructure:"APP_NAME"`
		Url   string `mapstructure:"APP_URL"`
		Stage string `mapstructure:"APP_STAGE"`
	}

	Db struct {
		Url string `mapstructure:"DB_URL"`
	}

	Jwt struct {
		AccessSecretKey   string        `mapstructure:"JWT_ACCESS_SECRET_KEY"`
		AccessDuration    time.Duration `mapstructure:"JWT_ACCESS_DURATION"`
		RefresehSecretKey string        `mapstructure:"JWT_REFRESEH_SECRET_KEY"`
		RefreshDuration   time.Duration `mapstructure:"JWT_REFRESH_DURATION"`
		ApiSecretKey      string        `mapstructure:"JWT_API_SECRET_KEY"`
	}

	Kafka struct {
		Url    string `mapstructure:"KAFKA_URL"`
		ApiKey string `mapstructure:"KAFKA_API_KEY"`
		Secret string `mapstructure:"KAFKA_SECRET"`
	}

	Grpc struct {
		AuthUrl      string `mapstructure:"GRPC_AUTH_URL"`
		PlayerUrl    string `mapstructure:"GRPC_PLAYER_URL"`
		ItemUrl      string `mapstructure:"GRPC_ITEM_URL"`
		InventoryUrl string `mapstructure:"GRPC_INVENTORY_URL"`
		PaymentUrl   string `mapstructure:"GRPC_PAYMENT_URL"`
	}

	Paginate struct {
		ItemNextPageBasedUrl      string `mapstructure:"PAGINATE_ITEM_NEXT_PAGE_BASED_URL"`
		InventoryNextPageBasedUrl string `mapstructure:"PAGINATE_INVENTORY_NEXT_PAGE_BASED_URL"`
	}
)

func (c *Config) Init(path string, filename string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.Unmarshal(&c.App)
	viper.Unmarshal(&c.Db)
	viper.Unmarshal(&c.Jwt)
	viper.Unmarshal(&c.Kafka)
	viper.Unmarshal(&c.Grpc)
	viper.Unmarshal(&c.Paginate)

	return nil
}

func LoadAppConfig(path string, filename string) *Config {
	cfg := Config{}

	err := cfg.Init(path, filename)
	if err != nil {
		panic(err)
	}

	return &cfg
}
