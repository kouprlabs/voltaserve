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
	fileSortSvc := NewMockFileSortService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)
	workspaceRepo := repo.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := guard.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache:      fileCache,
		fileRepo:       fileRepo,
		fileGuard:      fileGuard,
		fileCoreSvc:    fileCoreSvc,
		fileSortSvc:    fileSortSvc,
		fileMapper:     fileMapper,
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	workspace := repo.NewWorkspaceWithOptions(repo.NewWorkspaceOptions{ID: "workspace"})
	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder, WorkspaceID: workspace.GetID()})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile, WorkspaceID: workspace.GetID()})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)
	fileRepo.EXPECT().FindChildrenIDs(folder.GetID()).Return([]string{file.GetID()}, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileCoreSvc.EXPECT().authorize(gomock.Any(), []model.File{file}, model.PermissionViewer).Return([]model.File{file}, nil)
	fileSortSvc.EXPECT().sort([]model.File{file}, gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.File{file})
	fileMapper.EXPECT().mapMany([]model.File{file}, gomock.Any()).Return([]*File{{ID: file.GetID()}}, nil)
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
	fileSearch := search.NewMockFileSearch(ctrl)
	fileGuard := guard.NewMockFileGuard(ctrl)
	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileFilterSvc := NewMockFileFilterService(ctrl)
	fileSortSvc := NewMockFileSortService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)
	workspaceRepo := repo.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := guard.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache:      fileCache,
		fileSearch:     fileSearch,
		fileGuard:      fileGuard,
		fileMapper:     fileMapper,
		fileCoreSvc:    fileCoreSvc,
		fileFilterSvc:  fileFilterSvc,
		fileSortSvc:    fileSortSvc,
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	workspace := repo.NewWorkspaceWithOptions(repo.NewWorkspaceOptions{ID: "workspace"})
	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder, WorkspaceID: workspace.GetID()})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile, WorkspaceID: workspace.GetID()})
	query := FileQuery{Text: helper.ToPtr("search term")}

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)
	fileSearch.EXPECT().Query(*query.Text, gomock.Any()).Return([]model.File{file}, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileCoreSvc.EXPECT().authorize(gomock.Any(), []model.File{file}, model.PermissionViewer).Return([]model.File{file}, nil)
	fileFilterSvc.EXPECT().filterWithQuery([]model.File{file}, FileQuery{Text: query.Text}, folder).Return([]model.File{file}, nil)
	fileSortSvc.EXPECT().sort([]model.File{file}, gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.File{file})
	fileMapper.EXPECT().mapMany([]model.File{file}, gomock.Any()).Return([]*File{{ID: file.GetID()}}, nil)
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

func TestFileListService_createList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileSortSvc := NewMockFileSortService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &FileListService{fileCoreSvc: fileCoreSvc, fileSortSvc: fileSortSvc, fileMapper: fileMapper}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: "parent", Type: model.FileTypeFolder})
	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", Type: model.FileTypeFile}),
	}

	fileCoreSvc.EXPECT().authorize(gomock.Any(), files, model.PermissionViewer).Return(files, nil)
	fileSortSvc.EXPECT().sort(files, gomock.Any(), gomock.Any(), gomock.Any()).Return(files)
	fileMapper.EXPECT().mapMany([]model.File{files[0], files[1]}, gomock.Any()).Return([]*File{{ID: files[0].GetID()}, {ID: files[1].GetID()}}, nil)

	list, err := svc.createList(files, parent, FileListOptions{Page: 1, Size: 10, SortBy: SortByName, SortOrder: SortOrderAsc}, "")
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

func TestFileListService_createListWithQuery(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCoreSvc := NewMockFileCoreService(ctrl)
	fileSortSvc := NewMockFileSortService(ctrl)
	fileFilterSvc := NewMockFileFilterService(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &FileListService{
		fileCoreSvc:   fileCoreSvc,
		fileSortSvc:   fileSortSvc,
		fileFilterSvc: fileFilterSvc,
		fileMapper:    fileMapper,
	}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: "parent", Type: model.FileTypeFolder})
	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", Type: model.FileTypeFile}),
	}
	query := FileQuery{}

	fileCoreSvc.EXPECT().authorize(gomock.Any(), []model.File{files[1]}, model.PermissionViewer).Return([]model.File{files[1]}, nil)
	fileSortSvc.EXPECT().sort([]model.File{files[1]}, gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.File{files[1]})
	fileFilterSvc.EXPECT().filterWithQuery(files, query, parent).Return([]model.File{files[1]}, nil)
	fileMapper.EXPECT().mapMany([]model.File{files[1]}, gomock.Any()).Return([]*File{{ID: files[1].GetID()}}, nil)

	list, err := svc.createList(files, parent, FileListOptions{Page: 1, Size: 10, Query: &query}, "")
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
