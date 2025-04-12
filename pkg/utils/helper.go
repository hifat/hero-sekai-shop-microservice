package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gitnub.com/hifat/hero-sekai-shop-microservice/pkg/logger"
	"golang.org/x/crypto/bcrypt"
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

func MustStrToTime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", timeStr)
	if err != nil {
		logger.Warn(err.Error())
	}

	return t
}

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 12)
}
