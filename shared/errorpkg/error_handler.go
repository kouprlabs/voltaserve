// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package errorpkg

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/kouprlabs/voltaserve/shared/logger"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var errorResponse *ErrorResponse
	if errors.As(err, &errorResponse) {
		return c.Status(errorResponse.Status).JSON(errorResponse)
	} else {
		logger.GetLogger().Error(err)
		return c.Status(http.StatusInternalServerError).JSON(NewInternalServerError(err))
	}
}
