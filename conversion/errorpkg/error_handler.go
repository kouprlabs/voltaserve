// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package errorpkg

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var e *ErrorResponse
	if errors.As(err, &e) {
		var v *ErrorResponse
		errors.As(err, &v)
		return c.Status(v.Status).JSON(v)
	} else {
		log.Error(err)
		return c.Status(http.StatusInternalServerError).JSON(NewInternalServerError(err))
	}
}
