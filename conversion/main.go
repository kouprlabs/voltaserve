package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/pipeline"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var queue [][]core.PipelineOptions
var workerCount = 1

func pipelineWorker(index int) {
	dispatcher := pipeline.NewDispatcher()
	queue[index] = make([]core.PipelineOptions, 0)
	fmt.Printf("[Worker %d] Running...\n", index)
	for {
		if len(queue[index]) > 0 {
			opts := queue[index][0]
			queue[index] = queue[index][1:]
			fmt.Printf("[Worker %d] [🚀 Pipeline started] FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			pipelineResponse, err := dispatcher.Dispatch(opts)
			if err == nil {
				pipelineResponse.Options = opts
				fmt.Printf("[Worker %d] [👍 Pipeline completed] FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
				fmt.Printf("[Worker %d] [📃 Pipeline Result] Thumbnail=%t Preview=%t Text=%t OCR=%t\n", index, pipelineResponse.Thumbnail != nil, pipelineResponse.Preview != nil, pipelineResponse.Text != nil, pipelineResponse.OCR != nil)
				body, err := json.Marshal(pipelineResponse)
				if err != nil {
					log.Error(err)
					continue
				}
				cfg := config.GetConfig()
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/files/conversion_webhook/update_snapshot?api_key=%s", cfg.APIURL, cfg.Security.APIKey), bytes.NewBuffer(body))
				if err != nil {
					log.Error(err)
					continue
				}
				req.Header.Set("Content-Type", "application/json; charset=UTF-8")
				client := &http.Client{}
				res, err := client.Do(req)
				if err != nil {
					fmt.Printf("[Worker %d] [🔥 Request failed!]\n", index)
					log.Error(err)
					continue
				}
				res.Body.Close()
				fmt.Printf("[Worker %d] [🎉 Request succeeded!]\n", index)
			} else {
				fmt.Printf("[Worker %d] [🔥 Pipeline failed] FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s\n", index, opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			}
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func statusWorker() {
	for {
		time.Sleep(5 * time.Second)
		sum := 0
		for i := 0; i < workerCount; i++ {
			sum += len(queue[i])
		}
		if sum == 0 {
			fmt.Printf("[Status] 🌈 Queue empty!\n")
		} else {
			fmt.Printf("[Status] ⏳ Total items in queue: %d\n", sum)
		}
	}
}

func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		err := godotenv.Load(".env.local")
		if err != nil {
			panic(err)
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	cfg := config.GetConfig()

	app := fiber.New()

	app.Get("v1/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("v1/pipelines", func(c *fiber.Ctx) error {
		apiKey := c.Query("api_key")
		if apiKey == "" {
			c.SendStatus(http.StatusBadRequest)
			return errors.New("missing query param api_key")
		}
		if apiKey != cfg.Security.APIKey {
			c.SendStatus(http.StatusUnauthorized)
			return errors.New("invalid api_key")
		}
		opts := new(core.PipelineOptions)
		if err := c.BodyParser(opts); err != nil {
			return err
		}
		workerIndex := 0
		queueLength := len(queue[0])
		for i := 0; i < workerCount; i++ {
			if len(queue[i]) < queueLength {
				workerIndex = i
				queueLength = len(queue[i])
			}
		}
		fmt.Printf("Choosing worker 👉 %d\n", workerIndex)
		queue[workerIndex] = append(queue[workerIndex], *opts)
		return c.SendStatus(200)
	})

	fmt.Printf("Number of CPU cores: %d\n", runtime.NumCPU())

	workerCount = runtime.NumCPU()

	fmt.Printf("Setting the number of workers to: %d\n", workerCount)

	queue = make([][]core.PipelineOptions, workerCount)

	for i := 0; i < workerCount; i++ {
		go pipelineWorker(i)
	}

	go statusWorker()

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
