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

type GroupGuard struct {
	groupCache *cache.GroupCache
}

func NewGroupGuard() *GroupGuard {
	return &GroupGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *GroupGuard) IsAuthorized(userID string, group model.Group, permission string) bool {
	for _, p := range group.GetUserPermissions() {
		if p.GetUserID() == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range group.GetGroupPermissions() {
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

func (g *GroupGuard) Authorize(userID string, group model.Group, permission string) error {
	if !g.IsAuthorized(userID, group, permission) {
		err := errorpkg.NewGroupPermissionError(userID, group, permission)
		if g.IsAuthorized(userID, group, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewGroupNotFoundError(err)
		}
	}
	return nil
}
