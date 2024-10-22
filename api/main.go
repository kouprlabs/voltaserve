// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package main

import (
	"fmt"
	"os"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/router"
)

// @title		Voltaserve API
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
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.FileUploadMB)),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
	}))

	v3 := app.Group("v3")

	health := router.NewHealthRouter()
	health.AppendRoutes(v3)

	version := router.NewVersionRouter()
	version.AppendRoutes(app)

	filesGroup := v3.Group("files")
	files := router.NewFileRouter()
	files.AppendNonJWTRoutes(filesGroup)

	snapshotsGroup := v3.Group("snapshots")
	snapshots := router.NewSnapshotRouter()
	snapshots.AppendNonJWTRoutes(snapshotsGroup)

	insightsGroup := v3.Group("insights")
	insights := router.NewInsightsRouter()
	insights.AppendNonJWTRoutes(insightsGroup)

	mosaicGroup := v3.Group("mosaics")
	mosaic := router.NewMosaicRouter()
	mosaic.AppendNonJWTRoutes(mosaicGroup)

	tasksGroup := v3.Group("tasks")
	tasks := router.NewTaskRouter()
	tasks.AppendNonJWTRoutes(tasksGroup)

	usersGroup := v3.Group("users")
	users := router.NewUserRouter()
	users.AppendNonJWTRoutes(usersGroup)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Security.JWTSigningKey)},
	}))

	files.AppendRoutes(filesGroup)
	snapshots.AppendRoutes(snapshotsGroup)
	insights.AppendRoutes(insightsGroup)
	mosaic.AppendRoutes(mosaicGroup)
	tasks.AppendRoutes(tasksGroup)
	users.AppendRoutes(usersGroup)

	invitations := router.NewInvitationRouter()
	invitations.AppendRoutes(v3.Group("invitations"))

	organizations := router.NewOrganizationRouter()
	organizations.AppendRoutes(v3.Group("organizations"))

	storage := router.NewStorageRouter()
	storage.AppendRoutes(v3.Group("storage"))

	workspaces := router.NewWorkspaceRouter()
	workspaces.AppendRoutes(v3.Group("workspaces"))

	groups := router.NewGroupRouter()
	groups.AppendRoutes(v3.Group("groups"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
