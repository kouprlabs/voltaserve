package infra

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func GetLogger() (*zap.SugaredLogger, error) {
	if logger == nil {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableCaller = true
		if l, err := config.Build(); err != nil {
			return nil, err
		} else {
			logger = l.Sugar()
		}
	}
	return logger, nil
}
