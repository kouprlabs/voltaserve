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
	"net/http"
	"net/url"
	"path"
	"strings"
)

func GetTargetPath(req *http.Request) string {
	destination := req.Header.Get("Destination")
	if destination == "" {
		return ""
	}
	/* Check if the destination header is a full URL */
	if strings.HasPrefix(destination, "http://") || strings.HasPrefix(destination, "https://") {
		parsedURL, err := url.Parse(destination)
		if err != nil {
			return ""
		}
		return parsedURL.Path
	}
	/* Extract the path from the destination header */
	startIndex := strings.Index(destination, req.Host) + len(req.Host)
	if startIndex < len(req.Host) {
		return ""
	}
	return destination[startIndex:]
}

func Dirname(value string) string {
	trimmedValue := strings.TrimSuffix(value, "/")
	return path.Dir(trimmedValue)
}
