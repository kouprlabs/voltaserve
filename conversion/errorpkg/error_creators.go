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
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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

func NewMissingQueryParamError(param string) *ErrorResponse {
	return NewErrorResponse(
		"missing_query_param",
		http.StatusBadRequest,
		fmt.Sprintf("Query param '%s' is required.", param),
		MsgInvalidRequest,
		nil,
	)
}

func NewRequestBodyValidationError(err error) *ErrorResponse {
	var fields []string
	for _, e := range err.(validator.ValidationErrors) {
		fields = append(fields, e.Field())
	}
	return NewErrorResponse(
		"request_validation_error",
		http.StatusBadRequest,
		fmt.Sprintf("Failed validation for the following fields: %s.", strings.Join(fields, ",")),
		MsgInvalidRequest,
		err,
	)
}

func NewInvalidAPIKeyError() *ErrorResponse {
	return NewErrorResponse(
		"invalid_api_key",
		http.StatusUnauthorized,
		"Invalid API key.",
		"The API key is either missing or invalid.",
		nil,
	)
}
