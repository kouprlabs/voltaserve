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
	"fmt"
	"path/filepath"
)

func UniqueFilename(name string) string {
	return fmt.Sprintf("%s %s%s", FilenameWithoutExtension(name), NewID(), filepath.Ext(name))
}

func FilenameWithoutExtension(name string) string {
	withExt := filepath.Base(name)
	return withExt[0 : len(withExt)-len(filepath.Ext(name))]
}

func FilenameWithNewExtension(name string, ext string) string {
	withExt := filepath.Base(name)
	withoutExt := withExt[0 : len(withExt)-len(filepath.Ext(name))]
	if ext == "" {
		return withoutExt
	} else {
		return fmt.Sprintf("%s.%s", withoutExt, ext)
	}
}
