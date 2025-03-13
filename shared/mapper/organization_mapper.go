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

type OrganizationMapper struct {
	groupCache *cache.GroupCache
}

func NewOrganizationMapper(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *OrganizationMapper {
	return &OrganizationMapper{
		groupCache: cache.NewGroupCache(postgres, redis, environment),
	}
}

func (mp *OrganizationMapper) MapOne(m model.Organization, userID string) (*dto.Organization, error) {
	res := &dto.Organization{
		ID:         m.GetID(),
		Name:       m.GetName(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
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

func (mp *OrganizationMapper) MapMany(orgs []model.Organization, userID string) ([]*dto.Organization, error) {
	res := make([]*dto.Organization, 0)
	for _, org := range orgs {
		o, err := mp.MapOne(org, userID)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewOrganizationNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, o)
	}
	return res, nil
}
