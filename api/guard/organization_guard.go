package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"

	log "github.com/sirupsen/logrus"
)

type OrganizationGuard struct {
	groupCache *cache.GroupCache
}

func NewOrganizationGuard() *OrganizationGuard {
	return &OrganizationGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *OrganizationGuard) IsAuthorized(user model.UserModel, org model.OrganizationModel, permission string) bool {
	for _, p := range org.GetUserPermissions() {
		if p.GetUserID() == user.GetID() && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range org.GetGroupPermissions() {
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

func (g *OrganizationGuard) Authorize(user model.UserModel, org model.OrganizationModel, permission string) error {
	if !g.IsAuthorized(user, org, permission) {
		err := errorpkg.NewOrganizationPermissionError(user, org, permission)
		if g.IsAuthorized(user, org, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewOrganizationNotFoundError(err)
		}
	}
	return nil
}
