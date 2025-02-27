// Copyright (c) 2024 Mateusz Kaźmierczak.
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
	"github.com/gofiber/fiber/v2"
)

type VersionRouter struct{}

func NewVersionRouter() *VersionRouter {
	return &VersionRouter{}
}

func (r *VersionRouter) AppendRoutes(g fiber.Router) {
	g.Get("/version", r.Get)
}

// Get godoc
//
//	@Summary		Get
//	@Description	Get
//	@Tags			Version
//	@Id				version_get
//	@Produce		application/json
//	@Success		200	{object}	map[string]string
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/version [get]
func (r *VersionRouter) Get(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"version": "3.0.0",
	})
}
