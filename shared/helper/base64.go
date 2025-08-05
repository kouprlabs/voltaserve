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
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Base64ToBytes(value string) ([]byte, error) {
	var withoutPrefix string
	if strings.Contains(value, ",") {
		withoutPrefix = strings.Split(value, ",")[1]
	} else {
		withoutPrefix = value
	}
	b, err := base64.StdEncoding.DecodeString(withoutPrefix)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Base64ToMIME(value string) string {
	if !strings.HasPrefix(value, "data:image/") {
		return ""
	}
	colonIndex := strings.Index(value, ":")
	semicolonIndex := strings.Index(value, ";")
	if colonIndex == -1 || semicolonIndex == -1 {
		return ""
	}
	return value[colonIndex+1 : semicolonIndex]
}

func Base64ToExtension(value string) string {
	switch Base64ToMIME(value) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}

func FileToBase64(path string) (*string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // Safe path
	if err != nil {
		return nil, err
	}
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}
	base64 := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
	return &base64, nil
}
