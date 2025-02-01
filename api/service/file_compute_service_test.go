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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

func TestFileComputeService_ComputeSize(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)
	fileGuard := guard.NewMockFileGuard(ctrl)

	svc := &FileComputeService{
		fileCache: fileCache,
		fileRepo:  fileRepo,
		fileGuard: fileGuard,
	}

	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file"})
	expectedSize := int64(1024)

	fileCache.EXPECT().Get(file.GetID()).Return(file, nil).Times(1)
	fileGuard.EXPECT().Authorize(gomock.Any(), file, model.PermissionViewer).Return(nil)
	fileRepo.EXPECT().ComputeSize(file.GetID()).Return(expectedSize, nil)

	size, err := svc.ComputeSize(file.GetID(), "")
	if assert.NoError(t, err) {
		assert.Equal(t, expectedSize, *size)
	}
}

func TestFileComputeService_Count(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)
	fileGuard := guard.NewMockFileGuard(ctrl)

	svc := &FileComputeService{
		fileCache: fileCache,
		fileRepo:  fileRepo,
		fileGuard: fileGuard,
	}

	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file"})
	expectedCount := int64(10)

	fileCache.EXPECT().Get(file.GetID()).Return(file, nil).Times(1)
	fileGuard.EXPECT().Authorize("", file, model.PermissionViewer).Return(nil).Times(1)
	fileRepo.EXPECT().CountItems(file.GetID()).Return(expectedCount, nil).Times(1)

	count, err := svc.Count(file.GetID(), "")
	if assert.NoError(t, err) {
		assert.Equal(t, expectedCount, *count)
	}
}
