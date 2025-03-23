// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package mapper

import (
	"errors"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
)

type FileMapper struct {
	groupCache      *cache.GroupCache
	workspaceCache  *cache.WorkspaceCache
	workspaceMapper *WorkspaceMapper
	snapshotMapper  *SnapshotMapper
	snapshotCache   *cache.SnapshotCache
	snapshotRepo    *repo.SnapshotRepo
}

func NewFileMapper(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *FileMapper {
	return &FileMapper{
		groupCache:      cache.NewGroupCache(postgres, redis, environment),
		workspaceCache:  cache.NewWorkspaceCache(postgres, redis, environment),
		workspaceMapper: NewWorkspaceMapper(postgres, redis, environment),
		snapshotMapper:  NewSnapshotMapper(postgres, redis, environment),
		snapshotCache:   cache.NewSnapshotCache(postgres, redis, environment),
		snapshotRepo:    repo.NewSnapshotRepo(postgres, environment),
	}
}

func (mp *FileMapper) Map(m model.File, userID string) (*dto.File, error) {
	workspace, err := mp.findWorkspace(m.GetWorkspaceID(), userID)
	if err != nil {
		return nil, err
	}
	return mp.mapWithWorkspace(m, workspace, userID)
}

func (mp *FileMapper) MapMany(data []model.File, workspaceID string, userID string) ([]*dto.File, error) {
	res := make([]*dto.File, 0)
	workspace, err := mp.findWorkspace(workspaceID, userID)
	if err != nil {
		return nil, err
	}
	for _, file := range data {
		f, err := mp.mapWithWorkspace(file, workspace, userID)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, f)
	}
	return res, nil
}

func (mp *FileMapper) mapWithWorkspace(m model.File, workspace *dto.Workspace, userID string) (*dto.File, error) {
	res := &dto.File{
		ID:         m.GetID(),
		Workspace:  *workspace,
		Name:       m.GetName(),
		Type:       m.GetType(),
		ParentID:   m.GetParentID(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.GetSnapshotID() != nil {
		snapshot, err := mp.snapshotCache.Get(*m.GetSnapshotID())
		if err != nil {
			return nil, err
		}
		res.Snapshot = mp.snapshotMapper.Map(snapshot)
		res.Snapshot.IsActive = true
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
	shareCount := 0
	for _, p := range m.GetUserPermissions() {
		if p.GetUserID() != userID {
			shareCount++
		}
	}
	if res.Permission == model.PermissionOwner {
		shareCount += len(m.GetGroupPermissions())
		res.IsShared = new(bool)
		if shareCount > 0 {
			*res.IsShared = true
		} else {
			*res.IsShared = false
		}
	}
	return res, nil
}

func (mp *FileMapper) findWorkspace(workspaceID string, userID string) (*dto.Workspace, error) {
	workspace, err := mp.workspaceCache.Get(workspaceID)
	if err != nil {
		return nil, err
	}
	res, err := mp.workspaceMapper.Map(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
