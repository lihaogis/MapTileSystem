package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(level string) {
	var config zap.Config
	if level == "debug" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}
