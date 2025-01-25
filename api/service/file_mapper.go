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
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FileMapper struct {
	groupCache     *cache.GroupCache
	snapshotMapper *SnapshotMapper
	snapshotCache  *cache.SnapshotCache
	snapshotRepo   repo.SnapshotRepo
	config         *config.Config
}

func NewFileMapper() *FileMapper {
	return &FileMapper{
		groupCache:     cache.NewGroupCache(),
		snapshotMapper: NewSnapshotMapper(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		config:         config.GetConfig(),
	}
}

func (mp *FileMapper) mapOne(m model.File, userID string) (*File, error) {
	res := &File{
		ID:          m.GetID(),
		WorkspaceID: m.GetWorkspaceID(),
		Name:        m.GetName(),
		Type:        m.GetType(),
		ParentID:    m.GetParentID(),
		CreateTime:  m.GetCreateTime(),
		UpdateTime:  m.GetUpdateTime(),
	}
	if m.GetSnapshotID() != nil {
		snapshot, err := mp.snapshotCache.Get(*m.GetSnapshotID())
		if err != nil {
			return nil, err
		}
		res.Snapshot = mp.snapshotMapper.mapOne(snapshot)
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

func (mp *FileMapper) mapMany(data []model.File, userID string) ([]*File, error) {
	res := make([]*File, 0)
	for _, file := range data {
		f, err := mp.mapOne(file, userID)
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
