// Package helper contains various helper functions.
//
// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package helper

import "slices"

func Includes(arr []string, str string) bool {
	slices.Contains(arr, str)
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
