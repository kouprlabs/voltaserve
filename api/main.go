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
	"fmt"
	"net/http"
	"os"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
	"github.com/kouprlabs/voltaserve/api/router"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var errorResponse *errorpkg.ErrorResponse
	if errors.As(err, &errorResponse) {
		return c.Status(errorResponse.Status).JSON(errorResponse)
	} else {
		logger.GetLogger().Error(err)
		return c.Status(http.StatusInternalServerError).JSON(errorpkg.NewInternalServerError(err))
	}
}

//	@title		Voltaserve API
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

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.FileUploadMB)),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
	}))

	versionRouter := router.NewVersionRouter()
	versionRouter.AppendRoutes(app)

	v3 := app.Group("v3")

	healthRouter := router.NewHealthRouter()
	healthRouter.AppendRoutes(v3)

	workspaceGroup := v3.Group("workspaces")
	workspaceRouter := router.NewWorkspaceRouter()
	workspaceRouter.AppendNonJWTRoutes(workspaceGroup)

	fileGroup := v3.Group("files")
	fileRouter := router.NewFileRouter()
	fileRouter.AppendNonJWTRoutes(fileGroup)

	snapshotGroup := v3.Group("snapshots")
	snapshotRouter := router.NewSnapshotRouter()
	snapshotRouter.AppendNonJWTRoutes(snapshotGroup)

	mosaicGroup := v3.Group("mosaics")
	mosaicRouter := router.NewMosaicRouter()
	mosaicRouter.AppendNonJWTRoutes(mosaicGroup)

	taskGroup := v3.Group("tasks")
	taskRouter := router.NewTaskRouter()
	taskRouter.AppendNonJWTRoutes(taskGroup)

	userGroup := v3.Group("users")
	userRouter := router.NewUserRouter()
	userRouter.AppendNonJWTRoutes(userGroup)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Security.JWTSigningKey)},
	}))

	workspaceRouter.AppendRoutes(workspaceGroup)
	fileRouter.AppendRoutes(fileGroup)
	snapshotRouter.AppendRoutes(snapshotGroup)
	mosaicRouter.AppendRoutes(mosaicGroup)
	taskRouter.AppendRoutes(taskGroup)
	userRouter.AppendRoutes(userGroup)

	invitationRouter := router.NewInvitationRouter()
	invitationRouter.AppendRoutes(v3.Group("invitations"))

	orgRouter := router.NewOrganizationRouter()
	orgRouter.AppendRoutes(v3.Group("organizations"))

	storageRouter := router.NewStorageRouter()
	storageRouter.AppendRoutes(v3.Group("storage"))

	groupRouter := router.NewGroupRouter()
	groupRouter.AppendRoutes(v3.Group("groups"))

	entityRouter := router.NewEntityRouter()
	entityRouter.AppendRoutes(v3.Group("entities"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
