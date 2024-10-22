// Copyright 2024 Mateusz Kaźmierczak.
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
	"github.com/gofiber/fiber/v2"
)

type VersionRouter struct{}

func NewVersionRouter() *VersionRouter {
	return &VersionRouter{}
}

func (r *VersionRouter) AppendRoutes(g fiber.Router) {
	g.Get("/version", r.GetVersion)
}

// GetVersion godoc
//
//	@Summary		Get
//	@Description	Get
//	@Tags			Version
//	@Id				get_version
//	@Produce		json
//	@Success		200	{string}	string	"{Version}"
//	@Failure		503	{object}	errorpkg.ErrorResponse
//	@Router			/health [get]
func (r *VersionRouter) GetVersion(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"version": "3.0.0",
	})
}
