package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helpers"
	"voltaserve/router"

	jwtware "github.com/gofiber/jwt/v3"
	log "github.com/sirupsen/logrus"
)

// @title    Voltaserve API
// @version  1.0.0
// @BasePath /v1
func main() {
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	settings := config.GetConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helpers.MegabyteToByte(settings.Limits.MultipartBodyLengthLimitMb)),
	})

	v1 := app.Group("/v1")

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(settings.Security.CorsOrigins, ","),
	}))

	f := v1.Group("files")
	fileDownloads := router.NewFileDownloadRouter()
	fileDownloads.AppendRoutes(f)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(settings.Security.JwtSigningKey),
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

	url, err := url.Parse(settings.Url)
	if err != nil {
		panic(err)
	}
	if err := app.Listen(":" + url.Port()); err != nil {
		panic(err)
	}
}
