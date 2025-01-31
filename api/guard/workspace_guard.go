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

type WorkspaceGuard interface {
	IsAuthorized(userID string, workspace model.Workspace, permission string) bool
	Authorize(userID string, workspace model.Workspace, permission string) error
}

func NewWorkspaceGuard() WorkspaceGuard {
	return newWorkspaceGuard()
}

type workspaceGuard struct {
	groupCache cache.GroupCache
}

func newWorkspaceGuard() *workspaceGuard {
	return &workspaceGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *workspaceGuard) IsAuthorized(userID string, workspace model.Workspace, permission string) bool {
	for _, p := range workspace.GetUserPermissions() {
		if p.GetUserID() == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range workspace.GetGroupPermissions() {
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

func (g *workspaceGuard) Authorize(userID string, workspace model.Workspace, permission string) error {
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
