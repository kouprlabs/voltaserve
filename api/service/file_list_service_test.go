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
	"github.com/golang/mock/gomock"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/kouprlabs/voltaserve/api/mocks"
	"github.com/kouprlabs/voltaserve/api/model"
)

func TestFileListService_Probe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := mocks.NewMockFileCache(ctrl)
	fileRepo := mocks.NewMockFileRepo(ctrl)
	fileGuard := mocks.NewMockFileGuard(ctrl)

	svc := &FileListService{
		fileCache: fileCache,
		fileRepo:  fileRepo,
		fileGuard: fileGuard,
	}

	userID := "user_id"
	folder := repo.NewFileWithOptions(repo.NewFileOptions{
		ID:   helper.NewID(),
		Type: model.FileTypeFolder,
	})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil)
	fileGuard.EXPECT().Authorize(userID, folder, model.PermissionViewer).Return(nil)
	fileRepo.EXPECT().CountChildren(folder.GetID()).Return(int64(100), nil)

	probe, err := svc.Probe(folder.GetID(), FileListOptions{Page: 1, Size: 10}, userID)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(100), probe.TotalElements)
		assert.Equal(t, uint64(10), probe.TotalPages)
	}
}

func TestFileListService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := mocks.NewMockFileCache(ctrl)
	fileRepo := mocks.NewMockFileRepo(ctrl)
	fileGuard := mocks.NewMockFileGuard(ctrl)
	workspaceRepo := mocks.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := mocks.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache: fileCache,
		fileRepo:  fileRepo,
		fileGuard: fileGuard,
		fileCoreSvc: &fileCoreService{
			fileRepo:  fileRepo,
			fileCache: fileCache,
			fileGuard: fileGuard,
		},
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	userID := "user_id"
	folder := repo.NewFileWithOptions(repo.NewFileOptions{
		ID:          helper.NewID(),
		WorkspaceID: helper.NewID(),
		Type:        model.FileTypeFolder,
	})
	file := repo.NewFileWithOptions(repo.NewFileOptions{
		ID:          helper.NewID(),
		WorkspaceID: folder.GetWorkspaceID(),
		ParentID:    helper.ToPtr(folder.GetID()),
		Type:        model.FileTypeFile,
	})
	workspace := repo.NewWorkspaceWithOptions(repo.NewWorkspaceOptions{
		ID: folder.GetWorkspaceID(),
	})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil).Times(1)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil).Times(1)
	fileRepo.EXPECT().FindChildrenIDs(folder.GetID()).Return([]string{file.GetID()}, nil)
	fileGuard.EXPECT().Authorize(userID, folder, model.PermissionViewer).Return(nil)
	fileGuard.EXPECT().IsAuthorized(userID, file, model.PermissionViewer).Return(true)
	workspaceRepo.EXPECT().Find(folder.GetWorkspaceID()).Return(workspace, nil)
	workspaceGuard.EXPECT().Authorize(userID, workspace, model.PermissionViewer).Return(nil)

	list, err := svc.List(folder.GetID(), FileListOptions{Page: 1, Size: 10}, userID)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(list.Data))
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.TotalElements)
	}
}
