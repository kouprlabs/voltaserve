package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/log"
	"voltaserve/model"
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
