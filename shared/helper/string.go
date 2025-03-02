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
	"regexp"
	"strings"
)

func RemoveNonNumeric(s string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}

func RemoveWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
