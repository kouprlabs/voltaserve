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

import "time"

func ToUTCString(value *string) string {
	if value == nil {
		return ""
	}
	parsedTime, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return ""
	}
	return parsedTime.Format(time.RFC1123)
}
