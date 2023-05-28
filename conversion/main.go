package main

import (
	"fmt"
	"net/url"
	"os"

	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/pipeline"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var queue = make([]core.PipelineOptions, 0)

func worker() {
	dispatcher := pipeline.NewDispatcher()
	for {
		if len(queue) > 0 {
			opts := queue[0]
			queue = queue[1:]
			fmt.Printf("[Started 🚀] FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s\n", opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			res, err := dispatcher.Dispatch(opts)
			if err == nil {
				fmt.Printf("[Completed 🎉] FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s\n", opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
				fmt.Printf("[Result ☕️] Thumbnail=%t Preview=%t Text=%t OCR=%t\n", res.Thumbnail != nil, res.Preview != nil, res.Text != nil, res.OCR != nil)
			} else {
				fmt.Printf("[Failed ❌] FileID=%s SnapshotID=%s S3Bucket=%s S3Key=%s\n", opts.FileID, opts.SnapshotID, opts.Bucket, opts.Key)
			}
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

	settings := config.GetConfig()

	app := fiber.New()

	app.Post("/pipelines", func(c *fiber.Ctx) error {
		opts := new(core.PipelineOptions)
		if err := c.BodyParser(opts); err != nil {
			return err
		}
		queue = append(queue, *opts)
		return c.SendStatus(200)
	})

	go worker()

	url, err := url.Parse(settings.ConversionURL)
	if err != nil {
		panic(err)
	}
	if err := app.Listen(":" + url.Port()); err != nil {
		panic(err)
	}
}
