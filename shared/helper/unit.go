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

func MegabyteToByte(mb int64) int64 {
	return mb * 1000000
}

func ByteToMegabyte(b int64) float64 {
	return float64(b) / 1000000
}
