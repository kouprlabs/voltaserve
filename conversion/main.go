package main

import (
	"flag"
	"fmt"
	"os"

	"voltaserve/config"
	"voltaserve/router"
	"voltaserve/runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @title		Voltaserve Conversion
// @version	2.0.0
// @BasePath	/v1
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

	healthRouter := router.NewHealthRouter()
	healthRouter.AppendRoutes(v1)

	pipelineRouter := router.NewPipelineRouter(router.NewPipelineRouterOptions{
		Scheduler: scheduler,
	})
	pipelineRouter.AppendRoutes(v1)

	toolsRouter := router.NewToolRouter()
	toolsRouter.AppendRoutes(v1)

	scheduler.Start()

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
