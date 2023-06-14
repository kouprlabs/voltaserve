package runtime

import (
	"fmt"
	"runtime"
	"time"

	"voltaserve/builder"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/pipeline"

	"github.com/fatih/color"
)

type Scheduler struct {
	pipelineQueue       [][]core.PipelineOptions
	builderQueue        [][]core.PipelineOptions
	pipelineWorkerCount int
	builderWorkerCount  int
	activePipelineCount int
	activeBuilderCount  int
	apiClient           *client.APIClient
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
	return &Scheduler{
		pipelineQueue:       make([][]core.PipelineOptions, opts.PipelineWorkerCount),
		builderQueue:        make([][]core.PipelineOptions, opts.BuilderWorkerCount),
		pipelineWorkerCount: opts.PipelineWorkerCount,
		builderWorkerCount:  opts.BuilderWorkerCount,
		apiClient:           client.NewAPIClient(),
	}
}

func (s *Scheduler) Start() {
	fmt.Printf("[Scheduler] 🚀 Launching %d pipeline workers...\n", s.pipelineWorkerCount)
	for i := 0; i < s.pipelineWorkerCount; i++ {
		go s.pipelineWorker(i)
	}
	fmt.Printf("[Scheduler] 🚀 Launching %d builder workers...\n", s.builderWorkerCount)
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
	fmt.Printf("[Scheduler] 👉 Choosing pipline worker %d\n", index)
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
	fmt.Printf("[Scheduler] 👉 Choosing builder worker %d\n", index)
	s.builderQueue[index] = append(s.builderQueue[index], *opts)
}

func (s *Scheduler) pipelineWorker(index int) {
	dispatcher := pipeline.NewDispatcher()
	s.pipelineQueue[index] = make([]core.PipelineOptions, 0)
	fmt.Printf("[Pipeline Worker %d] Running\n", index)
	for {
		if len(s.pipelineQueue[index]) > 0 {
			s.activePipelineCount++
			opts := s.pipelineQueue[index][0]
			fmt.Printf("[Pipeline Worker %d] 🔨 Working... ", index)
			helper.PrintlnPipelineOptions(&opts)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			fmt.Printf("[Pipeline Worker %d] ⌚ Took ", index)
			color.Set(color.FgMagenta)
			fmt.Printf("%s ", elapsed)
			color.Unset()
			helper.PrintlnPipelineOptions(&opts)
			if err == nil {
				fmt.Printf("[Pipeline Worker %d] 🎉 Succeeded! ", index)
				helper.PrintlnPipelineOptions(&opts)
			} else {
				fmt.Printf("[Pipeline Worker %d] ⛈️ Failed! ", index)
				helper.PrintlnPipelineOptions(&opts)
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
	fmt.Printf("[Builder Worker %d] Running\n", index)
	for {
		if len(s.builderQueue[index]) > 0 {
			s.activeBuilderCount++
			opts := s.builderQueue[index][0]
			fmt.Printf("[Builder Worker %d] 🔨 Working... ", index)
			helper.PrintlnPipelineOptions(&opts)
			start := time.Now()
			err := dispatcher.Dispatch(opts)
			elapsed := time.Since(start)
			fmt.Printf("[Builder Worker %d] ⌚ Took ", index)
			color.Set(color.FgMagenta)
			fmt.Printf("%s ", elapsed)
			color.Unset()
			helper.PrintlnPipelineOptions(&opts)
			if err == nil {
				fmt.Printf("[Builder Worker %d] 🎉 Succeeded! ", index)
				helper.PrintlnPipelineOptions(&opts)
			} else {
				fmt.Printf("[Builder Worker %d] ⛈️ Failed! ", index)
				helper.PrintlnPipelineOptions(&opts)
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
				color.Set(color.FgGreen)
				fmt.Printf("[Pipeline Status] 🌈 Queue empty\n")
				color.Unset()
			} else {
				color.Set(color.FgBlue)
				fmt.Printf("[Pipeline Status] ⏳ Items in queue: %d\n", sum)
				color.Unset()
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
				color.Set(color.FgGreen)
				fmt.Printf("[Builder Status] 🌈 Queue empty\n")
				color.Unset()
			} else {
				color.Set(color.FgBlue)
				fmt.Printf("[Builder Status] ⏳ Items in queue: %d\n", sum)
				color.Unset()
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
				color.Set(color.FgGreen)
				fmt.Printf("[Pipeline Status] 🌤️ Workers idle\n")
				color.Unset()
			} else {
				color.Set(color.FgRed)
				fmt.Printf("[Pipeline Status] 🔥 Active workers: %d\n", s.activePipelineCount)
				color.Unset()
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
			if s.activeBuilderCount == 0 {
				color.Set(color.FgGreen)
				fmt.Printf("[Builder Status] 🌤️ Workers idle\n")
				color.Unset()
			} else {
				color.Set(color.FgRed)
				fmt.Printf("[Builder Status] 🔥 Active workers: %d\n", s.activeBuilderCount)
				color.Unset()
			}
		}
		previous = s.activeBuilderCount
	}
}
