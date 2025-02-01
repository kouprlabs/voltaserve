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
	mapper := NewMockFileMapper(ctrl)
	workspaceRepo := repo.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := guard.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache:      fileCache,
		fileRepo:       fileRepo,
		fileGuard:      fileGuard,
		fileCoreSvc:    &fileCoreService{fileRepo: fileRepo, fileCache: fileCache, fileGuard: fileGuard},
		fileMapper:     mapper,
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil).Times(1)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil).Times(1)
	fileRepo.EXPECT().FindChildrenIDs(folder.GetID()).Return([]string{file.GetID()}, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileGuard.EXPECT().IsAuthorized(gomock.Any(), file, model.PermissionViewer).Return(true)
	mapper.EXPECT().MapMany(gomock.Any(), gomock.Any()).Return([]*File{{ID: file.GetID()}}, nil)
	workspaceRepo.EXPECT().Find(gomock.Any()).Return(repo.NewWorkspace(), nil)
	workspaceGuard.EXPECT().Authorize(gomock.Any(), gomock.Any(), model.PermissionViewer).Return(nil)

	list, err := svc.List(folder.GetID(), FileListOptions{Page: 1, Size: 10}, "")
	if assert.NoError(t, err) {
		assert.Len(t, list.Data, 1)
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.TotalElements)
	}
}

func TestFileListService_Search(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileSearch := search.NewMockFileSearch(ctrl)

	svc := &FileListService{fileCache: fileCache, fileSearch: fileSearch}

	workspace := repo.NewWorkspaceWithOptions(repo.NewWorkspaceOptions{ID: helper.NewID()})
	query := &FileQuery{Text: helper.ToPtr("search term"), Type: helper.ToPtr(model.FileTypeFile)}
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), WorkspaceID: workspace.GetID(), Type: model.FileTypeFile})

	fileSearch.EXPECT().Query(*query.Text, gomock.Any()).Return([]model.File{file}, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)

	files, err := svc.search(query, workspace)
	if assert.NoError(t, err) {
		assert.Len(t, files, 1)
		assert.Equal(t, file.GetID(), files[0].GetID())
	}
}

func TestFileListService_GetChildren(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &FileListService{fileCache: fileCache, fileRepo: fileRepo}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID()})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: helper.NewID(), ParentID: helper.ToPtr(parent.GetID())})

	fileRepo.EXPECT().FindChildrenIDs(parent.GetID()).Return([]string{file.GetID()}, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)

	children, err := svc.getChildren(parent.GetID())
	if assert.NoError(t, err) {
		assert.Len(t, children, 1)
		assert.Equal(t, file.GetID(), children[0].GetID())
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
	mapper := NewMockFileMapper(ctrl)
	workspaceRepo := repo.NewMockWorkspaceRepo(ctrl)
	workspaceGuard := guard.NewMockWorkspaceGuard(ctrl)

	svc := &FileListService{
		fileCache:      fileCache,
		fileRepo:       fileRepo,
		fileSearch:     fileSearch,
		fileGuard:      fileGuard,
		fileMapper:     mapper,
		fileCoreSvc:    &fileCoreService{fileRepo: fileRepo, fileCache: fileCache, fileGuard: fileGuard},
		workspaceRepo:  workspaceRepo,
		workspaceGuard: workspaceGuard,
	}

	folder := repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder})
	file := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile})

	fileCache.EXPECT().Get(folder.GetID()).Return(folder, nil)
	fileCache.EXPECT().Get(file.GetID()).Return(file, nil)
	fileRepo.EXPECT().IsGrandChildOf(file.GetID(), folder.GetID()).Return(true, nil)
	fileSearch.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]model.File{file}, nil)
	fileGuard.EXPECT().Authorize(gomock.Any(), folder, model.PermissionViewer).Return(nil)
	fileGuard.EXPECT().IsAuthorized(gomock.Any(), file, model.PermissionViewer).Return(true)
	mapper.EXPECT().MapMany(gomock.Any(), gomock.Any()).Return([]*File{{ID: file.GetID()}}, nil)
	workspaceRepo.EXPECT().Find(gomock.Any()).Return(repo.NewWorkspace(), nil)
	workspaceGuard.EXPECT().Authorize(gomock.Any(), gomock.Any(), model.PermissionViewer).Return(nil)

	list, err := svc.List(folder.GetID(), FileListOptions{Page: 1, Size: 10, Query: &FileQuery{Text: helper.ToPtr("search term")}}, "")
	if assert.NoError(t, err) {
		assert.Len(t, list.Data, 1)
		assert.Equal(t, uint64(1), list.TotalPages)
		assert.Equal(t, uint64(1), list.TotalElements)
	}
}

func TestFileListService_SortByName(t *testing.T) {
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

func TestFileListService_SortBySize(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mapper := NewMockFileMapper(ctrl)
	svc := &FileListService{
		fileMapper: mapper,
	}

	fileA := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a"})
	fileB := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b"})
	fileC := repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c"})
	files := []model.File{fileA, fileB, fileC}

	mapper.EXPECT().MapOne(fileA, gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(100))}}}, nil).AnyTimes()
	mapper.EXPECT().MapOne(fileB, gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(200))}}}, nil).AnyTimes()
	mapper.EXPECT().MapOne(fileC, gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(50))}}}, nil).AnyTimes()

	sorted := svc.sortBySize(files, SortOrderAsc, "")
	assert.Equal(t, "file_c", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_b", sorted[2].GetID())

	sorted = svc.sortBySize(files, SortOrderDesc, "")
	assert.Equal(t, "file_b", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_c", sorted[2].GetID())
}

func TestFileListService_SortByDateCreated(t *testing.T) {
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

func TestFileListService_SortByDateModified(t *testing.T) {
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

func TestFileListService_SortByKind(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mapper := NewMockFileMapper(ctrl)
	svc := &FileListService{fileMapper: mapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_a", Type: model.FileTypeFolder}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_b", Type: model.FileTypeFolder}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_c", Type: model.FileTypeFolder}),
	}

	mapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	mapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	mapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()

	sorted := svc.sortByKind(files, "")
	assert.Equal(t, "folder_a", sorted[0].GetID())
	assert.Equal(t, "folder_b", sorted[1].GetID())
	assert.Equal(t, "folder_c", sorted[2].GetID())
	assert.Equal(t, "file_a", sorted[3].GetID())
	assert.Equal(t, "file_b", sorted[4].GetID())
	assert.Equal(t, "file_c", sorted[5].GetID())
}

func TestFileListService_FilterWithQuery(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileCache := cache.NewMockFileCache(ctrl)
	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &FileListService{fileCache: fileCache, fileRepo: fileRepo}

	parent := repo.NewFileWithOptions(repo.NewFileOptions{ID: "parent", Type: model.FileTypeFolder})
	file := repo.NewFileWithOptions(repo.NewFileOptions{
		ID:         "file",
		ParentID:   helper.ToPtr(parent.GetID()),
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

func TestFileListService_Paginate(t *testing.T) {
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
