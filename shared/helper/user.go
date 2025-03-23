// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
)

func GetUserID(c *fiber.Ctx) (string, error) {
	user := c.Locals("user")
	if user == nil {
		return "", errorpkg.NewUnauthorizedUserError()
	}
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return claims.GetSubject()
}
