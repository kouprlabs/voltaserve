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
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/router"
)

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

	v := "v3"

	cfg := config.GetConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorpkg.ErrorHandler,
		BodyLimit:    int(helper.MegabyteToByte(cfg.Limits.FileUploadMB)),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.Security.CORSOrigins, ","),
	}))

	app.Use(func(c *fiber.Ctx) error {
		for _, route := range []struct {
			Path   string
			Method string
		}{
			{Path: "/version", Method: "GET"},
			{Path: "/" + v + "/health", Method: "GET"},
			{Path: "/" + v + "/workspaces/:id/bucket", Method: "GET"},
			{Path: "/" + v + "/workspaces/:id/image.:extension", Method: "GET"},
			{Path: "/" + v + "/organizations/:id/image.:extension", Method: "GET"},
			{Path: "/" + v + "/groups/:id/image.:extension", Method: "GET"},
			{Path: "/" + v + "/files/:id/original.:extension", Method: "GET"},
			{Path: "/" + v + "/files/:id/preview.:extension", Method: "GET"},
			{Path: "/" + v + "/files/:id/text.:extension", Method: "GET"},
			{Path: "/" + v + "/files/:id/ocr.:extension", Method: "GET"},
			{Path: "/" + v + "/files/:id/thumbnail.:extension", Method: "GET"},
			{Path: "/" + v + "/files/create_from_s3", Method: "POST"},
			{Path: "/" + v + "/files/:id/patch_from_s3", Method: "PATCH"},
			{Path: "/" + v + "/snapshots/:id", Method: "GET"},
			{Path: "/" + v + "/snapshots/:id", Method: "PATCH"},
			{Path: "/" + v + "/mosaics/:file_id/zoom_level/:zoom_level/row/:row/column/:column/extension/:extension", Method: "GET"},
			{Path: "/" + v + "/tasks", Method: "POST"},
			{Path: "/" + v + "/tasks/:id", Method: "DELETE"},
			{Path: "/" + v + "/tasks/:id", Method: "PATCH"},
			{Path: "/" + v + "/users/:id/picture.:extension", Method: "GET"},
			{Path: "/" + v + "/webhooks/users", Method: "POST"},
		} {
			if helper.MatchPath(route.Path, c.Path()) && c.Method() == route.Method {
				return c.Next()
			}
		}
		return jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(cfg.Security.JWTSigningKey)},
		})(c)
	})

	router.NewVersionRouter().AppendRoutes(app)

	group := app.Group(v)

	router.NewHealthRouter().AppendRoutes(group.Group("health"))
	router.NewWorkspaceRouter().AppendRoutes(group.Group("workspaces"))
	router.NewFileRouter().AppendRoutes(group.Group("files"))
	router.NewSnapshotRouter().AppendRoutes(group.Group("snapshots"))
	router.NewMosaicRouter().AppendRoutes(group.Group("mosaics"))
	router.NewTaskRouter().AppendRoutes(group.Group("tasks"))
	router.NewUserRouter().AppendRoutes(group.Group("users"))
	router.NewInvitationRouter().AppendRoutes(group.Group("invitations"))
	router.NewOrganizationRouter().AppendRoutes(group.Group("organizations"))
	router.NewStorageRouter().AppendRoutes(group.Group("storage"))
	router.NewGroupRouter().AppendRoutes(group.Group("groups"))
	router.NewEntityRouter().AppendRoutes(group.Group("entities"))
	router.NewWebhookRouter().AppendRoutes(group.Group("webhooks"))

	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
