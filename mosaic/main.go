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
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/mosaic/config"
	"github.com/kouprlabs/voltaserve/mosaic/errorpkg"
	"github.com/kouprlabs/voltaserve/mosaic/helper"
	"github.com/kouprlabs/voltaserve/mosaic/router"
)

// @title		Voltaserve Mosaic
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

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.MultipartBodyLengthLimitMB)),
	})

	versionRouter := router.NewVersionRouter()
	versionRouter.AppendRoutes(app)

	v3 := app.Group("v3")

	healthRouter := router.NewHealthRouter()
	healthRouter.AppendRoutes(v3)

	mosaicRouter := router.NewMosaicRouter()
	mosaicRouter.AppendRoutes(v3.Group("mosaics"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
