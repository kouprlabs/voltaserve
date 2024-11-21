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
	"net/url"
	"strings"
)

func DecodeURIComponent(value string) string {
	res, err := url.PathUnescape(value)
	if err != nil {
		return ""
	}
	return res
}

func EncodeURIComponent(value string) string {
	encoded := url.QueryEscape(value)
	encoded = strings.ReplaceAll(encoded, "%2F", "/")
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	return encoded
}
