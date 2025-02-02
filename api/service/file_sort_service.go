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
	sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File
	sortByName(data []model.File, sortOrder string) []model.File
	sortBySize(data []model.File, sortOrder string, userID string) []model.File
	sortByDateCreated(data []model.File, sortOrder string) []model.File
	sortByDateModified(data []model.File, sortOrder string) []model.File
	sortByKind(data []model.File, userID string) []model.File
}

type fileSortService struct {
	fileMapper    FileMapper
	fileFilterSvc FileFilterService
}

func newFileSortService() FileSortService {
	return &fileSortService{
		fileMapper:    newFileMapper(),
		fileFilterSvc: newFileFilterService(),
	}
}

func (svc *fileSortService) sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
	if sortBy == SortByName {
		return svc.sortByName(data, sortOrder)
	} else if sortBy == SortBySize {
		return svc.sortBySize(data, sortOrder, userID)
	} else if sortBy == SortByDateCreated {
		return svc.sortByDateCreated(data, sortOrder)
	} else if sortBy == SortByDateModified {
		return svc.sortByDateModified(data, sortOrder)
	} else if sortBy == SortByKind {
		return svc.sortByKind(data, userID)
	}
	return data
}

func (svc *fileSortService) sortByName(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if sortOrder == SortOrderDesc {
			return data[i].GetName() > data[j].GetName()
		} else {
			return data[i].GetName() < data[j].GetName()
		}
	})
	return data
}

func (svc *fileSortService) sortBySize(data []model.File, sortOrder string, userID string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		fileA, err := svc.fileMapper.mapOne(data[i], userID)
		if err != nil {
			return false
		}
		fileB, err := svc.fileMapper.mapOne(data[j], userID)
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

func (svc *fileSortService) sortByDateCreated(data []model.File, sortOrder string) []model.File {
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

func (svc *fileSortService) sortByDateModified(data []model.File, sortOrder string) []model.File {
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

func (svc *fileSortService) sortByKind(data []model.File, userID string) []model.File {
	var res []model.File
	folders := svc.fileFilterSvc.filterFolders(data)
	files := svc.fileFilterSvc.filterFiles(data)
	res = append(res, folders...)
	res = append(res, files...)
	res = append(res, svc.fileFilterSvc.filterImages(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterPDFs(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterDocuments(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterVideos(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterTexts(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterOthers(files, userID)...)
	return res
}
