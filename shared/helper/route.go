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

import "strings"

func MatchPath(pattern, path string) bool {
	// Split pattern and path into segments
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(path, "/")
	if len(patternSegments) != len(pathSegments) {
		return false
	}
	for i := range patternSegments {
		patternPart := patternSegments[i]
		pathPart := pathSegments[i]
		// If the pattern part is a dynamic segment (starts with ":"), skip comparison
		if strings.HasPrefix(patternPart, ":") {
			continue
		}
		// If the pattern part contains a wildcard (e.g., "thumbnail.:extension"), handle it
		if strings.Contains(patternPart, ":") {
			// Split the pattern part into subparts (e.g., "thumbnail.:extension" -> ["thumbnail", ":extension"])
			patternSubparts := strings.Split(patternPart, ".")
			pathSubparts := strings.Split(pathPart, ".")
			if len(patternSubparts) != len(pathSubparts) {
				return false
			}
			for j := range patternSubparts {
				if !strings.HasPrefix(patternSubparts[j], ":") && patternSubparts[j] != pathSubparts[j] {
					return false
				}
			}
			continue
		}
		// If the pattern part is not dynamic and doesn't match the path part, return false
		if patternPart != pathPart {
			return false
		}
	}
	return true
}
