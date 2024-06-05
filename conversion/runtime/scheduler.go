package runtime

import (
	"runtime"
	"time"
	"voltaserve/pipeline"

	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/infra"

	"go.uber.org/zap"
)

type Scheduler struct {
	pipelineQueue       [][]core.PipelineRunOptions
	pipelineWorkerCount int
	activePipelineCount int
	apiClient           *client.APIClient
	logger              *zap.SugaredLogger
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
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &Scheduler{
		pipelineQueue:       make([][]core.PipelineRunOptions, opts.PipelineWorkerCount),
		pipelineWorkerCount: opts.PipelineWorkerCount,
		apiClient:           client.NewAPIClient(),
		logger:              logger,
	}
}

func (s *Scheduler) Start() {
	s.logger.Named(infra.StrScheduler).Infow("🚀  launching", "type", "pipeline", "count", s.pipelineWorkerCount)
	for i := 0; i < s.pipelineWorkerCount; i++ {
		go s.pipelineWorker(i)
	}
	go s.pipelineQueueStatus()
	go s.pipelineWorkerStatus()
}

func (s *Scheduler) SchedulePipeline(opts *core.PipelineRunOptions) {
	index := s.choosePipeline()
	s.logger.Named(infra.StrScheduler).Infow("👉  choosing", "pipeline", index)
	s.pipelineQueue[index] = append(s.pipelineQueue[index], *opts)
}

/* Choose the pipeline with the least number of items in the queue */
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
