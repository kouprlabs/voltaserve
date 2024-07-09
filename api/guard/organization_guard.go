// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package guard

import (
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
)

type OrganizationGuard struct {
	groupCache *cache.GroupCache
}

func NewOrganizationGuard() *OrganizationGuard {
	return &OrganizationGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *OrganizationGuard) IsAuthorized(userID string, org model.Organization, permission string) bool {
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
		for _, u := range g.GetUsers() {
			if u == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
				return true
			}
		}
	}
	return false
}

func (g *OrganizationGuard) Authorize(userID string, org model.Organization, permission string) error {
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
