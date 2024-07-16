// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	StrScheduler            = fmt.Sprintf("%-13s", "scheduler")
	StrPipeline             = fmt.Sprintf("%-13s", "pipeline")
	StrWorkerStatus         = fmt.Sprintf("%-13s", "worker_status")
	StrQueueStatus          = fmt.Sprintf("%-13s", "queue_status")
	StrDependencyDownloader = fmt.Sprintf("%-13s", "installer")
)

var logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.DisableCaller = true
		if l, err := config.Build(); err != nil {
			panic(err)
		} else {
			logger = l.Sugar()
		}
	}
	return logger
}
