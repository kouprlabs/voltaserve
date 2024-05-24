package infra

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var StrScheduler = fmt.Sprintf("%-13s", "scheduler")
var StrPipeline = fmt.Sprintf("%-13s", "pipeline")
var StrBuilder = fmt.Sprintf("%-13s", "builder")
var StrWorkerStatus = fmt.Sprintf("%-13s", "worker_status")
var StrQueueStatus = fmt.Sprintf("%-13s", "queue_status")
var StrToolRunner = fmt.Sprintf("%-13s", "tool_runner")

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
