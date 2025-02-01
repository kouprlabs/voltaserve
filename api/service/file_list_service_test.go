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

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

func TestFileListService_Probe(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)
	fileGuard := guard.NewMockFileGuard(ctrl)

	svc := &FileListService{
		fileCache: fileCache,
		fileRepo:  fileRepo,
		fileGuard: fileGuard,
	}

	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), Type: model.FileTypeFolder})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileRepo.EXPECT().CountChildren(folder.GetID()).Return(int64(100), nil)

	probe, err := svc.Probe(folder.GetID(), FileListOptions{Page: 1, Size: 10}, "")
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(100), probe.TotalElements)
		assert.Equal(t, uint64(10), probe.TotalPages)
	}
}

func TestFileListService_List(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)
	fileGuard := guard.NewMockFileGuard(ctrl)
	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)
	workspaceRepo := repo.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := guard.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache:      fileCache,
		fileRepo:       fileRepo,
		fileGuard:      fileGuard,
		fileCoreSvc:    fileCoreSvc,
		fileMapper:     fileMapper,
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	workspace := repo.NewWorkspaceWithOptions(repo.NewWorkspaceOptions{ID: "workspace"})
	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder, WorkspaceID: workspace.GetID()})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile, WorkspaceID: workspace.GetID()})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil).Times(1)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil).Times(1)
	fileRepo.EXPECT().FindChildrenIDs(folder.GetID()).Return([]string{file.GetID()}, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileCoreSvc.EXPECT().Authorize(gomock.Any(), []model.File{file}, model.PermissionViewer).Return([]model.File{file}, nil)
	fileMapper.EXPECT().MapMany([]model.File{file}, gomock.Any()).Return([]*File{{ID: file.GetID()}}, nil)
	workspaceRepo.EXPECT().Find(workspace.GetID()).Return(workspace, nil)
	workspaceGuard.EXPECT().Authorize(gomock.Any(), workspace, model.PermissionViewer).Return(nil)

	list, err := svc.List(folder.GetID(), FileListOptions{Page: 1, Size: 10}, "")
	if assert.NoError(t, err) {
		assert.Len(t, list.Data, 1)
		assert.Equal(t, list.Data[0].ID, file.GetID())
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.TotalElements)
		assert.Equal(t, uint64(1), list.Page)
		assert.Equal(t, uint64(10), list.Size)
	}
}

func TestFileListService_ListWithQuery(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)
	fileSearch := search.NewMockFileSearch(ctrl)
	fileGuard := guard.NewMockFileGuard(ctrl)
	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)
	workspaceRepo := repo.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := guard.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache:      fileCache,
		fileRepo:       fileRepo,
		fileSearch:     fileSearch,
		fileGuard:      fileGuard,
		fileMapper:     fileMapper,
		fileCoreSvc:    fileCoreSvc,
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	workspace := repo.NewWorkspaceWithOptions(repo.NewWorkspaceOptions{ID: "workspace"})
	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder, WorkspaceID: workspace.GetID()})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile, WorkspaceID: workspace.GetID()})
	query := FileQuery{Text: helper.ToPtr("search term")}

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)
	fileRepo.EXPECT().IsGrandChildOf(file.GetID(), folder.GetID()).Return(true, nil)
	fileSearch.EXPECT().Query(*query.Text, gomock.Any()).Return([]model.File{file}, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileCoreSvc.EXPECT().Authorize(gomock.Any(), []model.File{file}, model.PermissionViewer).Return([]model.File{file}, nil)
	fileMapper.EXPECT().MapMany([]model.File{file}, gomock.Any()).Return([]*File{{ID: file.GetID()}}, nil)
	workspaceRepo.EXPECT().Find(workspace.GetID()).Return(workspace, nil)
	workspaceGuard.EXPECT().Authorize(gomock.Any(), workspace, model.PermissionViewer).Return(nil)

	list, err := svc.List(folder.GetID(), FileListOptions{Page: 1, Size: 10, Query: &query}, "")
	if assert.NoError(t, err) {
		assert.Len(t, list.Data, 1)
		assert.Equal(t, list.Data[0].ID, file.GetID())
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.TotalElements)
		assert.Equal(t, uint64(1), list.Page)
		assert.Equal(t, uint64(10), list.Size)
	}
}

func TestFileListService_list(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &FileListService{fileCoreSvc: fileCoreSvc, fileMapper: fileMapper}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: "parent", Type: model.FileTypeFolder})
	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", Type: model.FileTypeFile}),
	}

	fileCoreSvc.EXPECT().Authorize(gomock.Any(), files, model.PermissionViewer).Return(files, nil)
	fileMapper.EXPECT().MapMany([]model.File{files[0], files[1]}, gomock.Any()).Return([]*File{{ID: files[0].GetID()}, {ID: files[1].GetID()}}, nil)

	list, err := svc.list(files, parent, FileListOptions{Page: 1, Size: 10, SortBy: SortByName, SortOrder: SortOrderAsc}, "")
	if assert.NoError(t, err) {
		assert.Len(t, list.Data, 2)
		assert.Equal(t, files[0].GetID(), list.Data[0].ID)
		assert.Equal(t, files[1].GetID(), list.Data[1].ID)
		assert.Equal(t, uint64(2), list.TotalElements)
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.Page)
		assert.Equal(t, uint64(10), list.Size)
	}
}

func TestFileListService_listWithQuery(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)
	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &FileListService{fileRepo: fileRepo, fileCoreSvc: fileCoreSvc, fileMapper: fileMapper}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: "parent", Type: model.FileTypeFolder})
	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{
			ID:         "file_a",
			Type:       model.FileTypeFile,
			CreateTime: "2022-12-01T00:00:00Z",
			UpdateTime: helper.ToPtr("2022-12-02T00:00:00Z"),
		}),
		repo.NewFileWithOptions(repo.NewFileOptions{
			ID:         "file_b",
			Type:       model.FileTypeFile,
			CreateTime: "2023-01-01T00:00:00Z",
			UpdateTime: helper.ToPtr("2023-01-01T00:00:00Z"),
		}),
	}
	query := FileQuery{
		Type:             helper.ToPtr(model.FileTypeFile),
		CreateTimeAfter:  helper.ToPtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()),
		CreateTimeBefore: helper.ToPtr(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli()),
		UpdateTimeAfter:  helper.ToPtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()),
		UpdateTimeBefore: helper.ToPtr(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).UnixMilli()),
	}

	fileRepo.EXPECT().IsGrandChildOf(files[0].GetID(), parent.GetID()).Return(true, nil)
	fileRepo.EXPECT().IsGrandChildOf(files[1].GetID(), parent.GetID()).Return(true, nil)
	fileCoreSvc.EXPECT().Authorize(gomock.Any(), []model.File{files[1]}, model.PermissionViewer).Return([]model.File{files[1]}, nil)
	fileMapper.EXPECT().MapMany([]model.File{files[1]}, gomock.Any()).Return([]*File{{ID: files[1].GetID()}}, nil)

	list, err := svc.list(files, parent, FileListOptions{Page: 1, Size: 10, Query: &query}, "")
	if assert.NoError(t, err) {
		assert.Len(t, list.Data, 1)
		assert.Equal(t, files[1].GetID(), list.Data[0].ID)
		assert.Equal(t, uint64(1), list.TotalElements)
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.Page)
		assert.Equal(t, uint64(10), list.Size)
	}
}

func TestFileListService_search(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileSearch := search.NewMockFileSearch(ctrl)

	svc := &FileListService{fileCache: fileCache, fileSearch: fileSearch}

	query := &FileQuery{Text: helper.ToPtr("search term"), Type: helper.ToPtr(model.FileTypeFile)}
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), Type: model.FileTypeFile})

	fileSearch.EXPECT().Query(*query.Text, gomock.Any()).Return([]model.File{file}, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)

	files, err := svc.search(query, repo.NewWorkspace())
	if assert.NoError(t, err) {
		assert.Len(t, files, 1)
		assert.Equal(t, file.GetID(), files[0].GetID())
	}
}

func TestFileListService_getChildren(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &FileListService{fileCache: fileCache, fileRepo: fileRepo}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID()})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID()})

	fileRepo.EXPECT().FindChildrenIDs(parent.GetID()).Return([]string{file.GetID()}, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)

	children, err := svc.getChildren(parent.GetID())
	if assert.NoError(t, err) {
		assert.Len(t, children, 1)
		assert.Equal(t, file.GetID(), children[0].GetID())
	}
}

func TestFileListService_sortByName(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &FileListService{}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), Name: "b"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), Name: "a"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), Name: "c"}),
	}

	sorted := svc.sortByName(files, SortOrderAsc)
	assert.Equal(t, "a", sorted[0].GetName())
	assert.Equal(t, "b", sorted[1].GetName())
	assert.Equal(t, "c", sorted[2].GetName())

	sorted = svc.sortByName(files, SortOrderDesc)
	assert.Equal(t, "c", sorted[0].GetName())
	assert.Equal(t, "b", sorted[1].GetName())
	assert.Equal(t, "a", sorted[2].GetName())
}

func TestFileListService_sortBySize(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMapper := NewMockFileMapper(ctrl)
	svc := &FileListService{
		fileMapper: fileMapper,
	}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c"}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(100))}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(200))}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(50))}}}, nil).AnyTimes()

	sorted := svc.sortBySize(files, SortOrderAsc, "")
	assert.Equal(t, "file_c", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_b", sorted[2].GetID())

	sorted = svc.sortBySize(files, SortOrderDesc, "")
	assert.Equal(t, "file_b", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_c", sorted[2].GetID())
}

func TestFileListService_sortByDateCreated(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &FileListService{}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", CreateTime: "2023-01-02T00:00:00Z"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", CreateTime: "2023-01-01T00:00:00Z"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c", CreateTime: "2023-01-03T00:00:00Z"}),
	}

	sorted := svc.sortByDateCreated(files, SortOrderAsc)
	assert.Equal(t, "file_b", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_c", sorted[2].GetID())

	sorted = svc.sortByDateCreated(files, SortOrderDesc)
	assert.Equal(t, "file_c", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_b", sorted[2].GetID())
}

func TestFileListService_sortByDateModified(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &FileListService{}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", UpdateTime: helper.ToPtr("2023-01-02T00:00:00Z")}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", UpdateTime: helper.ToPtr("2023-01-01T00:00:00Z")}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c", UpdateTime: helper.ToPtr("2023-01-03T00:00:00Z")}),
	}

	sorted := svc.sortByDateModified(files, SortOrderAsc)
	assert.Equal(t, "file_b", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_c", sorted[2].GetID())

	sorted = svc.sortByDateModified(files, SortOrderDesc)
	assert.Equal(t, "file_c", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_b", sorted[2].GetID())
}

func TestFileListService_sortByKind(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMapper := NewMockFileMapper(ctrl)
	svc := &FileListService{fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_a", Type: model.FileTypeFolder}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_b", Type: model.FileTypeFolder}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_c", Type: model.FileTypeFolder}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()

	sorted := svc.sortByKind(files, "")
	assert.Equal(t, "folder_a", sorted[0].GetID())
	assert.Equal(t, "folder_b", sorted[1].GetID())
	assert.Equal(t, "folder_c", sorted[2].GetID())
	assert.Equal(t, "file_a", sorted[3].GetID())
	assert.Equal(t, "file_b", sorted[4].GetID())
	assert.Equal(t, "file_c", sorted[5].GetID())
}

func TestFileListService_filterWithQuery(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &FileListService{fileRepo: fileRepo}

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

	filtered, err := svc.filterWithQuery([]model.File{file}, query, parent)
	if assert.NoError(t, err) {
		assert.Len(t, filtered, 1)
		assert.Equal(t, file.GetID(), filtered[0].GetID())
	}
}

func TestFileListService_paginate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &FileListService{}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID()}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID()}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID()}),
	}

	paged, totalElements, totalPages := svc.paginate(files, 1, 2)
	assert.Len(t, paged, 2)
	assert.Equal(t, uint64(3), totalElements)
	assert.Equal(t, uint64(2), totalPages)
}
