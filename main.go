package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gitnub.com/hifat/hero-sekai-shop-microservice/config"
)

func main() {
	cfg := config.LoadAppConfig(func() (string, string) {
		if len(os.Args) < 2 {
			log.Fatal("Err: env path is required")
		}

		splitPaths := strings.Split(os.Args[1], "/")
		path := strings.Join(splitPaths[:len(splitPaths)-1], "/") + "/"
		filename := splitPaths[len(splitPaths)-1]

		return path, filename
	}())

	fmt.Printf("\n %+v \n", cfg)
}
