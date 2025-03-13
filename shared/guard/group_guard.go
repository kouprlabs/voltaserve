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
	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/logger"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type GroupGuard struct {
	groupCache *cache.GroupCache
}

func NewGroupGuard(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *GroupGuard {
	return &GroupGuard{
		groupCache: cache.NewGroupCache(postgres, redis, environment),
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
			logger.GetLogger().Error(err)
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
