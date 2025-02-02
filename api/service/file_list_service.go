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
	fileCoreSvc    FileCoreService
	fileFilterSvc  FileFilterService
	fileSortSvc    FileSortService
	fileMapper     FileMapper
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
		fileFilterSvc:  newFileFilterService(),
		fileSortSvc:    newFileSortService(),
		fileMapper:     newFileMapper(),
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
	return svc.createList(data, file, opts, userID)
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

func (svc *FileListService) createList(data []model.File, parent model.File, opts FileListOptions, userID string) (*FileList, error) {
	var filtered []model.File
	var err error
	if opts.Query != nil {
		filtered, err = svc.fileFilterSvc.filterWithQuery(data, *opts.Query, parent)
		if err != nil {
			return nil, err
		}
	} else {
		filtered = data
	}
	authorized, err := svc.fileCoreSvc.authorize(userID, filtered, model.PermissionViewer)
	if err != nil {
		return nil, err
	}
	sorted := svc.fileSortSvc.sort(authorized, opts.SortBy, opts.SortOrder, userID)
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
