// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package runtime

import (
	"runtime"
	"time"

	"github.com/kouprlabs/voltaserve/conversion/client"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/pipeline"
)

type Scheduler struct {
	pipelineQueue       [][]client.PipelineRunOptions
	pipelineWorkerCount int
	activePipelineCount int
	apiClient           *client.APIClient
}

type SchedulerOptions struct {
	PipelineWorkerCount int
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
		pipelineQueue:       make([][]client.PipelineRunOptions, opts.PipelineWorkerCount),
		pipelineWorkerCount: opts.PipelineWorkerCount,
		apiClient:           client.NewAPIClient(),
	}
}

func (s *Scheduler) Start() {
	infra.GetLogger().Named(infra.StrScheduler).Infow("🚀  launching", "type", "pipeline", "count", s.pipelineWorkerCount)
	for i := 0; i < s.pipelineWorkerCount; i++ {
		go s.pipelineWorker(i)
	}
	go s.pipelineQueueStatus()
	go s.pipelineWorkerStatus()
}

func (s *Scheduler) SchedulePipeline(opts *client.PipelineRunOptions) {
	index := s.choosePipeline()
	infra.GetLogger().Named(infra.StrScheduler).Infow("👉  choosing", "pipeline", index)
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
	s.pipelineQueue[index] = make([]client.PipelineRunOptions, 0)
	infra.GetLogger().Named(infra.StrPipeline).Infow("⚙️  running", "worker", index)
	for {
		if len(s.pipelineQueue[index]) > 0 {
			s.activePipelineCount++
			opts := s.pipelineQueue[index][0]
			infra.GetLogger().Named(infra.StrPipeline).Infow("🔨  working", "worker", index, "bucket", opts.Bucket, "key", opts.Key)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			if err == nil {
				infra.GetLogger().Named(infra.StrPipeline).Infow("🎉  succeeded", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			} else {
				infra.GetLogger().Named(infra.StrPipeline).Errorw("⛈️  failed", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key, "error", err.Error())
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
				infra.GetLogger().Named(infra.StrQueueStatus).Infow("🌈  empty", "type", "pipeline")
			} else {
				infra.GetLogger().Named(infra.StrQueueStatus).Infow("⏳  items", "type", "pipeline", "count", sum)
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
				infra.GetLogger().Named(infra.StrWorkerStatus).Infow("🌤️  all idle", "type", "pipeline")
			} else {
				infra.GetLogger().Named(infra.StrWorkerStatus).Infow("🔥  active", "type", "pipeline", "count", s.activePipelineCount)
			}
		}
		previous = s.activePipelineCount
	}
}
