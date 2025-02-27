// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package router

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/shared/client"
	_ "github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/infra"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/runtime"
)

type HealthRouter struct {
	installer *runtime.Installer
}

type HealthRouterOptions struct {
	Installer *runtime.Installer
}

func NewHealthRouter(opts HealthRouterOptions) *HealthRouter {
	return &HealthRouter{
		installer: opts.Installer,
	}
}

func (r *HealthRouter) AppendRoutes(g fiber.Router) {
	g.Get("/health", r.Get)
}

// Get godoc
//
//	@Summary		Get
//	@Description	Get
//	@Tags			Health
//	@Id				health_get
//	@Produce		text/plain
//	@Produce		application/json
//	@Success		200	{string}	string	"OK"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/health [get]
func (r *HealthRouter) Get(c *fiber.Ctx) error {
	if r.installer.IsRunning() {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if err := infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment).Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if ok, err := client.NewHealthClient(config.GetConfig().APIURL).Get(); err != nil || ok != "OK" {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	return c.SendString("OK")
}
