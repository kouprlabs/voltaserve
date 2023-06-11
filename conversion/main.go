package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/runtime"

	log "github.com/sirupsen/logrus"

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

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	cfg := config.GetConfig()
	scheduler := runtime.NewScheduler()
	app := fiber.New()

	app.Get("v1/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("v1/run_pipeline", func(c *fiber.Ctx) error {
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
		opts := new(core.PipelineOptions)
		if err := c.BodyParser(opts); err != nil {
			return err
		}
		scheduler.Schedule(opts)
		return c.SendStatus(200)
	})

	scheduler.Start()

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
