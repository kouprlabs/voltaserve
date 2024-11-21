// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/errorpkg"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/router"
	"github.com/kouprlabs/voltaserve/conversion/runtime"
)

// @title		Voltaserve Conversion
// @version	3.0.0
// @BasePath	/v3
//
// .
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
	installer := runtime.NewInstaller()
	scheduler := runtime.NewScheduler(runtime.SchedulerOptions{
		PipelineWorkerCount: *pipelineWorkers,
		Installer:           installer,
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.MultipartBodyLengthLimitMB)),
	})

	v3 := app.Group("v3")

	healthRouter := router.NewHealthRouter(router.HealthRouterOptions{
		Installer: installer,
	})
	healthRouter.AppendRoutes(v3)

	version := router.NewVersionRouter()
	version.AppendRoutes(app)

	pipelineRouter := router.NewPipelineRouter(router.NewPipelineRouterOptions{
		Scheduler: scheduler,
	})
	pipelineRouter.AppendRoutes(v3)

	scheduler.Start()
	installer.Start()

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
