package runtime

import (
	"fmt"
	"runtime"
	"time"

	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/pipeline"

	"github.com/fatih/color"
)

type Scheduler struct {
	queue       [][]core.PipelineOptions
	workerCount int
	activeCount int
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
	go s.queueStatus()
	go s.workerStatus()
}

func (s *Scheduler) Schedule(opts *core.PipelineOptions) {
	index := 0
	length := len(s.queue[0])
	for i := 0; i < s.workerCount; i++ {
		if len(s.queue[i]) < length {
			index = i
			length = len(s.queue[i])
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
			s.activeCount++
			opts := s.queue[index][0]
			fmt.Printf("[Worker %d] 🚀 Pipeline started ", index)
			helper.PrintlnPipelineOptions(&opts)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			fmt.Printf("[Worker %d] ⌚ Pipeline took ", index)
			color.Set(color.FgMagenta)
			fmt.Printf("%s ", elapsed)
			color.Unset()
			helper.PrintlnPipelineOptions(&opts)
			if err == nil {
				fmt.Printf("[Worker %d] 🎉 Pipeline succeeded! ", index)
				helper.PrintlnPipelineOptions(&opts)
			} else {
				fmt.Printf("[Worker %d] ⛈️  Pipeline failed! ", index)
				helper.PrintlnPipelineOptions(&opts)
			}
			s.queue[index] = s.queue[index][1:]
			s.activeCount--
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (s *Scheduler) queueStatus() {
	previous := -1
	for {
		time.Sleep(5 * time.Second)
		sum := 0
		for i := 0; i < s.workerCount; i++ {
			sum += len(s.queue[i])
		}
		if sum != previous {
			if sum == 0 {
				color.Set(color.FgGreen)
				fmt.Printf("[Status] 🌈 Queue empty\n")
				color.Unset()
			} else {
				color.Set(color.FgBlue)
				fmt.Printf("[Status] ⏳ Items waiting in queue: %d\n", sum)
				color.Unset()
			}
		}
		previous = sum
	}
}

func (s *Scheduler) workerStatus() {
	previous := -1
	for {
		time.Sleep(3 * time.Second)
		if previous != s.activeCount {
			if s.activeCount == 0 {
				color.Set(color.FgGreen)
				fmt.Printf("[Status] 🌤️  Workers idle\n")
				color.Unset()
			} else {
				color.Set(color.FgRed)
				fmt.Printf("[Status] 🔥 Active workers: %d\n", s.activeCount)
				color.Unset()
			}
		}
		previous = s.activeCount
	}
}
