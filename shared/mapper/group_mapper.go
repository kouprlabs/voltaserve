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
)

type GroupMapper struct {
	orgCache   *cache.OrganizationCache
	orgMapper  *OrganizationMapper
	groupCache *cache.GroupCache
}

func NewGroupMapper(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *GroupMapper {
	return &GroupMapper{
		orgCache:   cache.NewOrganizationCache(postgres, redis, environment),
		orgMapper:  NewOrganizationMapper(postgres, redis, environment),
		groupCache: cache.NewGroupCache(postgres, redis, environment),
	}
}

func (mp *GroupMapper) Map(m model.Group, userID string) (*dto.Group, error) {
	org, err := mp.findOrganization(m.GetOrganizationID(), userID)
	if err != nil {
		return nil, err
	}
	res := &dto.Group{
		ID:           m.GetID(),
		Name:         m.GetName(),
		Organization: *org,
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
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

func (mp *GroupMapper) MapMany(groups []model.Group, userID string) ([]*dto.Group, error) {
	res := make([]*dto.Group, 0)
	for _, group := range groups {
		g, err := mp.Map(group, userID)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewGroupNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, g)
	}
	return res, nil
}

func (mp *GroupMapper) findOrganization(orgID string, userID string) (*dto.Organization, error) {
	org, err := mp.orgCache.Get(orgID)
	if err != nil {
		return nil, err
	}
	return mp.orgMapper.Map(org, userID)
}
