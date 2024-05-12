package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

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

	cfg := config.GetConfig()

	schedulerOpts := runtime.NewDefaultSchedulerOptions()
	pipelineWorkers := flag.Int("pipeline-workers", schedulerOpts.PipelineWorkerCount, "Number of pipeline workers")
	flag.Parse()
	scheduler := runtime.NewScheduler(runtime.SchedulerOptions{
		PipelineWorkerCount: *pipelineWorkers,
	})

	app := fiber.New()
	v1 := app.Group("v1")

	v1.Get("health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	v1.Post("pipelines/run", func(c *fiber.Ctx) error {
		apiKey := c.Query("api_key")
		if apiKey == "" {
			if err := c.SendStatus(http.StatusBadRequest); err != nil {
				return err
			}
			return errors.New("missing query param api_key")
		}
		if apiKey != cfg.Security.APIKey {
			if err := c.SendStatus(http.StatusUnauthorized); err != nil {
				return err
			}
			return errors.New("invalid api_key")
		}
		opts := new(core.PipelineRunOptions)
		if err := c.BodyParser(opts); err != nil {
			return err
		}
		scheduler.SchedulePipeline(opts)
		return c.SendStatus(200)
	})

	scheduler.Start()

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
