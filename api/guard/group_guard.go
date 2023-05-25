package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"

	log "github.com/sirupsen/logrus"
)

type GroupGuard struct {
	groupCache *cache.GroupCache
}

func NewGroupGuard() *GroupGuard {
	return &GroupGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *GroupGuard) IsAuthorized(user model.UserModel, group model.GroupModel, permission string) bool {
	for _, p := range group.GetUserPermissions() {
		if p.GetUserID() == user.GetID() && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range group.GetGroupPermissions() {
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

func (g *GroupGuard) Authorize(user model.UserModel, group model.GroupModel, permission string) error {
	if !g.IsAuthorized(user, group, permission) {
		err := errorpkg.NewGroupPermissionError(user, group, permission)
		if g.IsAuthorized(user, group, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewGroupNotFoundError(err)
		}
	}
	return nil
}
