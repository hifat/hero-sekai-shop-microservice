package whydoweneedtest

import "gitnub.com/hifat/hero-sekai-shop-microservice/config"

func NewTestConfig() *config.Config {
	return config.LoadAppConfig("../env/test", ".env")
}
