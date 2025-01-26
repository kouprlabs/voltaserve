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
	"strings"
	"time"

	"github.com/reactivex/rxgo/v2"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileFetch struct {
	fileCache      *cache.FileCache
	fileRepo       repo.FileRepo
	fileSearch     *search.FileSearch
	fileGuard      *guard.FileGuard
	fileMapper     *FileMapper
	fileIdent      *infra.FileIdentifier
	userRepo       repo.UserRepo
	workspaceRepo  repo.WorkspaceRepo
	workspaceSvc   *WorkspaceService
	workspaceGuard *guard.WorkspaceGuard
}

func NewFileFind() *FileFetch {
	return &FileFetch{
		fileCache:     cache.NewFileCache(),
		fileRepo:      repo.NewFileRepo(),
		fileSearch:    search.NewFileSearch(),
		fileGuard:     guard.NewFileGuard(),
		fileMapper:    NewFileMapper(),
		fileIdent:     infra.NewFileIdentifier(),
		userRepo:      repo.NewUserRepo(),
		workspaceRepo: repo.NewWorkspaceRepo(),
		workspaceSvc:  NewWorkspaceService(),
	}
}

func (svc *FileFetch) Find(ids []string, userID string) ([]*File, error) {
	var res []*File
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			continue
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		f, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (svc *FileFetch) FindByPath(path string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		return &File{
			ID:          user.GetID(),
			WorkspaceID: "",
			Name:        "/",
			Type:        model.FileTypeFolder,
			Permission:  "owner",
			CreateTime:  user.GetCreateTime(),
			UpdateTime:  nil,
		}, nil
	}
	components := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	workspace, err := svc.workspaceSvc.Find(helper.WorkspaceIDFromSlug(components[0]), userID)
	if err != nil {
		return nil, err
	}
	if len(components) == 1 {
		return &File{
			ID:          workspace.RootID,
			WorkspaceID: workspace.ID,
			Name:        helper.SlugFromWorkspace(workspace.ID, workspace.Name),
			Type:        model.FileTypeFolder,
			Permission:  workspace.Permission,
			CreateTime:  workspace.CreateTime,
			UpdateTime:  workspace.UpdateTime,
		}, nil
	}
	currentID := workspace.RootID
	components = components[1:]
	for _, component := range components {
		ids, err := svc.fileRepo.FindChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
		var filtered []model.File
		for _, f := range authorized {
			if f.GetName() == component {
				filtered = append(filtered, f)
			}
		}
		if len(filtered) > 0 {
			item := filtered[0]
			currentID = item.GetID()
			if item.GetType() == model.FileTypeFolder {
				continue
			} else if item.GetType() == model.FileTypeFile {
				break
			}
		} else {
			return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("component not found '%s'", component))
		}
	}
	result, err := svc.Find([]string{currentID}, userID)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("item not found '%s'", currentID))
	}
	return result[0], nil
}

func (svc *FileFetch) ListByPath(path string, userID string) ([]*File, error) {
	if path == "/" {
		workspaces, err := svc.workspaceSvc.findAllWithoutOptions(userID)
		if err != nil {
			return nil, err
		}
		result := make([]*File, 0)
		for _, w := range workspaces {
			result = append(result, &File{
				ID:          w.RootID,
				WorkspaceID: w.ID,
				Name:        helper.SlugFromWorkspace(w.ID, w.Name),
				Type:        model.FileTypeFolder,
				Permission:  w.Permission,
				CreateTime:  w.CreateTime,
				UpdateTime:  w.UpdateTime,
			})
		}
		return result, nil
	}
	components := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	workspace, err := svc.workspaceRepo.Find(helper.WorkspaceIDFromSlug(components[0]))
	if err != nil {
		return nil, err
	}
	currentID := workspace.GetRootID()
	currentType := model.FileTypeFolder
	components = components[1:]
	for _, component := range components {
		ids, err := svc.fileRepo.FindChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
		var filtered []model.File
		for _, f := range authorized {
			if f.GetName() == component {
				filtered = append(filtered, f)
			}
		}
		if len(filtered) > 0 {
			item := filtered[0]
			currentID = item.GetID()
			currentType = item.GetType()
			if item.GetType() == model.FileTypeFolder {
				continue
			} else if item.GetType() == model.FileTypeFile {
				break
			}
		} else {
			return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("component not found '%s'", component))
		}
	}
	if currentType == model.FileTypeFolder {
		ids, err := svc.fileRepo.FindChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
		result, err := svc.fileMapper.mapMany(authorized, userID)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else if currentType == model.FileTypeFile {
		result, err := svc.Find([]string{currentID}, userID)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("invalid file type %s", currentType))
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

func (svc *FileFetch) Probe(id string, opts FileListOptions, userID string) (*FileProbe, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
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

func (svc *FileFetch) List(id string, opts FileListOptions, userID string) (*FileList, error) {
	parent, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceRepo.Find(parent.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	var data []model.File
	if opts.Query != nil && opts.Query.Text != nil {
		filter := fmt.Sprintf("workspaceId=\"%s\"", workspace.GetID())
		if opts.Query.Type != nil {
			filter += fmt.Sprintf(" AND type=\"%s\"", *opts.Query.Type)
		}
		hits, err := svc.fileSearch.Query(*opts.Query.Text, infra.QueryOptions{Filter: filter})
		if err != nil {
			return nil, err
		}
		for _, hit := range hits {
			var f model.File
			f, err := svc.fileCache.Get(hit.GetID())
			if err != nil {
				var e *errorpkg.ErrorResponse
				// We don't want to break if the search engine contains files that shouldn't be there
				if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
					continue
				} else {
					return nil, err
				}
			}
			data = append(data, f)
		}
	} else {
		ids, err := svc.fileRepo.FindChildrenIDs(id)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			var f model.File
			f, err := svc.fileCache.Get(id)
			if err != nil {
				return nil, err
			}
			data = append(data, f)
		}
	}
	var filtered []model.File
	if opts.Query != nil {
		filtered, err = svc.doQueryFiltering(data, *opts.Query, parent)
	} else {
		filtered = data
	}
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorization(filtered, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(authorized, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
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

func (svc *FileFetch) FindPath(id string, userID string) ([]*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	path, err := svc.fileRepo.FindPath(id)
	if err != nil {
		return nil, err
	}
	res := make([]*File, 0)
	for _, file := range path {
		f, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append([]*File{f}, res...)
	}
	return res, nil
}

func (svc *FileFetch) doAuthorization(data []model.File, userID string) ([]model.File, error) {
	var res []model.File
	for _, f := range data {
		if svc.fileGuard.IsAuthorized(userID, f, model.PermissionViewer) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *FileFetch) doAuthorizationByIDs(ids []string, userID string) ([]model.File, error) {
	var res []model.File
	for _, id := range ids {
		var f model.File
		f, err := svc.fileCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if svc.fileGuard.IsAuthorized(userID, f, model.PermissionViewer) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *FileFetch) doSorting(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == SortBySize {
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
	} else if sortBy == SortByDateCreated {
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
	} else if sortBy == SortByDateModified {
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
	} else if sortBy == SortByKind {
		folders, _ := rxgo.Just(data)().
			Filter(func(v interface{}) bool {
				return v.(model.File).GetType() == model.FileTypeFolder
			}).
			ToSlice(0)
		files, _ := rxgo.Just(data)().
			Filter(func(v interface{}) bool {
				return v.(model.File).GetType() == model.FileTypeFile
			}).
			ToSlice(0)
		images, _ := rxgo.Just(files)().
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
		pdfs, _ := rxgo.Just(files)().
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
		documents, _ := rxgo.Just(files)().
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
		videos, _ := rxgo.Just(files)().
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
		texts, _ := rxgo.Just(files)().
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
		others, _ := rxgo.Just(files)().
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
		for _, v := range folders {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range images {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range pdfs {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range documents {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range videos {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range texts {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range others {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		return res
	}
	return data
}

func (svc *FileFetch) doPagination(data []model.File, page, size uint64) (pageData []model.File, totalElements uint64, totalPages uint64) {
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

func (svc *FileFetch) doQueryFiltering(data []model.File, opts FileQuery, parent model.File) ([]model.File, error) {
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
