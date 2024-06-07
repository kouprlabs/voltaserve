package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/log"
	"voltaserve/model"
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
		for _, u := range g.GetUsers() {
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
