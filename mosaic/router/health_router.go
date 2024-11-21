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

	"github.com/kouprlabs/voltaserve/mosaic/infra"
)

type HealthRouter struct{}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
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
	if err := infra.NewS3Manager().Connect(); err != nil {
		return c.SendStatus(http.StatusServiceUnavailable)
	}
	return c.SendString("OK")
}
