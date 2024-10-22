// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import "github.com/gabriel-vasile/mimetype"

func DetectMIMEFromPath(path string) string {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return "application/octet-stream"
	}
	return mime.String()
}

func DetectMIMEFromBytes(b []byte) string {
	return mimetype.Detect(b).String()
}
