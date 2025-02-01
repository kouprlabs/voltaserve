// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

func TestFileFilterService_FilterWithQuery(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: "parent", Type: model.FileTypeFolder})
	file := repo.NewFileWithOptions(repo.NewFileOptions{
		ID:         "file",
		Type:       model.FileTypeFile,
		CreateTime: "2023-01-01T00:00:00Z",
		UpdateTime: helper.ToPtr("2023-01-01T00:00:00Z"),
	})
	query := FileQuery{
		Type:             helper.ToPtr(model.FileTypeFile),
		CreateTimeAfter:  helper.ToPtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()),
		CreateTimeBefore: helper.ToPtr(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli()),
		UpdateTimeAfter:  helper.ToPtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()),
		UpdateTimeBefore: helper.ToPtr(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli()),
	}

	fileRepo.EXPECT().IsGrandChildOf(file.GetID(), parent.GetID()).Return(true, nil)

	filtered, err := svc.FilterWithQuery([]model.File{file}, query, parent)
	if assert.NoError(t, err) {
		assert.Len(t, filtered, 1)
		assert.Equal(t, file.GetID(), filtered[0].GetID())
	}
}
