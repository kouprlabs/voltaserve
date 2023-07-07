package runtime

import (
	"fmt"
	"runtime"
	"time"

	"voltaserve/builder"
	"voltaserve/client"
	"voltaserve/core"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Scheduler struct {
	pipelineQueue       [][]core.PipelineOptions
	builderQueue        [][]core.PipelineOptions
	pipelineWorkerCount int
	builderWorkerCount  int
	activePipelineCount int
	activeBuilderCount  int
	apiClient           *client.APIClient
	logger              *zap.SugaredLogger
}

type SchedulerOptions struct {
	PipelineWorkerCount int
	BuilderWorkerCount  int
}

var StrScheduler = fmt.Sprintf("%-13s", "scheduler")
var StrPipeline = fmt.Sprintf("%-13s", "pipeline")
var StrBuilder = fmt.Sprintf("%-13s", "builder")
var StrWorkerStatus = fmt.Sprintf("%-13s", "worker_status")
var StrQueueStatus = fmt.Sprintf("%-13s", "queue_status")

func NewDefaultSchedulerOptions() SchedulerOptions {
	opts := SchedulerOptions{}
	if runtime.NumCPU() == 1 {
		opts.PipelineWorkerCount = 1
		opts.BuilderWorkerCount = 1
	} else {
		opts.PipelineWorkerCount = runtime.NumCPU() / 2
		opts.BuilderWorkerCount = runtime.NumCPU() / 2
	}
	return opts
}

func NewScheduler(opts SchedulerOptions) *Scheduler {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableCaller = true
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return &Scheduler{
		pipelineQueue:       make([][]core.PipelineOptions, opts.PipelineWorkerCount),
		builderQueue:        make([][]core.PipelineOptions, opts.BuilderWorkerCount),
		pipelineWorkerCount: opts.PipelineWorkerCount,
		builderWorkerCount:  opts.BuilderWorkerCount,
		apiClient:           client.NewAPIClient(),
		logger:              logger.Sugar(),
	}
}

func (s *Scheduler) Start() {
	s.logger.Named(StrScheduler).Infow("🚀  launching", "type", "pipeline", "count", s.pipelineWorkerCount)
	for i := 0; i < s.pipelineWorkerCount; i++ {
		go s.pipelineWorker(i)
	}
	s.logger.Named(StrScheduler).Infow("🚀  launching", "type", "builder", "count", s.builderWorkerCount)
	for i := 0; i < s.builderWorkerCount; i++ {
		go s.builderWorker(i)
	}
	go s.pipelineQueueStatus()
	go s.pipelineWorkerStatus()
	go s.builderQueueStatus()
	go s.builderWorkerStatus()
}

func (s *Scheduler) SchedulePipeline(opts *core.PipelineOptions) {
	index := 0
	length := len(s.pipelineQueue[0])
	for i := 0; i < s.pipelineWorkerCount; i++ {
		if len(s.pipelineQueue[i]) < length {
			index = i
			length = len(s.pipelineQueue[i])
		}
	}
	s.logger.Named(StrScheduler).Infow("👉  choosing", "pipeline", index)
	s.pipelineQueue[index] = append(s.pipelineQueue[index], *opts)
}

func (s *Scheduler) ScheduleBuilder(opts *core.PipelineOptions) {
	index := 0
	length := len(s.builderQueue[0])
	for i := 0; i < s.builderWorkerCount; i++ {
		if len(s.builderQueue[i]) < length {
			index = i
			length = len(s.builderQueue[i])
		}
	}
	s.logger.Named(StrScheduler).Infow("👉  choosing", "builder", index)
	s.builderQueue[index] = append(s.builderQueue[index], *opts)
}

func (s *Scheduler) pipelineWorker(index int) {
	dispatcher := NewDispatcher()
	s.pipelineQueue[index] = make([]core.PipelineOptions, 0)
	s.logger.Named(StrPipeline).Infow("⚙️  running", "worker", index)
	for {
		if len(s.pipelineQueue[index]) > 0 {
			s.activePipelineCount++
			opts := s.pipelineQueue[index][0]
			s.logger.Named(StrPipeline).Infow("🔨  working", "worker", index, "bucket", opts.Bucket, "key", opts.Key)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			if err == nil {
				s.logger.Named(StrPipeline).Infow("🎉  succeeded", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			} else {
				s.logger.Named(StrPipeline).Errorw("⛈️  failed", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			}
			s.pipelineQueue[index] = s.pipelineQueue[index][1:]
			s.activePipelineCount--
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (s *Scheduler) builderWorker(index int) {
	dispatcher := builder.NewDispatcher()
	s.builderQueue[index] = make([]core.PipelineOptions, 0)
	s.logger.Named(StrBuilder).Infow("⚙️  running", "worker", index)
	for {
		if len(s.builderQueue[index]) > 0 {
			s.activeBuilderCount++
			opts := s.builderQueue[index][0]
			s.logger.Named(StrBuilder).Infow("🔨  working", "worker", index, "bucket", opts.Bucket, "key", opts.Key)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			if err == nil {
				s.logger.Named(StrBuilder).Infow("🎉  succeeded", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			} else {
				s.logger.Named(StrBuilder).Errorw("⛈️  failed", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			}
			s.builderQueue[index] = s.builderQueue[index][1:]
			s.activeBuilderCount--
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
				s.logger.Named(StrQueueStatus).Infow("🌈  empty", "type", "pipeline")
			} else {
				s.logger.Named(StrQueueStatus).Infow("⏳  items", "type", "pipeline", "count", sum)
			}
		}
		previous = sum
	}
}

func (s *Scheduler) builderQueueStatus() {
	previous := -1
	for {
		time.Sleep(5 * time.Second)
		sum := 0
		for i := 0; i < s.builderWorkerCount; i++ {
			sum += len(s.builderQueue[i])
		}
		if sum != previous {
			if sum == 0 {
				s.logger.Named(StrQueueStatus).Infow("🌈  empty", "type", "builder")
			} else {
				s.logger.Named(StrQueueStatus).Infow("⏳  items", "type", "builder", "count", sum)
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
				s.logger.Named(StrWorkerStatus).Infow("🌤️  all idle", "type", "pipeline")
			} else {
				s.logger.Named(StrWorkerStatus).Infow("🔥  active", "type", "pipeline", "count", s.activePipelineCount)
			}
		}
		previous = s.activePipelineCount
	}
}

func (s *Scheduler) builderWorkerStatus() {
	previous := -1
	for {
		time.Sleep(3 * time.Second)
		if previous != s.activeBuilderCount {
			if s.activePipelineCount == 0 {
				s.logger.Named(StrWorkerStatus).Infow("🌤️  all idle", "type", "builder")
			} else {
				s.logger.Named(StrWorkerStatus).Infow("🔥  active", "type", "builder", "count", s.activeBuilderCount)
			}
		}
		previous = s.activeBuilderCount
	}
}
