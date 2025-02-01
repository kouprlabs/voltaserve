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
	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{
			ID:         "file",
			Type:       model.FileTypeFile,
			CreateTime: "2023-01-01T00:00:00Z",
			UpdateTime: helper.ToPtr("2023-01-01T00:00:00Z"),
		}),
		repo.NewFileWithOptions(repo.NewFileOptions{
			ID:         "folder",
			Type:       model.FileTypeFolder,
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

	filtered, err := svc.FilterWithQuery(files, query, parent)
	if assert.NoError(t, err) {
		assert.Len(t, filtered, 1)
		assert.Equal(t, files[0].GetID(), filtered[0].GetID())
	}
}

func TestFileFilterService_FilterFolders(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder}),
	}

	filtered := svc.FilterFolders(files)
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[1].GetID(), filtered[0].GetID())
}

func TestFileFilterService_FilterFiles(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "file", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "folder", Type: model.FileTypeFolder}),
	}

	filtered := svc.FilterFiles(files)
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[0].GetID(), filtered[0].GetID())
}

func TestFileFilterService_FilterImages(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo, fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "image", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "pdf", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "text", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "office", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "video", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "other", Type: model.FileTypeFile}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[3], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".docx"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".mp4"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[5], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".zip"}}}, nil).AnyTimes()

	filtered := svc.FilterImages(files, "")
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[0].GetID(), filtered[0].GetID())
}

func TestFileFilterService_FilterPDFs(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo, fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "image", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "pdf", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "text", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "office", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "video", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "other", Type: model.FileTypeFile}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[3], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".docx"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".mp4"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[5], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".zip"}}}, nil).AnyTimes()

	filtered := svc.FilterPDFs(files, "")
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[1].GetID(), filtered[0].GetID())
}

func TestFileFilterService_FilterDocuments(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo, fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "image", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "pdf", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "text", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "office", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "video", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "other", Type: model.FileTypeFile}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[3], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".docx"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".mp4"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[5], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".zip"}}}, nil).AnyTimes()

	filtered := svc.FilterDocuments(files, "")
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[3].GetID(), filtered[0].GetID())
}

func TestFileFilterService_FilterVideos(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo, fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "image", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "pdf", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "text", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "office", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "video", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "other", Type: model.FileTypeFile}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[3], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".docx"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".mp4"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[5], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".zip"}}}, nil).AnyTimes()

	filtered := svc.FilterVideos(files, "")
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[4].GetID(), filtered[0].GetID())
}

func TestFileFilterService_FilterOthers(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileRepo := repo.NewMockFileRepo(ctrl)
	fileMapper := NewMockFileMapper(ctrl)

	svc := &fileFilterService{fileRepo: fileRepo, fileMapper: fileMapper}

	files := []model.File{
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "image", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "pdf", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "text", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "office", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "video", Type: model.FileTypeFile}),
		repo.NewFileWithOptions(repo.NewFileOptions{ID: "other", Type: model.FileTypeFile}),
	}

	fileMapper.EXPECT().MapOne(files[0], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".jpg"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[1], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".pdf"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[2], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".txt"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[3], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".docx"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[4], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".mp4"}}}, nil).AnyTimes()
	fileMapper.EXPECT().MapOne(files[5], gomock.Any()).Return(&File{Snapshot: &Snapshot{Original: &Download{Extension: ".zip"}}}, nil).AnyTimes()

	filtered := svc.FilterOthers(files, "")
	assert.Len(t, filtered, 1)
	assert.Equal(t, files[5].GetID(), filtered[0].GetID())
}
