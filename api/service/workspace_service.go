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

	"github.com/google/uuid"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/guard"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/mapper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
	"github.com/kouprlabs/voltaserve/shared/search"

	"github.com/kouprlabs/voltaserve/api/config"
)

type WorkspaceService struct {
	workspaceRepo          *repo.WorkspaceRepo
	workspaceCache         *cache.WorkspaceCache
	workspaceGuard         *guard.WorkspaceGuard
	workspaceSearch        *search.WorkspaceSearch
	workspaceMapper        *mapper.WorkspaceMapper
	workspaceWebhookClient *client.WorkspaceWebhookClient
	orgCache               *cache.OrganizationCache
	orgGuard               *guard.OrganizationGuard
	fileRepo               *repo.FileRepo
	fileCache              *cache.FileCache
	fileGuard              *guard.FileGuard
	fileMapper             *mapper.FileMapper
	fileDelete             *fileDelete
	s3                     infra.S3Manager
	config                 *config.Config
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceSearch: search.NewWorkspaceSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		workspaceGuard: guard.NewWorkspaceGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceMapper: mapper.NewWorkspaceMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceWebhookClient: client.NewWorkspaceWebhookClient(
			config.GetConfig().Security,
		),
		orgCache: cache.NewOrganizationCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgGuard: guard.NewOrganizationGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileDelete: newFileDelete(),
		s3:         infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		config:     config.GetConfig(),
	}
}

func (svc *WorkspaceService) Create(opts dto.WorkspaceCreateOptions, userID string) (*dto.Workspace, error) {
	if opts.StorageCapacity == 0 {
		opts.StorageCapacity = helper.MegabyteToByte(svc.config.Defaults.WorkspaceStorageCapacityMB)
	}
	org, err := svc.orgCache.Get(opts.OrganizationID)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionEditor); err != nil {
		return nil, err
	}
	if svc.config.WorkspaceWebhook != "" {
		if err := svc.workspaceWebhookClient.Call(config.GetConfig().WorkspaceWebhook, dto.WorkspaceWebhookOptions{
			EventType: dto.WorkspaceWebhookEventTypeCreate,
			UserID:    userID,
			Create:    &opts,
		}); err != nil {
			return nil, err
		}
	}
	workspace, err := svc.create(opts, userID)
	if err != nil {
		return nil, err
	}
	root, err := svc.createRoot(workspace, userID)
	if err != nil {
		return nil, err
	}
	workspace, err = svc.associateWithRoot(workspace, root)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Index([]model.Workspace{workspace}); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.Map(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Find(id string, userID string) (*dto.Workspace, error) {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.Map(workspace, userID)
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

func (svc *WorkspaceService) List(opts WorkspaceListOptions, userID string) (*dto.WorkspaceList, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.WorkspaceSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.WorkspaceSortOrderAsc
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.workspaceMapper.MapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &dto.WorkspaceList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *WorkspaceService) Probe(opts WorkspaceListOptions, userID string) (*dto.WorkspaceProbe, error) {
	all, err := svc.load(userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &dto.WorkspaceProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *WorkspaceService) PatchName(id string, name string, userID string) (*dto.Workspace, error) {
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
	res, err := svc.workspaceMapper.Map(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) PatchStorageCapacity(id string, storageCapacity int64, userID string) (*dto.Workspace, error) {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionOwner); err != nil {
		return nil, err
	}
	if svc.config.WorkspaceWebhook != "" {
		if err := svc.workspaceWebhookClient.Call(config.GetConfig().WorkspaceWebhook, dto.WorkspaceWebhookOptions{
			EventType:   dto.WorkspaceWebhookEventTypePatchStorageCapacity,
			UserID:      userID,
			WorkspaceID: &id,
			PatchStorageCapacity: &dto.WorkspacePatchStorageCapacityOptions{
				StorageCapacity: storageCapacity,
			},
		}); err != nil {
			return nil, err
		}
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
	res, err := svc.workspaceMapper.Map(workspace, userID)
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
	return svc.delete(id)
}

func (svc *WorkspaceService) HasEnoughSpaceForByteSize(id string, byteSize int64, userID string) (bool, error) {
	workspace, err := svc.workspaceRepo.Find(id)
	if err != nil {
		return false, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return false, err
	}
	root, err := svc.fileRepo.Find(workspace.GetRootID())
	if err != nil {
		return false, err
	}
	usage, err := svc.fileRepo.ComputeSize(root.GetID())
	if err != nil {
		return false, err
	}
	expectedUsage := usage + byteSize
	if expectedUsage > workspace.GetStorageCapacity() {
		return false, nil
	}
	return true, nil
}

func (svc *WorkspaceService) GetBucket(id string) (string, error) {
	workspace, err := svc.workspaceRepo.Find(id)
	if err != nil {
		return "", err
	}
	return workspace.GetBucket(), nil
}

func (svc *WorkspaceService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.WorkspaceSortByName ||
		value == dto.WorkspaceSortByDateCreated ||
		value == dto.WorkspaceSortByDateModified
}

func (svc *WorkspaceService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.WorkspaceSortOrderAsc || value == dto.WorkspaceSortOrderDesc
}

func (svc *WorkspaceService) findAll(opts WorkspaceListOptions, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	var err error
	if opts.Query == "" {
		res, err = svc.load(userID)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = svc.search(opts, userID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *WorkspaceService) load(userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	ids, err := svc.workspaceRepo.FindIDs()
	if err != nil {
		return nil, err
	}
	res, err = svc.authorizeIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) search(opts WorkspaceListOptions, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	count, err := svc.workspaceRepo.Count()
	if err != nil {
		return nil, err
	}
	hits, err := svc.workspaceSearch.Query(opts.Query, infra.SearchQueryOptions{Limit: count})
	if err != nil {
		return nil, err
	}
	var workspaces []model.Workspace
	for _, hit := range hits {
		workspace, err := svc.workspaceCache.Get(hit.GetID())
		if err != nil {
			var e *errorpkg.ErrorResponse
			// We don't want to break if the search engine contains workspaces that shouldn't be there
			if errors.As(err, &e) && e.Code == errorpkg.NewWorkspaceNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		workspaces = append(workspaces, workspace)
	}
	res, err = svc.authorize(workspaces, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) create(opts dto.WorkspaceCreateOptions, userID string) (model.Workspace, error) {
	var bucket string
	if config.GetConfig().S3.Bucket == "" {
		bucket = strings.ReplaceAll(uuid.NewString(), "-", "")
		if err := svc.s3.CreateBucket(bucket); err != nil {
			return nil, err
		}
	} else {
		bucket = config.GetConfig().S3.Bucket
	}
	res, err := svc.workspaceRepo.Insert(repo.WorkspaceInsertOptions{
		ID:              helper.NewID(),
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		OrganizationID:  opts.OrganizationID,
		Image:           opts.Image,
		Bucket:          bucket,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceRepo.GrantUserPermission(res.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	res, err = svc.workspaceRepo.Find(res.GetID())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) createRoot(workspace model.Workspace, userID string) (model.File, error) {
	res, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        "root",
		WorkspaceID: workspace.GetID(),
		Type:        model.FileTypeFolder,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(res.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	if _, err := svc.fileCache.Refresh(res.GetID()); err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) associateWithRoot(workspace model.Workspace, root model.File) (model.Workspace, error) {
	if err := svc.workspaceRepo.UpdateRootID(workspace.GetID(), root.GetID()); err != nil {
		return nil, err
	}
	res, err := svc.workspaceCache.Refresh(workspace.GetID())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) authorize(data []model.Workspace, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	for _, w := range data {
		if svc.workspaceGuard.IsAuthorized(userID, w, model.PermissionViewer) {
			res = append(res, w)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) authorizeIDs(ids []string, userID string) ([]model.Workspace, error) {
	var res []model.Workspace
	for _, id := range ids {
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

func (svc *WorkspaceService) sort(data []model.Workspace, sortBy string, sortOrder string) []model.Workspace {
	if sortBy == dto.WorkspaceSortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.WorkspaceSortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == dto.WorkspaceSortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a := helper.StringToTime(data[i].GetCreateTime())
			b := helper.StringToTime(data[j].GetCreateTime())
			if sortOrder == dto.WorkspaceSortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == dto.WorkspaceSortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a := helper.StringToTime(*data[i].GetUpdateTime())
				b := helper.StringToTime(*data[j].GetUpdateTime())
				if sortOrder == dto.WorkspaceSortOrderDesc {
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

func (svc *WorkspaceService) paginate(data []model.Workspace, page, size uint64) ([]model.Workspace, uint64, uint64) {
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
	return data[startIndex:endIndex], totalElements, totalPages
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

func (svc *WorkspaceService) delete(id string) error {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.deleteFiles(id); err != nil {
		return err
	}
	if err := svc.workspaceRepo.Delete(id); err != nil {
		return err
	}
	if err := svc.workspaceCache.Delete(id); err != nil {
		return err
	}
	if err := svc.workspaceSearch.Delete([]string{workspace.GetID()}); err != nil {
		return err
	}
	if config.GetConfig().S3.Bucket == "" {
		if err := svc.s3.RemoveBucket(workspace.GetBucket()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *WorkspaceService) deleteFiles(id string) error {
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.workspaceRepo.ClearRootID(id); err != nil {
		return err
	} else {
		if err := svc.fileDelete.deleteFolder(workspace.GetRootID()); err != nil {
			return err
		}
	}
	return nil
}
