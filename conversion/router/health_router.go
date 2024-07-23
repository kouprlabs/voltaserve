// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package router

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/infra"
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
	g.Get("health", r.GetHealth)
}

// GetHealth godoc
//
//	@Summary		Get Health
//	@Description	Get Health
//	@Tags			Health
//	@Id				get_health
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/health [get]
func (r *HealthRouter) GetHealth(c *fiber.Ctx) error {
	if r.installer.IsRunning() {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if err := infra.NewS3Manager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if ok, err := api_client.NewHealthClient().Get(); err != nil || ok != "OK" {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	return c.SendString("OK")
}
