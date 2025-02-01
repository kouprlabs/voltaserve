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
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FilePermissionService struct {
	fileCache      cache.FileCache
	fileRepo       repo.FileRepo
	fileGuard      guard.FileGuard
	fileCoreSvc    *fileCoreService
	userRepo       repo.UserRepo
	userMapper     *userMapper
	workspaceRepo  repo.WorkspaceRepo
	workspaceCache cache.WorkspaceCache
	groupCache     cache.GroupCache
	groupGuard     guard.GroupGuard
	groupMapper    *groupMapper
	permissionRepo repo.PermissionRepo
}

func NewFilePermissionService() *FilePermissionService {
	return &FilePermissionService{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileGuard:      guard.NewFileGuard(),
		fileCoreSvc:    newFileCoreService(),
		userRepo:       repo.NewUserRepo(),
		userMapper:     newUserMapper(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceCache: cache.NewWorkspaceCache(),
		groupCache:     cache.NewGroupCache(),
		groupGuard:     guard.NewGroupGuard(),
		groupMapper:    newGroupMapper(),
		permissionRepo: repo.NewPermissionRepo(),
	}
}

func (svc *FilePermissionService) GrantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		if _, err := svc.userRepo.Find(assigneeID); err != nil {
			return err
		}
		if err = svc.fileRepo.GrantUserPermission(id, assigneeID, permission); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		workspace, err := svc.workspaceRepo.Find(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		if err = svc.workspaceCache.Set(workspace); err != nil {
			return err
		}
		path, err := svc.fileRepo.FindPath(id)
		if err != nil {
			return err
		}
		for _, f := range path {
			if err := svc.fileCoreSvc.Sync(f); err != nil {
				return err
			}
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		for _, f := range tree {
			if err := svc.fileCoreSvc.Sync(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FilePermissionService) RevokeUserPermission(ids []string, assigneeID string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		if _, err := svc.userRepo.Find(assigneeID); err != nil {
			return err
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeUserPermission(tree, assigneeID); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		for _, f := range tree {
			if _, err := svc.fileCache.Refresh(f.GetID()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FilePermissionService) GrantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		group, err := svc.groupCache.Get(groupID)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return err
		}
		if err = svc.fileRepo.GrantGroupPermission(id, groupID, permission); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		workspace, err := svc.workspaceRepo.Find(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		if err = svc.workspaceCache.Set(workspace); err != nil {
			return err
		}
		path, err := svc.fileRepo.FindPath(id)
		if err != nil {
			return err
		}
		for _, f := range path {
			if err := svc.fileCoreSvc.Sync(f); err != nil {
				return err
			}
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		for _, f := range tree {
			if err := svc.fileCoreSvc.Sync(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FilePermissionService) RevokeGroupPermission(ids []string, groupID string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		group, err := svc.groupCache.Get(groupID)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return err
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeGroupPermission(tree, groupID); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		for _, f := range tree {
			if _, err := svc.fileCache.Refresh(f.GetID()); err != nil {
				return err
			}
		}
	}
	return nil
}

type UserPermission struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	Permission string `json:"permission"`
}

func (svc *FilePermissionService) FindUserPermissions(id string, userID string) ([]*UserPermission, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.FindUserPermissions(id)
	if err != nil {
		return nil, err
	}
	res := make([]*UserPermission, 0)
	for _, p := range permissions {
		if p.GetUserID() == userID {
			continue
		}
		u, err := svc.userRepo.Find(p.GetUserID())
		if err != nil {
			return nil, err
		}
		res = append(res, &UserPermission{
			ID:         p.GetID(),
			User:       svc.userMapper.mapOne(u),
			Permission: p.GetPermission(),
		})
	}
	return res, nil
}

type GroupPermission struct {
	ID         string `json:"id"`
	Group      *Group `json:"group"`
	Permission string `json:"permission"`
}

func (svc *FilePermissionService) FindGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.FindGroupPermissions(id)
	if err != nil {
		return nil, err
	}
	res := make([]*GroupPermission, 0)
	for _, p := range permissions {
		m, err := svc.groupCache.Get(p.GetGroupID())
		if err != nil {
			return nil, err
		}
		g, err := svc.groupMapper.mapOne(m, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, &GroupPermission{
			ID:         p.GetID(),
			Group:      g,
			Permission: p.GetPermission(),
		})
	}
	return res, nil
}
