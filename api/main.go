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

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.MultipartBodyLengthLimitMB)),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
	}))

	v1 := app.Group("v1")

	health := router.NewHealthRouter()
	health.AppendRoutes(v1)

	f := v1.Group("files")

	downloads := router.NewDownloadsRouter(router.NewDownloadsRouterOptions{})
	downloads.AppendNonJWTRoutes(f)

	conversionWebhook := router.NewConversionWebhookRouter(router.NewConversionWebhookRouterOptions{})
	conversionWebhook.AppendInternalRoutes(f)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Security.JWTSigningKey)},
	}))

	files := router.NewFileRouter(router.NewFileRouterOptions{})
	files.AppendRoutes(f)

	snapshots := router.NewSnapshotRouter(router.NewSnapshotRouterOptions{})
	snapshots.AppendRoutes(f)

	invitations := router.NewInvitationRouter(router.NewInvitationRouterOptions{})
	invitations.AppendRoutes(v1.Group("invitations"))

	notifications := router.NewNotificationRouter(router.NewNotificationRouterOptions{})
	notifications.AppendRoutes(v1.Group("notifications"))

	organizations := router.NewOrganizationRouter(router.NewOrganizationRouterOptions{})
	organizations.AppendRoutes(v1.Group("organizations"))

	storage := router.NewStorageRouter(router.NewStorageRouterOptions{})
	storage.AppendRoutes(v1.Group("storage"))

	workspaces := router.NewWorkspaceRouter(router.NewWorkspaceRouterOptions{})
	workspaces.AppendRoutes(v1.Group("workspaces"))

	groups := router.NewGroupRouter(router.NewGroupRouterOptions{})
	groups.AppendRoutes(v1.Group("groups"))

	users := router.NewUserRouter(router.NewUserRouterOptions{})
	users.AppendRoutes(v1.Group("users"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
