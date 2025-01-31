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
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/reactivex/rxgo/v2"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileListService struct {
	fileCache      cache.FileCache
	fileRepo       repo.FileRepo
	fileSearch     search.FileSearch
	fileGuard      guard.FileGuard
	fileCoreSvc    *fileCoreService
	fileMapper     *fileMapper
	fileIdent      *infra.FileIdentifier
	workspaceRepo  repo.WorkspaceRepo
	workspaceGuard guard.WorkspaceGuard
}

func NewFileListService() *FileListService {
	return &FileListService{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileSearch:     search.NewFileSearch(),
		fileGuard:      guard.NewFileGuard(),
		fileCoreSvc:    newFileCoreService(),
		fileMapper:     newFileMapper(),
		fileIdent:      infra.NewFileIdentifier(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceGuard: guard.NewWorkspaceGuard(),
	}
}

type FileQuery struct {
	Text             *string `json:"text"                       validate:"required"`
	Type             *string `json:"type,omitempty"             validate:"omitempty,oneof=file folder"`
	CreateTimeAfter  *int64  `json:"createTimeAfter,omitempty"`
	CreateTimeBefore *int64  `json:"createTimeBefore,omitempty"`
	UpdateTimeAfter  *int64  `json:"updateTimeAfter,omitempty"`
	UpdateTimeBefore *int64  `json:"updateTimeBefore,omitempty"`
}

type FileList struct {
	Data          []*File    `json:"data"`
	TotalPages    uint64     `json:"totalPages"`
	TotalElements uint64     `json:"totalElements"`
	Page          uint64     `json:"page"`
	Size          uint64     `json:"size"`
	Query         *FileQuery `json:"query,omitempty"`
}

type FileListOptions struct {
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
	Query     *FileQuery
}

type FileProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

func (svc *FileListService) Probe(id string, opts FileListOptions, userID string) (*FileProbe, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFolder {
		return nil, errorpkg.NewFileIsNotAFolderError(file)
	}
	totalElements, err := svc.fileRepo.CountChildren(id)
	if err != nil {
		return nil, err
	}
	return &FileProbe{
		TotalElements: uint64(totalElements),                               // #nosec G115 integer overflow conversion
		TotalPages:    (uint64(totalElements) + opts.Size - 1) / opts.Size, // #nosec G115 integer overflow conversion
	}, nil
}

func (svc *FileListService) List(id string, opts FileListOptions, userID string) (*FileList, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFolder {
		return nil, errorpkg.NewFileIsNotAFolderError(file)
	}
	workspace, err := svc.workspaceRepo.Find(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	var data []model.File
	if opts.Query != nil && opts.Query.Text != nil {
		data, err = svc.search(opts.Query, workspace)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = svc.getChildren(id)
		if err != nil {
			return nil, err
		}
	}
	return svc.list(data, file, opts, userID)
}

func (svc *FileListService) search(query *FileQuery, workspace model.Workspace) ([]model.File, error) {
	var res []model.File
	filter := fmt.Sprintf("workspaceId=\"%s\"", workspace.GetID())
	if query.Type != nil {
		filter += fmt.Sprintf(" AND type=\"%s\"", *query.Type)
	}
	hits, err := svc.fileSearch.Query(*query.Text, infra.QueryOptions{Filter: filter})
	if err != nil {
		return nil, err
	}
	for _, hit := range hits {
		var file model.File
		file, err := svc.fileCache.Get(hit.GetID())
		if err != nil {
			var e *errorpkg.ErrorResponse
			// We don't want to break if the search engine contains files that shouldn't be there
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, file)
	}
	return res, nil
}

func (svc *FileListService) getChildren(id string) ([]model.File, error) {
	var res []model.File
	ids, err := svc.fileRepo.FindChildrenIDs(id)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		var file model.File
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		res = append(res, file)
	}
	return res, nil
}

func (svc *FileListService) list(data []model.File, parent model.File, opts FileListOptions, userID string) (*FileList, error) {
	var filtered []model.File
	var err error
	if opts.Query != nil {
		filtered, err = svc.filterWithQuery(data, *opts.Query, parent)
		if err != nil {
			return nil, err
		}
	} else {
		filtered = data
	}
	authorized, err := svc.fileCoreSvc.authorize(filtered, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.sort(authorized, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mappedData, err := svc.fileMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	res := &FileList{
		Data:          mappedData,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          opts.Page,
		Size:          opts.Size,
		Query:         opts.Query,
	}
	return res, nil
}

func (svc *FileListService) sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
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

func (svc *FileListService) sortByName(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if sortOrder == SortOrderDesc {
			return data[i].GetName() > data[j].GetName()
		} else {
			return data[i].GetName() < data[j].GetName()
		}
	})
	return data
}

func (svc *FileListService) sortBySize(data []model.File, sortOrder string, userID string) []model.File {
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

func (svc *FileListService) sortByDateCreated(data []model.File, sortOrder string) []model.File {
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

func (svc *FileListService) sortByDateModified(data []model.File, sortOrder string) []model.File {
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

func (svc *FileListService) sortByKind(data []model.File, userID string) []model.File {
	var res []model.File
	res = append(res, svc.filterFolders(data)...)
	res = append(res, svc.filterFiles(data)...)
	res = append(res, svc.filterImages(data, userID)...)
	res = append(res, svc.filterPDFs(data, userID)...)
	res = append(res, svc.filterDocuments(data, userID)...)
	res = append(res, svc.filterVideos(data, userID)...)
	res = append(res, svc.filterTexts(data, userID)...)
	res = append(res, svc.filterOthers(data, userID)...)
	return res
}

func (svc *FileListService) getFromCache(data []interface{}) []model.File {
	var res []model.File
	for _, v := range data {
		var file model.File
		file, err := svc.fileCache.Get(v.(model.File).GetID())
		if err != nil {
			continue
		}
		res = append(res, file)
	}
	return res
}

func (svc *FileListService) filterFolders(data []model.File) []model.File {
	folders, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetType() == model.FileTypeFolder
		}).
		ToSlice(0)
	return svc.getFromCache(folders)
}

func (svc *FileListService) filterFiles(data []model.File) []model.File {
	files, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetType() == model.FileTypeFile
		}).
		ToSlice(0)
	return svc.getFromCache(files)
}

func (svc *FileListService) filterImages(data []model.File, userID string) []model.File {
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
	return svc.getFromCache(images)
}

func (svc *FileListService) filterPDFs(data []model.File, userID string) []model.File {
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
	return svc.getFromCache(pdfs)
}

func (svc *FileListService) filterDocuments(data []model.File, userID string) []model.File {
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
	return svc.getFromCache(documents)
}

func (svc *FileListService) filterVideos(data []model.File, userID string) []model.File {
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
	return svc.getFromCache(videos)
}

func (svc *FileListService) filterTexts(data []model.File, userID string) []model.File {
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
	return svc.getFromCache(texts)
}

func (svc *FileListService) filterOthers(data []model.File, userID string) []model.File {
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
	return svc.getFromCache(others)
}

func (svc *FileListService) paginate(data []model.File, page, size uint64) (pageData []model.File, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.File{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *FileListService) filterWithQuery(data []model.File, opts FileQuery, parent model.File) ([]model.File, error) {
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
		var file model.File
		file, err := svc.fileCache.Get(v.(model.File).GetID())
		if err != nil {
			return nil, err
		}
		res = append(res, file)
	}
	return res, nil
}
