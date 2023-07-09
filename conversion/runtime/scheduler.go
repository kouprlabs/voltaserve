package runtime

import (
	"runtime"
	"time"

	"voltaserve/builder"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/infra"

	"go.uber.org/zap"
)

type Scheduler struct {
	pipelineQueue       [][]core.PipelineRunOptions
	builderQueue        [][]core.PipelineRunOptions
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
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &Scheduler{
		pipelineQueue:       make([][]core.PipelineRunOptions, opts.PipelineWorkerCount),
		builderQueue:        make([][]core.PipelineRunOptions, opts.BuilderWorkerCount),
		pipelineWorkerCount: opts.PipelineWorkerCount,
		builderWorkerCount:  opts.BuilderWorkerCount,
		apiClient:           client.NewAPIClient(),
		logger:              logger,
	}
}

func (s *Scheduler) Start() {
	s.logger.Named(infra.StrScheduler).Infow("🚀  launching", "type", "pipeline", "count", s.pipelineWorkerCount)
	for i := 0; i < s.pipelineWorkerCount; i++ {
		go s.pipelineWorker(i)
	}
	s.logger.Named(infra.StrScheduler).Infow("🚀  launching", "type", "builder", "count", s.builderWorkerCount)
	for i := 0; i < s.builderWorkerCount; i++ {
		go s.builderWorker(i)
	}
	go s.pipelineQueueStatus()
	go s.pipelineWorkerStatus()
	go s.builderQueueStatus()
	go s.builderWorkerStatus()
}

func (s *Scheduler) SchedulePipeline(opts *core.PipelineRunOptions) {
	index := 0
	length := len(s.pipelineQueue[0])
	for i := 0; i < s.pipelineWorkerCount; i++ {
		if len(s.pipelineQueue[i]) < length {
			index = i
			length = len(s.pipelineQueue[i])
		}
	}
	s.logger.Named(infra.StrScheduler).Infow("👉  choosing", "pipeline", index)
	s.pipelineQueue[index] = append(s.pipelineQueue[index], *opts)
}

func (s *Scheduler) ScheduleBuilder(opts *core.PipelineRunOptions) {
	index := 0
	length := len(s.builderQueue[0])
	for i := 0; i < s.builderWorkerCount; i++ {
		if len(s.builderQueue[i]) < length {
			index = i
			length = len(s.builderQueue[i])
		}
	}
	s.logger.Named(infra.StrScheduler).Infow("👉  choosing", "builder", index)
	s.builderQueue[index] = append(s.builderQueue[index], *opts)
}

func (s *Scheduler) pipelineWorker(index int) {
	dispatcher := NewDispatcher()
	s.pipelineQueue[index] = make([]core.PipelineRunOptions, 0)
	s.logger.Named(infra.StrPipeline).Infow("⚙️  running", "worker", index)
	for {
		if len(s.pipelineQueue[index]) > 0 {
			s.activePipelineCount++
			opts := s.pipelineQueue[index][0]
			s.logger.Named(infra.StrPipeline).Infow("🔨  working", "worker", index, "bucket", opts.Bucket, "key", opts.Key)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			if err == nil {
				s.logger.Named(infra.StrPipeline).Infow("🎉  succeeded", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			} else {
				s.logger.Named(infra.StrPipeline).Errorw("⛈️  failed", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key, "error", err.Error())
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
	s.builderQueue[index] = make([]core.PipelineRunOptions, 0)
	s.logger.Named(infra.StrBuilder).Infow("⚙️  running", "worker", index)
	for {
		if len(s.builderQueue[index]) > 0 {
			s.activeBuilderCount++
			opts := s.builderQueue[index][0]
			s.logger.Named(infra.StrBuilder).Infow("🔨  working", "worker", index, "bucket", opts.Bucket, "key", opts.Key)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			if err == nil {
				s.logger.Named(infra.StrBuilder).Infow("🎉  succeeded", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key)
			} else {
				s.logger.Named(infra.StrBuilder).Errorw("⛈️  failed", "worker", index, "elapsed", elapsed, "bucket", opts.Bucket, "key", opts.Key, "error", err.Error())
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
				s.logger.Named(infra.StrQueueStatus).Infow("🌈  empty", "type", "pipeline")
			} else {
				s.logger.Named(infra.StrQueueStatus).Infow("⏳  items", "type", "pipeline", "count", sum)
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
				s.logger.Named(infra.StrQueueStatus).Infow("🌈  empty", "type", "builder")
			} else {
				s.logger.Named(infra.StrQueueStatus).Infow("⏳  items", "type", "builder", "count", sum)
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
				s.logger.Named(infra.StrWorkerStatus).Infow("🌤️  all idle", "type", "pipeline")
			} else {
				s.logger.Named(infra.StrWorkerStatus).Infow("🔥  active", "type", "pipeline", "count", s.activePipelineCount)
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
				s.logger.Named(infra.StrWorkerStatus).Infow("🌤️  all idle", "type", "builder")
			} else {
				s.logger.Named(infra.StrWorkerStatus).Infow("🔥  active", "type", "builder", "count", s.activeBuilderCount)
			}
		}
		previous = s.activeBuilderCount
	}
}
