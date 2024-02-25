package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
)

func Debug(object any) {
	raw, _ := json.MarshalIndent(object, "", "\t")
	fmt.Println(string(raw))
}

func GetEnvPath() (string, string) {
	if len(os.Args) < 2 {
		logger.Error("Err: env path is required")
	}

	splitPaths := strings.Split(os.Args[1], "/")
	path := strings.Join(splitPaths[:len(splitPaths)-1], "/") + "/"
	filename := splitPaths[len(splitPaths)-1]

	return path, filename
}

func TimeNow() *time.Time {
	t := time.Now()
	return &t
}
