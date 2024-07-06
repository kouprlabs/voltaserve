// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

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
