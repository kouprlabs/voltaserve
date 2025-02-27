// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package handler

import (
	"errors"
	"net/http"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"

	"github.com/kouprlabs/voltaserve/webdav/logger"
)

func handleError(err error, w http.ResponseWriter) {
	var errorResponse *errorpkg.ErrorResponse
	switch {
	case errors.As(err, &errorResponse):
		w.WriteHeader(errorResponse.Status)
		if _, err := w.Write([]byte(errorResponse.UserMessage)); err != nil {
			logger.GetLogger().Error(err)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Internal Server Error")); err != nil {
			return
		}
	}
	logger.GetLogger().Error(err)
}
