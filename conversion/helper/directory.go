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
	"os"
	"path/filepath"
)

func FindFileWithExtension(dirPath string, ext string) (*string, error) {
	var res string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ext {
			res = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if res == "" {
		return nil, nil
	}
	return &res, nil
}
