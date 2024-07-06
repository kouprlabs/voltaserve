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
	"net/http"
)

func NewInternalServerError(err error) *ErrorResponse {
	return NewErrorResponse(
		"internal_server_error",
		http.StatusInternalServerError,
		"Internal server error.",
		MsgSomethingWentWrong,
		err,
	)
}

func NewResourceNotFoundError(err error) *ErrorResponse {
	return &ErrorResponse{
		Code:        "resource_not_found",
		Status:      http.StatusNotFound,
		Message:     "Resource not found.",
		UserMessage: "The requested resource could not be found.",
		MoreInfo:    err.Error(),
		Err:         err,
	}
}
