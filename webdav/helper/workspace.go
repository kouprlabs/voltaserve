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

import "strings"

func ExtractWorkspaceIDFromPath(path string) string {
	slashParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	dashParts := strings.Split(slashParts[0], "-")
	if len(dashParts) > 1 {
		return dashParts[len(dashParts)-1]
	}
	return ""
}
