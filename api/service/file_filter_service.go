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
	"time"

	"github.com/reactivex/rxgo/v2"

	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FileFilterService interface {
	filterFolders(data []model.File) []model.File
	filterFiles(data []model.File) []model.File
	filterImages(data []model.File, userID string) []model.File
	filterPDFs(data []model.File, userID string) []model.File
	filterDocuments(data []model.File, userID string) []model.File
	filterVideos(data []model.File, userID string) []model.File
	filterTexts(data []model.File, userID string) []model.File
	filterOthers(data []model.File, userID string) []model.File
	filterWithQuery(data []model.File, opts FileQuery, parent model.File) ([]model.File, error)
}

type fileFilterService struct {
	fileRepo   repo.FileRepo
	fileGuard  guard.FileGuard
	fileMapper FileMapper
	fileIdent  *infra.FileIdentifier
}

func newFileFilterService() FileFilterService {
	return &fileFilterService{
		fileRepo:   repo.NewFileRepo(),
		fileGuard:  guard.NewFileGuard(),
		fileMapper: newFileMapper(),
		fileIdent:  infra.NewFileIdentifier(),
	}
}

func (svc *fileFilterService) filterFolders(data []model.File) []model.File {
	folders, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetType() == model.FileTypeFolder
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range folders {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterFiles(data []model.File) []model.File {
	files, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetType() == model.FileTypeFile
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range files {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterImages(data []model.File, userID string) []model.File {
	images, _ := rxgo.Just(data)().
		Filter(func(file interface{}) bool {
			f, err := svc.fileMapper.mapOne(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot != nil && f.Snapshot.Original == nil {
				return false
			}
			if f.Snapshot != nil && svc.fileIdent.IsImage(f.Snapshot.Original.Extension) {
				return true
			}
			return false
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range images {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterPDFs(data []model.File, userID string) []model.File {
	pdfs, _ := rxgo.Just(data)().
		Filter(func(file interface{}) bool {
			f, err := svc.fileMapper.mapOne(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot != nil && f.Snapshot.Original == nil {
				return false
			}
			if f.Snapshot != nil && svc.fileIdent.IsPDF(f.Snapshot.Original.Extension) {
				return true
			}
			return false
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range pdfs {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterDocuments(data []model.File, userID string) []model.File {
	documents, _ := rxgo.Just(data)().
		Filter(func(file interface{}) bool {
			f, err := svc.fileMapper.mapOne(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot != nil && f.Snapshot.Original == nil {
				return false
			}
			if f.Snapshot != nil && svc.fileIdent.IsOffice(f.Snapshot.Original.Extension) {
				return true
			}
			return false
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range documents {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterVideos(data []model.File, userID string) []model.File {
	videos, _ := rxgo.Just(data)().
		Filter(func(file interface{}) bool {
			f, err := svc.fileMapper.mapOne(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot != nil && f.Snapshot.Original == nil {
				return false
			}
			if f.Snapshot != nil && svc.fileIdent.IsVideo(f.Snapshot.Original.Extension) {
				return true
			}
			return false
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range videos {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterTexts(data []model.File, userID string) []model.File {
	texts, _ := rxgo.Just(data)().
		Filter(func(file interface{}) bool {
			f, err := svc.fileMapper.mapOne(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot != nil && f.Snapshot.Original == nil {
				return false
			}
			if f.Snapshot != nil && svc.fileIdent.IsPlainText(f.Snapshot.Original.Extension) {
				return true
			}
			return false
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range texts {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterOthers(data []model.File, userID string) []model.File {
	others, _ := rxgo.Just(data)().
		Filter(func(file interface{}) bool {
			f, err := svc.fileMapper.mapOne(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot != nil && f.Snapshot.Original == nil {
				return false
			}
			if f.Snapshot != nil &&
				!svc.fileIdent.IsImage(f.Snapshot.Original.Extension) &&
				!svc.fileIdent.IsPDF(f.Snapshot.Original.Extension) &&
				!svc.fileIdent.IsOffice(f.Snapshot.Original.Extension) &&
				!svc.fileIdent.IsVideo(f.Snapshot.Original.Extension) &&
				!svc.fileIdent.IsPlainText(f.Snapshot.Original.Extension) {
				return true
			}
			return false
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range others {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterWithQuery(data []model.File, opts FileQuery, parent model.File) ([]model.File, error) {
	filtered, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetWorkspaceID() == parent.GetWorkspaceID()
		}).
		Filter(func(v interface{}) bool {
			if opts.Type != nil {
				return v.(model.File).GetType() == *opts.Type
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			file := v.(model.File)
			res, err := svc.fileRepo.IsGrandChildOf(file.GetID(), parent.GetID())
			if err != nil {
				return false
			}
			return res
		}).
		Filter(func(v interface{}) bool {
			if opts.CreateTimeBefore != nil {
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return t.UnixMilli() >= *opts.CreateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.CreateTimeBefore != nil {
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return t.UnixMilli() <= *opts.CreateTimeBefore
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeAfter != nil {
				file := v.(model.File)
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return file.GetUpdateTime() != nil && t.UnixMilli() >= *opts.UpdateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeBefore != nil {
				file := v.(model.File)
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return file.GetUpdateTime() != nil && t.UnixMilli() <= *opts.UpdateTimeBefore
			} else {
				return true
			}
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range filtered {
		res = append(res, v.(model.File))
	}
	return res, nil
}
