// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package test

import (
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

func TestRedis_SetAndGet(t *testing.T) {
	fileCache := cache.NewFileCache()
	opts := repo.NewFileOptions{
		ID:   helper.NewID(),
		Name: "file",
		Type: model.FileTypeFile,
	}
	if err := fileCache.Set(repo.NewFileWithOptions(opts)); err != nil {
		t.Fatal(err)
	}
	file, err := fileCache.Get(opts.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, file.GetID(), opts.ID)
	assert.Equal(t, file.GetName(), opts.Name)
	assert.Equal(t, file.GetType(), opts.Type)
}
