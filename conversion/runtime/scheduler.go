package runtime

import (
	"fmt"
	"runtime"
	"time"

	"voltaserve/core"
	"voltaserve/infra"
	"voltaserve/pipeline"

	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	queue       [][]core.PipelineOptions
	workerCount int
	apiClient   *infra.APIClient
}

func NewScheduler() *Scheduler {
	workerCount := runtime.NumCPU()
	return &Scheduler{
		queue:       make([][]core.PipelineOptions, workerCount),
		workerCount: workerCount,
		apiClient:   infra.NewAPIClient(),
	}
}

func (s *Scheduler) Start() {
	fmt.Printf("[Scheduler] Starting %d workers\n", s.workerCount)
	for i := 0; i < s.workerCount; i++ {
		go s.worker(i)
	}
	go s.status()
}

func (s *Scheduler) Schedule(opts *core.PipelineOptions) {
	index := 0
	queueLength := len(s.queue[0])
	for i := 0; i < s.workerCount; i++ {
		if len(s.queue[i]) < queueLength {
			index = i
			queueLength = len(s.queue[i])
		}
	}
	fmt.Printf("[Scheduler] 👉 Choosing worker %d\n", index)
	s.queue[index] = append(s.queue[index], *opts)
}

func (s *Scheduler) worker(index int) {
	dispatcher := pipeline.NewDispatcher()
	s.queue[index] = make([]core.PipelineOptions, 0)
	fmt.Printf("[Worker %d] Started\n", index)
	for {
		if len(s.queue[index]) > 0 {
			opts := s.queue[index][0]
			s.queue[index] = s.queue[index][1:]
			fmt.Printf("[Worker %d] 🚀 Pipeline started (FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			start := time.Now()
			pr, err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			fmt.Printf("[Worker %d] ⌚ Pipeline took %s (FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)\n", index, elapsed, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			if err == nil {
				pr.Options = opts
				fmt.Printf("[Worker %d] 📝 Updating snapshot (FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
				if err := s.apiClient.UpdateSnapshot(&pr); err != nil {
					fmt.Printf("[Worker %d] 🔥 Failed to update snapshot!\n", index)
					log.Error(err)
				}
				fmt.Printf("[Worker %d] 🎉 Succeeded! (FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			} else {
				fmt.Printf("[Worker %d] 🔥 Pipeline failed! (FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s)\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			}
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (s *Scheduler) status() {
	previousSum := -1
	for {
		time.Sleep(5 * time.Second)
		sum := 0
		for i := 0; i < s.workerCount; i++ {
			sum += len(s.queue[i])
		}
		if sum != previousSum {
			if sum == 0 {
				fmt.Printf("[Status] 🌈 Queue empty\n")
			} else {
				fmt.Printf("[Status] ⏳ Items waiting in queue: %d\n", sum)
			}
		}
		previousSum = sum
	}
}
