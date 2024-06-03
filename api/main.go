package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/router"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/joho/godotenv"
)

// @title		Voltaserve API
// @version	2.0.0
// @BasePath	/v2
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
	}))

	v2 := app.Group("v2")

	health := router.NewHealthRouter()
	health.AppendRoutes(v2)

	filesGroup := v2.Group("files")
	files := router.NewFileRouter()
	files.AppendNonJWTRoutes(filesGroup)

	snapshotsGroup := v2.Group("snapshots")
	snapshots := router.NewSnapshotRouter()
	snapshots.AppendNonJWTRoutes(snapshotsGroup)

	insightsGroup := v2.Group("insights")
	insights := router.NewInsightsRouter()
	insights.AppendNonJWTRoutes(insightsGroup)

	mosaicGroup := v2.Group("mosaics")
	mosaic := router.NewMosaicRouter()
	mosaic.AppendNonJWTRoutes(mosaicGroup)

	watermarkGroup := v2.Group("watermarks")
	watermark := router.NewWatermarkRouter()
	watermark.AppendNonJWTRoutes(watermarkGroup)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Security.JWTSigningKey)},
	}))

	files.AppendRoutes(filesGroup)
	snapshots.AppendRoutes(snapshotsGroup)
	insights.AppendRoutes(insightsGroup)
	mosaic.AppendRoutes(mosaicGroup)
	watermark.AppendRoutes(watermarkGroup)

	invitations := router.NewInvitationRouter()
	invitations.AppendRoutes(v2.Group("invitations"))

	notifications := router.NewNotificationRouter()
	notifications.AppendRoutes(v2.Group("notifications"))

	organizations := router.NewOrganizationRouter()
	organizations.AppendRoutes(v2.Group("organizations"))

	tasks := router.NewTaskRouter()
	tasks.AppendRoutes(v2.Group("tasks"))

	storage := router.NewStorageRouter()
	storage.AppendRoutes(v2.Group("storage"))

	workspaces := router.NewWorkspaceRouter()
	workspaces.AppendRoutes(v2.Group("workspaces"))

	groups := router.NewGroupRouter()
	groups.AppendRoutes(v2.Group("groups"))

	users := router.NewUserRouter()
	users.AppendRoutes(v2.Group("users"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
