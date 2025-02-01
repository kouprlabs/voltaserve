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
	"sort"
	"time"

	"github.com/kouprlabs/voltaserve/api/model"
)

type FileSortService interface {
	Sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File
	SortBySize(data []model.File, sortOrder string, userID string) []model.File
	SortByDateCreated(data []model.File, sortOrder string) []model.File
	SortByDateModified(data []model.File, sortOrder string) []model.File
	SortByKind(data []model.File, userID string) []model.File
}

func NewFileSortService() FileSortService {
	return newFileSortService()
}

type fileSortService struct {
	fileMapper    FileMapper
	fileFilterSvc FileFilterService
}

func newFileSortService() *fileSortService {
	return &fileSortService{
		fileMapper:    NewFileMapper(),
		fileFilterSvc: NewFileFilterService(),
	}
}

func (svc *fileSortService) Sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
	if sortBy == SortByName {
		return svc.SortByName(data, sortOrder)
	} else if sortBy == SortBySize {
		return svc.SortBySize(data, sortOrder, userID)
	} else if sortBy == SortByDateCreated {
		return svc.SortByDateCreated(data, sortOrder)
	} else if sortBy == SortByDateModified {
		return svc.SortByDateModified(data, sortOrder)
	} else if sortBy == SortByKind {
		return svc.SortByKind(data, userID)
	}
	return data
}

func (svc *fileSortService) SortByName(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if sortOrder == SortOrderDesc {
			return data[i].GetName() > data[j].GetName()
		} else {
			return data[i].GetName() < data[j].GetName()
		}
	})
	return data
}

func (svc *fileSortService) SortBySize(data []model.File, sortOrder string, userID string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		fileA, err := svc.fileMapper.MapOne(data[i], userID)
		if err != nil {
			return false
		}
		fileB, err := svc.fileMapper.MapOne(data[j], userID)
		if err != nil {
			return false
		}
		var sizeA int64 = 0
		if fileA.Snapshot != nil && fileA.Snapshot.Original != nil {
			sizeA = *fileA.Snapshot.Original.Size
		}
		var sizeB int64 = 0
		if fileB.Snapshot != nil && fileB.Snapshot.Original != nil {
			sizeB = *fileB.Snapshot.Original.Size
		}
		if sortOrder == SortOrderDesc {
			return sizeA > sizeB
		} else {
			return sizeA < sizeB
		}
	})
	return data
}

func (svc *fileSortService) SortByDateCreated(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		a, _ := time.Parse(time.RFC3339, data[i].GetCreateTime())
		b, _ := time.Parse(time.RFC3339, data[j].GetCreateTime())
		if sortOrder == SortOrderDesc {
			return a.UnixMilli() > b.UnixMilli()
		} else {
			return a.UnixMilli() < b.UnixMilli()
		}
	})
	return data
}

func (svc *fileSortService) SortByDateModified(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
			a, _ := time.Parse(time.RFC3339, *data[i].GetUpdateTime())
			b, _ := time.Parse(time.RFC3339, *data[j].GetUpdateTime())
			if sortOrder == SortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		} else {
			return false
		}
	})
	return data
}

func (svc *fileSortService) SortByKind(data []model.File, userID string) []model.File {
	var res []model.File
	folders := svc.fileFilterSvc.FilterFolders(data)
	files := svc.fileFilterSvc.FilterFiles(data)
	res = append(res, folders...)
	res = append(res, files...)
	res = append(res, svc.fileFilterSvc.FilterImages(files, userID)...)
	res = append(res, svc.fileFilterSvc.FilterPDFs(files, userID)...)
	res = append(res, svc.fileFilterSvc.FilterDocuments(files, userID)...)
	res = append(res, svc.fileFilterSvc.FilterVideos(files, userID)...)
	res = append(res, svc.fileFilterSvc.FilterTexts(files, userID)...)
	res = append(res, svc.fileFilterSvc.FilterOthers(files, userID)...)
	return res
}
