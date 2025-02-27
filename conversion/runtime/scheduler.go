// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package runtime

import (
	"runtime"
	"time"

	"github.com/kouprlabs/voltaserve/shared/dto"

	"github.com/kouprlabs/voltaserve/conversion/logger"
	"github.com/kouprlabs/voltaserve/conversion/pipeline"
)

type Scheduler struct {
	pipelineQueue       [][]dto.PipelineRunOptions
	pipelineWorkerCount int
	activePipelineCount int
	installer           *Installer
}

type SchedulerOptions struct {
	PipelineWorkerCount int
	Installer           *Installer
}

func NewDefaultSchedulerOptions() SchedulerOptions {
	opts := SchedulerOptions{}
	if runtime.NumCPU() == 1 {
		opts.PipelineWorkerCount = 1
	} else {
		opts.PipelineWorkerCount = runtime.NumCPU()
	}
	return opts
}

func NewScheduler(opts SchedulerOptions) *Scheduler {
	return &Scheduler{
		pipelineQueue:       make([][]dto.PipelineRunOptions, opts.PipelineWorkerCount),
		pipelineWorkerCount: opts.PipelineWorkerCount,
		installer:           opts.Installer,
	}
}

func (s *Scheduler) Start() {
	logger.GetLogger().Named(logger.StrScheduler).Infow("🚀  launching", "type", "pipeline", "count", s.pipelineWorkerCount)
	for i := 0; i < s.pipelineWorkerCount; i++ {
		go s.pipelineWorker(i)
	}
	go s.pipelineQueueStatus()
	go s.pipelineWorkerStatus()
}

func (s *Scheduler) SchedulePipeline(opts *dto.PipelineRunOptions) {
	index := s.choosePipeline()
	logger.GetLogger().Named(logger.StrScheduler).Infow("👉  choosing", "pipeline", index)
	s.pipelineQueue[index] = append(s.pipelineQueue[index], *opts)
}

// Choose the pipeline with the least number of items in the queue.
func (s *Scheduler) choosePipeline() int {
	index := 0
	length := len(s.pipelineQueue[0])
	for i := 0; i < s.pipelineWorkerCount; i++ {
		if len(s.pipelineQueue[i]) < length {
			index = i
			length = len(s.pipelineQueue[i])
		}
	}
	return index
}

func (s *Scheduler) pipelineWorker(index int) {
	dispatcher := pipeline.NewDispatcher()
	s.pipelineQueue[index] = make([]dto.PipelineRunOptions, 0)
	logger.GetLogger().Named(logger.StrPipeline).Infow("⚙️  running", "worker", index)
	for {
		if len(s.pipelineQueue[index]) > 0 && !s.installer.IsRunning() {
			s.activePipelineCount++
			opts := s.pipelineQueue[index][0]
			logger.GetLogger().Named(logger.StrPipeline).Infow("🔨  working", "worker", index, "bucket", opts.Bucket, "key", opts.Key)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			if err == nil {
				logger.GetLogger().Named(logger.StrPipeline).
					Infow("🎉  succeeded", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			} else {
				logger.GetLogger().Named(logger.StrPipeline).
					Errorw("⛈️  failed", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key, "error", err.Error())
			}
			s.pipelineQueue[index] = s.pipelineQueue[index][1:]
			s.activePipelineCount--
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (s *Scheduler) pipelineQueueStatus() {
	previous := -1
	for {
		time.Sleep(5 * time.Second)
		sum := 0
		for i := 0; i < s.pipelineWorkerCount; i++ {
			sum += len(s.pipelineQueue[i])
		}
		if sum != previous {
			if sum == 0 {
				logger.GetLogger().Named(logger.StrQueueStatus).Infow("🌈  empty", "type", "pipeline")
			} else {
				logger.GetLogger().Named(logger.StrQueueStatus).Infow("⏳  items", "type", "pipeline", "count", sum)
			}
		}
		previous = sum
	}
}

func (s *Scheduler) pipelineWorkerStatus() {
	previous := -1
	for {
		time.Sleep(3 * time.Second)
		if previous != s.activePipelineCount {
			if s.activePipelineCount == 0 {
				logger.GetLogger().Named(logger.StrWorkerStatus).Infow("🌤️  all idle", "type", "pipeline")
			} else {
				logger.GetLogger().Named(logger.StrWorkerStatus).
					Infow("🔥  active", "type", "pipeline", "count", s.activePipelineCount)
			}
		}
		previous = s.activePipelineCount
	}
}
