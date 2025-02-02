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
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/model"
)

type OrganizationMapper interface {
	mapOne(model.Organization, string) (*Organization, error)
	mapMany([]model.Organization, string) ([]*Organization, error)
}

type organizationMapper struct {
	groupCache cache.GroupCache
}

func newOrganizationMapper() OrganizationMapper {
	return &organizationMapper{
		groupCache: cache.NewGroupCache(),
	}
}

func (mp *organizationMapper) mapOne(m model.Organization, userID string) (*Organization, error) {
	res := &Organization{
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

func (mp *organizationMapper) mapMany(orgs []model.Organization, userID string) ([]*Organization, error) {
	res := make([]*Organization, 0)
	for _, org := range orgs {
		o, err := mp.mapOne(org, userID)
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
