package guard

import (
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"

	log "github.com/sirupsen/logrus"
)

type FileGuard struct {
	groupCache *cache.GroupCache
}

func NewFileGuard() *FileGuard {
	return &FileGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *FileGuard) IsAuthorized(user model.UserModel, file model.FileModel, permission string) bool {
	for _, p := range file.GetUserPermissions() {
		if p.GetUserId() == user.GetId() && model.IsEquivalentPermission(p.GetValue(), permission) {
			return true
		}
	}
	for _, p := range file.GetGroupPermissions() {
		g, err := g.groupCache.Get(p.GetGroupId())
		if err != nil {
			log.Error(err)
			return false
		}
		for _, u := range g.GetUsers() {
			if u == user.GetId() && model.IsEquivalentPermission(p.GetValue(), permission) {
				return true
			}
		}
	}
	return false
}

func (g *FileGuard) Authorize(user model.UserModel, file model.FileModel, permission string) error {
	if !g.IsAuthorized(user, file, permission) {
		err := errorpkg.NewFilePermissionError(user, file, permission)
		if g.IsAuthorized(user, file, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewOrganizationNotFoundError(err)
		}
	}
	return nil
}
