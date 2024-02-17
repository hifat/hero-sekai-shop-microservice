package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	// object key support cloud run logging
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "receiveTimestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.LevelKey = "severity"

	var err error
	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Skip(skip int) *zap.Logger {
	return log.WithOptions(zap.AddCallerSkip(skip))
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	log.Warn(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	log.Fatal(message, fields...)
}

func Error(message interface{}, fields ...zap.Field) {
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), fields...)
	case string:
		log.Error(v, fields...)
	}
}
