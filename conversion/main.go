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
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/router"
	"github.com/kouprlabs/voltaserve/conversion/runtime"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var e *errorpkg.ErrorResponse
	if errors.As(err, &e) {
		var v *errorpkg.ErrorResponse
		errors.As(err, &v)
		return c.Status(v.Status).JSON(v)
	} else {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).JSON(errorpkg.NewInternalServerError(err))
	}
}

//	@title		Voltaserve Conversion
//	@version	3.0.0
//	@BasePath	/v3
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
		ErrorHandler: ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.MultipartBodyLengthLimitMB)),
	})

	router.NewVersionRouter().AppendRoutes(app)

	v3 := app.Group("v3")

	router.NewHealthRouter(router.HealthRouterOptions{
		Installer: installer,
	}).AppendRoutes(v3)

	router.NewPipelineRouter(router.NewPipelineRouterOptions{
		Scheduler: scheduler,
	}).AppendRoutes(v3)

	scheduler.Start()
	installer.Start()

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
