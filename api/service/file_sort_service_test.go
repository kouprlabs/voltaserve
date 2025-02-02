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

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

func TestFileSortService_SortByName(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &fileSortService{}

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

func TestFileSortService_SortBySize(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileMapper := NewMockFileMapper(ctrl)
	svc := &fileSortService{fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b"}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c"}),
	}

	fileMapper.EXPECT().mapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(100))}}}, nil).AnyTimes()
	fileMapper.EXPECT().mapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(200))}}}, nil).AnyTimes()
	fileMapper.EXPECT().mapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Size: helper.ToPtr(int64(50))}}}, nil).AnyTimes()

	sorted := svc.sortBySize(files, SortOrderAsc, "")
	assert.Equal(t, "file_c", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_b", sorted[2].GetID())

	sorted = svc.sortBySize(files, SortOrderDesc, "")
	assert.Equal(t, "file_b", sorted[0].GetID())
	assert.Equal(t, "file_a", sorted[1].GetID())
	assert.Equal(t, "file_c", sorted[2].GetID())
}

func TestFileSortService_SortByDateCreated(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &fileSortService{}

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

func TestFileSortService_SortByDateModified(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := &fileSortService{}

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

func TestFileSortService_SortByKind(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileFilterSvc := NewMockFileFilterService(ctrl)
	svc := &fileSortService{fileFilterSvc: fileFilterSvc}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_a", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_a", Type: model.FileTypeFolder}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_b", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_b", Type: model.FileTypeFolder}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file_c", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder_c", Type: model.FileTypeFolder}),
	}

	fileFilterSvc.EXPECT().filterFolders(files).Return([]model.File{files[1], files[3], files[5]})
	fileFilterSvc.EXPECT().filterFiles(files).Return([]model.File{files[0], files[2], files[4]})
	fileFilterSvc.EXPECT().filterImages(gomock.Any(), gomock.Any()).Return([]model.File{})
	fileFilterSvc.EXPECT().filterPDFs(gomock.Any(), gomock.Any()).Return([]model.File{})
	fileFilterSvc.EXPECT().filterDocuments(gomock.Any(), gomock.Any()).Return([]model.File{})
	fileFilterSvc.EXPECT().filterVideos(gomock.Any(), gomock.Any()).Return([]model.File{})
	fileFilterSvc.EXPECT().filterTexts(gomock.Any(), gomock.Any()).Return([]model.File{})
	fileFilterSvc.EXPECT().filterOthers(gomock.Any(), gomock.Any()).Return([]model.File{})

	sorted := svc.sortByKind(files, "")
	assert.Equal(t, "folder_a", sorted[0].GetID())
	assert.Equal(t, "folder_b", sorted[1].GetID())
	assert.Equal(t, "folder_c", sorted[2].GetID())
	assert.Equal(t, "file_a", sorted[3].GetID())
	assert.Equal(t, "file_b", sorted[4].GetID())
	assert.Equal(t, "file_c", sorted[5].GetID())
}
