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
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type WorkspaceGuard struct {
	groupCache *cache.GroupCache
}

func NewWorkspaceGuard() *WorkspaceGuard {
	return &WorkspaceGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *WorkspaceGuard) IsAuthorized(userID string, workspace model.Workspace, permission string) bool {
	for _, p := range workspace.GetUserPermissions() {
		if p.GetUserID() == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range workspace.GetGroupPermissions() {
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

func (g *WorkspaceGuard) Authorize(userID string, workspace model.Workspace, permission string) error {
	if !g.IsAuthorized(userID, workspace, permission) {
		err := errorpkg.NewWorkspacePermissionError(userID, workspace, permission)
		if g.IsAuthorized(userID, workspace, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewWorkspaceNotFoundError(err)
		}
	}
	return nil
}
