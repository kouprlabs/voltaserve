package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"

	log "github.com/sirupsen/logrus"
)

type WorkspaceGuard struct {
	groupCache *cache.GroupCache
}

func NewWorkspaceGuard() *WorkspaceGuard {
	return &WorkspaceGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *WorkspaceGuard) IsAuthorized(user model.UserModel, workspace model.WorkspaceModel, permission string) bool {
	for _, p := range workspace.GetUserPermissions() {
		if p.GetUserID() == user.GetID() && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range workspace.GetGroupPermissions() {
		g, err := g.groupCache.Get(p.GetGroupId())
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

func (g *WorkspaceGuard) Authorize(user model.UserModel, workspace model.WorkspaceModel, permission string) error {
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
