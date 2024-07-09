// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package router

import "github.com/kouprlabs/voltaserve/api/service"

func IsValidSortBy(value string) bool {
	return value == "" ||
		value == service.SortByName ||
		value == service.SortByKind ||
		value == service.SortBySize ||
		value == service.SortByEmail ||
		value == service.SortByFullName ||
		value == service.SortByVersion ||
		value == service.SortByFrequency ||
		value == service.SortByDateCreated ||
		value == service.SortByDateModified
}

func IsValidSortOrder(value string) bool {
	return value == "" || value == service.SortOrderAsc || value == service.SortOrderDesc
}
