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

	"github.com/kouprlabs/voltaserve/api/infra"
)

type HealthRouter struct{}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

func (r *HealthRouter) AppendRoutes(g fiber.Router) {
	g.Get("/health", r.Check)
}

// Check godoc
//
//	@Summary		Check
//	@Description	Check
//	@Tags			Health
//	@Id				health_check
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/health [get]
func (r *HealthRouter) Check(c *fiber.Ctx) error {
	if err := infra.NewPostgresManager().Connect(true); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if err := infra.NewRedisManager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	if err := infra.NewS3Manager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	return c.SendString("OK")
}
