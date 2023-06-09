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

	jwtware "github.com/gofiber/jwt/v3"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

//	@title		Voltaserve API
//	@version	1.0.0
//	@BasePath	/v1
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

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.MultipartBodyLengthLimitMB)),
	})

	v1 := app.Group("/v1")

	app.Get("v1/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
	}))

	f := v1.Group("files")

	fileDownloads := router.NewFileDownloadRouter()
	fileDownloads.AppendRoutes(f)

	conversionWebhook := router.NewConversionWebhookRouter()
	conversionWebhook.AppendRoutes(f)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.Security.JWTSigningKey),
	}))

	files := router.NewFileRouter()
	files.AppendRoutes(f)

	invitations := router.NewInvitationRouter()
	invitations.AppendRoutes(v1.Group("invitations"))

	notifications := router.NewNotificationRouter()
	notifications.AppendRoutes(v1.Group("notifications"))

	organizations := router.NewOrganizationRouter()
	organizations.AppendRoutes(v1.Group("organizations"))

	storage := router.NewStorageRouter()
	storage.AppendRoutes(v1.Group("storage"))

	workspaces := router.NewWorkspaceRouter()
	workspaces.AppendRoutes(v1.Group("workspaces"))

	groups := router.NewGroupRouter()
	groups.AppendRoutes(v1.Group("groups"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
