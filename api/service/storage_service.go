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
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type StorageService struct {
	workspaceRepo  *repo.WorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	workspaceGuard *guard.WorkspaceGuard
	fileRepo       *repo.FileRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	storageMapper  *storageMapper
}

func NewStorageService() *StorageService {
	return &StorageService{
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceCache: cache.NewWorkspaceCache(),
		workspaceGuard: guard.NewWorkspaceGuard(),
		fileRepo:       repo.NewFileRepo(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		storageMapper:  newStorageMapper(),
	}
}

func (svc *StorageService) GetAccountUsage(userID string) (*dto.StorageUsage, error) {
	ids, err := svc.workspaceRepo.FindIDs()
	if err != nil {
		return nil, err
	}
	workspaces := make([]model.Workspace, 0)
	for _, id := range ids {
		var workspace model.Workspace
		workspace, err = svc.workspaceCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.workspaceGuard.IsAuthorized(userID, workspace, model.PermissionOwner) {
			workspaces = append(workspaces, workspace)
		}
	}
	var maxBytes int64 = 0
	var b int64 = 0
	for _, w := range workspaces {
		root, err := svc.fileCache.Get(w.GetRootID())
		if err != nil {
			return nil, err
		}
		size, err := svc.fileRepo.ComputeSize(root.GetID())
		if err != nil {
			return nil, err
		}
		b = b + size
		maxBytes = maxBytes + w.GetStorageCapacity()
	}
	return svc.storageMapper.mapStorageUsage(b, maxBytes), nil
}

func (svc *StorageService) GetWorkspaceUsage(workspaceID string, userID string) (*dto.StorageUsage, error) {
	workspace, err := svc.workspaceCache.Get(workspaceID)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	root, err := svc.fileCache.Get(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, root, model.PermissionViewer); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.ComputeSize(root.GetID())
	if err != nil {
		return nil, err
	}
	return svc.storageMapper.mapStorageUsage(size, workspace.GetStorageCapacity()), nil
}

func (svc *StorageService) GetFileUsage(fileID string, userID string) (*dto.StorageUsage, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.ComputeSize(file.GetID())
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	return svc.storageMapper.mapStorageUsage(size, workspace.GetStorageCapacity()), nil
}

type storageMapper struct{}

func newStorageMapper() *storageMapper {
	return &storageMapper{}
}

func (mp *storageMapper) mapStorageUsage(byteCount int64, maxBytes int64) *dto.StorageUsage {
	res := dto.StorageUsage{
		Bytes:    byteCount,
		MaxBytes: maxBytes,
	}
	if maxBytes != 0 {
		res.Percentage = int(byteCount * 100 / maxBytes)
	}
	return &res
}
