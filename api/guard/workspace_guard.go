package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"

	"github.com/gofiber/fiber/v2/log"
)

type WorkspaceGuard struct {
	groupCache *cache.GroupCache
}

func NewWorkspaceGuard() *WorkspaceGuard {
	return &WorkspaceGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *WorkspaceGuard) IsAuthorized(user model.User, workspace model.Workspace, permission string) bool {
	for _, p := range workspace.GetUserPermissions() {
		if p.GetUserID() == user.GetID() && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range workspace.GetGroupPermissions() {
		g, err := g.groupCache.Get(p.GetGroupID())
		if err != nil {
			log.Error(err)
			return false
		}
		for _, u := range g.GetUsers() {
			if u == user.GetID() && model.IsEquivalentPermission(p.GetValue(), permission) {
				return true
			}
		}
	}
	return false
}

func (g *WorkspaceGuard) Authorize(user model.User, workspace model.Workspace, permission string) error {
	if !g.IsAuthorized(user, workspace, permission) {
		err := errorpkg.NewWorkspacePermissionError(user, workspace, permission)
		if g.IsAuthorized(user, workspace, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewWorkspaceNotFoundError(err)
		}
	}
	return nil
}
