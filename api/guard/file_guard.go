package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/log"
	"voltaserve/model"
)

type FileGuard struct {
	groupCache *cache.GroupCache
}

func NewFileGuard() *FileGuard {
	return &FileGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *FileGuard) IsAuthorized(userID string, file model.File, permission string) bool {
	for _, p := range file.GetUserPermissions() {
		if p.GetUserID() == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range file.GetGroupPermissions() {
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

func (g *FileGuard) Authorize(userID string, file model.File, permission string) error {
	if !g.IsAuthorized(userID, file, permission) {
		err := errorpkg.NewFilePermissionError(userID, file, permission)
		if g.IsAuthorized(userID, file, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewOrganizationNotFoundError(err)
		}
	}
	return nil
}
