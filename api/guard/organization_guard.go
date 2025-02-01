// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package guard

import (
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
)

type OrganizationGuard interface {
	IsAuthorized(userID string, org model.Organization, permission string) bool
	Authorize(userID string, org model.Organization, permission string) error
}

func NewOrganizationGuard() OrganizationGuard {
	return newOrganizationGuard()
}

type organizationGuard struct {
	groupCache cache.GroupCache
}

func newOrganizationGuard() *organizationGuard {
	return &organizationGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *organizationGuard) IsAuthorized(userID string, org model.Organization, permission string) bool {
	for _, p := range org.GetUserPermissions() {
		if p.GetUserID() == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range org.GetGroupPermissions() {
		g, err := g.groupCache.Get(p.GetGroupID())
		if err != nil {
			log.GetLogger().Error(err)
			return false
		}
		for _, u := range g.GetMembers() {
			if u == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
				return true
			}
		}
	}
	return false
}

func (g *organizationGuard) Authorize(userID string, org model.Organization, permission string) error {
	if !g.IsAuthorized(userID, org, permission) {
		err := errorpkg.NewOrganizationPermissionError(userID, org, permission)
		if g.IsAuthorized(userID, org, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewOrganizationNotFoundError(err)
		}
	}
	return nil
}
