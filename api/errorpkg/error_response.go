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

import "fmt"

type ErrorResponse struct {
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
	MoreInfo    string `json:"moreInfo"`
	Err         error  `json:"-"`
}

func NewErrorResponse(code string, status int, message string, userMessage string, err error) *ErrorResponse {
	return &ErrorResponse{
		Code:        code,
		Status:      status,
		Message:     message,
		UserMessage: userMessage,
		MoreInfo:    fmt.Sprintf("https://voltaserve.com/docs/api/errors/%s", code),
		Err:         err,
	}
}

func (err ErrorResponse) Error() string {
	return err.Code
}

func (err ErrorResponse) Unwrap() error {
	return err.Err
}
