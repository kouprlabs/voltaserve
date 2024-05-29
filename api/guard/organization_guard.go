package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"

	"github.com/gofiber/fiber/v2/log"
)

type OrganizationGuard struct {
	groupCache *cache.GroupCache
}

func NewOrganizationGuard() *OrganizationGuard {
	return &OrganizationGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *OrganizationGuard) IsAuthorized(userID string, org model.Organization, permission string) bool {
	for _, p := range org.GetUserPermissions() {
		if p.GetUserID() == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range org.GetGroupPermissions() {
		g, err := g.groupCache.Get(p.GetGroupID())
		if err != nil {
			log.Error(err)
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

func (g *OrganizationGuard) Authorize(userID string, org model.Organization, permission string) error {
	if !g.IsAuthorized(userID, org, permission) {
		err := errorpkg.NewOrganizationPermissionError(userID, org, permission)
		if g.IsAuthorized(userID, org, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewOrganizationNotFoundError(err)
		}
	}
	return nil
}
