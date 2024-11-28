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
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type WorkspaceService struct {
	workspaceRepo   repo.WorkspaceRepo
	workspaceCache  *cache.WorkspaceCache
	workspaceGuard  *guard.WorkspaceGuard
	workspaceSearch *search.WorkspaceSearch
	workspaceMapper *workspaceMapper
	fileRepo        repo.FileRepo
	fileCache       *cache.FileCache
	fileGuard       *guard.FileGuard
	fileMapper      *FileMapper
	s3              *infra.S3Manager
	config          *config.Config
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo:   repo.NewWorkspaceRepo(),
		workspaceCache:  cache.NewWorkspaceCache(),
		workspaceSearch: search.NewWorkspaceSearch(),
		workspaceGuard:  guard.NewWorkspaceGuard(),
		workspaceMapper: newWorkspaceMapper(),
		fileRepo:        repo.NewFileRepo(),
		fileCache:       cache.NewFileCache(),
		fileGuard:       guard.NewFileGuard(),
		fileMapper:      NewFileMapper(),
		s3:              infra.NewS3Manager(),
		config:          config.GetConfig(),
	}
}

type Workspace struct {
	ID              string       `json:"id"`
	Image           *string      `json:"image,omitempty"`
	Name            string       `json:"name"`
	RootID          string       `json:"rootId,omitempty"`
	StorageCapacity int64        `json:"storageCapacity"`
	Permission      string       `json:"permission"`
	Organization    Organization `json:"organization"`
	CreateTime      string       `json:"createTime"`
	UpdateTime      *string      `json:"updateTime,omitempty"`
}

type WorkspaceCreateOptions struct {
	Name            string  `json:"name"            validate:"required,max=255"`
	Image           *string `json:"image"`
	OrganizationID  string  `json:"organizationId"  validate:"required"`
	StorageCapacity int64   `json:"storageCapacity"`
}

func (svc *WorkspaceService) Create(opts WorkspaceCreateOptions, userID string) (*Workspace, error) {
	id := helper.NewID()
	bucket := strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := svc.s3.CreateBucket(bucket); err != nil {
		return nil, err
	}
	if opts.StorageCapacity == 0 {
		opts.StorageCapacity = helper.MegabyteToByte(svc.config.Defaults.WorkspaceStorageCapacityMB)
	}
	workspace, err := svc.workspaceRepo.Insert(repo.WorkspaceInsertOptions{
		ID:              id,
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		OrganizationID:  opts.OrganizationID,
		Image:           opts.Image,
		Bucket:          bucket,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceRepo.GrantUserPermission(workspace.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	workspace, err = svc.workspaceRepo.Find(workspace.GetID())
	if err != nil {
		return nil, err
	}
	root, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        "root",
		WorkspaceID: workspace.GetID(),
		Type:        model.FileTypeFolder,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(root.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	if _, err := svc.fileCache.Refresh(root.GetID()); err != nil {
		return nil, err
	}
	if err = svc.workspaceRepo.UpdateRootID(workspace.GetID(), root.GetID()); err != nil {
		return nil, err
	}
	workspace, err = svc.workspaceCache.Refresh(workspace.GetID())
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Index([]model.Workspace{workspace}); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Find(id string, userID string) (*Workspace, error) {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type WorkspaceListOptions struct {
	Query     string
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
}

type WorkspaceList struct {
	Data          []*Workspace `json:"data"`
	TotalPages    uint64       `json:"totalPages"`
	TotalElements uint64       `json:"totalElements"`
	Page          uint64       `json:"page"`
	Size          uint64       `json:"size"`
}

func (svc *WorkspaceService) List(opts WorkspaceListOptions, userID string) (*WorkspaceList, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.workspaceMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &WorkspaceList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

type WorkspaceProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

func (svc *WorkspaceService) Probe(opts WorkspaceListOptions, userID string) (*WorkspaceProbe, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &WorkspaceProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *WorkspaceService) PatchName(id string, name string, userID string) (*Workspace, error) {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionEditor); err != nil {
		return nil, err
	}
	if workspace, err = svc.workspaceRepo.UpdateName(id, name); err != nil {
		return nil, err
	}
	if err = svc.sync(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) PatchStorageCapacity(id string, storageCapacity int64, userID string) (*Workspace, error) {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionOwner); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.ComputeSize(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	if storageCapacity < size {
		return nil, errorpkg.NewInsufficientStorageCapacityError()
	}
	if workspace, err = svc.workspaceRepo.UpdateStorageCapacity(id, storageCapacity); err != nil {
		return nil, err
	}
	if err = svc.sync(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Delete(id string, userID string) error {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionOwner); err != nil {
		return err
	}
	if err = svc.workspaceRepo.Delete(id); err != nil {
		return err
	}
	if err = svc.workspaceSearch.Delete([]string{workspace.GetID()}); err != nil {
		return err
	}
	if err = svc.workspaceCache.Delete(id); err != nil {
		return err
	}
	if err = svc.s3.RemoveBucket(workspace.GetBucket()); err != nil {
		return err
	}
	return nil
}

func (svc *WorkspaceService) HasEnoughSpaceForByteSize(id string, byteSize int64) (*bool, error) {
	workspace, err := svc.workspaceRepo.Find(id)
	if err != nil {
		return nil, err
	}
	root, err := svc.fileRepo.Find(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	usage, err := svc.fileRepo.ComputeSize(root.GetID())
	if err != nil {
		return nil, err
	}
	expectedUsage := usage + byteSize
	if expectedUsage > workspace.GetStorageCapacity() {
		return helper.ToPtr(false), err
	}
	return helper.ToPtr(true), nil
}

func (svc *WorkspaceService) findAllWithoutOptions(userID string) ([]*Workspace, error) {
	ids, err := svc.workspaceRepo.FindIDs()
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorizationByIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	mapped, err := svc.workspaceMapper.mapMany(authorized, userID)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (svc *WorkspaceService) findAll(opts WorkspaceListOptions, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	if opts.Query == "" {
		ids, err := svc.workspaceRepo.FindIDs()
		if err != nil {
			return nil, err
		}
		res, err = svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
	} else {
		count, err := svc.workspaceRepo.Count()
		if err != nil {
			return nil, err
		}
		hits, err := svc.workspaceSearch.Query(opts.Query, infra.QueryOptions{Limit: count})
		if err != nil {
			return nil, err
		}
		var workspaces []model.Workspace
		for _, hit := range hits {
			workspace, err := svc.workspaceCache.Get(hit.GetID())
			if err != nil {
				continue
			}
			workspaces = append(workspaces, workspace)
		}
		res, err = svc.doAuthorization(workspaces, userID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *WorkspaceService) doAuthorization(data []model.Workspace, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	for _, w := range data {
		if svc.workspaceGuard.IsAuthorized(userID, w, model.PermissionViewer) {
			res = append(res, w)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) doAuthorizationByIDs(ids []string, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	for _, id := range ids {
		var w model.Workspace
		w, err := svc.workspaceCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewWorkspaceNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if svc.workspaceGuard.IsAuthorized(userID, w, model.PermissionViewer) {
			res = append(res, w)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) doSorting(data []model.Workspace, sortBy string, sortOrder string) []model.Workspace {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
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
	}
	return data
}

func (svc *WorkspaceService) doPagination(data []model.Workspace, page, size uint64) ([]model.Workspace, uint64, uint64) {
	totalElements := uint64(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Workspace{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

func (svc *WorkspaceService) sync(workspace model.Workspace) error {
	if err := svc.workspaceCache.Set(workspace); err != nil {
		return err
	}
	if err := svc.workspaceSearch.Update([]model.Workspace{workspace}); err != nil {
		return err
	}
	return nil
}

type workspaceMapper struct {
	orgCache   *cache.OrganizationCache
	orgMapper  *organizationMapper
	groupCache *cache.GroupCache
}

func newWorkspaceMapper() *workspaceMapper {
	return &workspaceMapper{
		orgCache:   cache.NewOrganizationCache(),
		orgMapper:  newOrganizationMapper(),
		groupCache: cache.NewGroupCache(),
	}
}

func (mp *workspaceMapper) mapOne(m model.Workspace, userID string) (*Workspace, error) {
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	o, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	res := &Workspace{
		ID:              m.GetID(),
		Name:            m.GetName(),
		RootID:          m.GetRootID(),
		StorageCapacity: m.GetStorageCapacity(),
		Organization:    *o,
		CreateTime:      m.GetCreateTime(),
		UpdateTime:      m.GetUpdateTime(),
	}
	res.Permission = model.PermissionNone
	for _, p := range m.GetUserPermissions() {
		if p.GetUserID() == userID && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
			res.Permission = p.GetValue()
		}
	}
	for _, p := range m.GetGroupPermissions() {
		g, err := mp.groupCache.Get(p.GetGroupID())
		if err != nil {
			return nil, err
		}
		for _, u := range g.GetMembers() {
			if u == userID && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
				res.Permission = p.GetValue()
			}
		}
	}
	return res, nil
}

func (mp *workspaceMapper) mapMany(workspaces []model.Workspace, userID string) ([]*Workspace, error) {
	res := make([]*Workspace, 0)
	for _, workspace := range workspaces {
		w, err := mp.mapOne(workspace, userID)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewWorkspaceNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, w)
	}
	return res, nil
}
