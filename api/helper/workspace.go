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

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

func SlugFromWorkspace(id string, name string) string {
	return fmt.Sprintf("%s-%s", slug.Make(name), id)
}

func WorkspaceIDFromSlug(slug string) string {
	parts := strings.Split(slug, "-")
	return parts[len(parts)-1]
}
