// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package guard

import (
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
)

type FileGuard interface {
	IsAuthorized(userID string, file model.File, permission string) bool
	Authorize(userID string, file model.File, permission string) error
}

func NewFileGuard() FileGuard {
	return newFileGuard()
}

type fileGuard struct {
	groupCache cache.GroupCache
}

func newFileGuard() *fileGuard {
	return &fileGuard{
		groupCache: cache.NewGroupCache(),
	}
}

func (g *fileGuard) IsAuthorized(userID string, file model.File, permission string) bool {
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
		for _, u := range g.GetMembers() {
			if u == userID && model.IsEquivalentPermission(p.GetValue(), permission) {
				return true
			}
		}
	}
	return false
}

func (g *fileGuard) Authorize(userID string, file model.File, permission string) error {
	if !g.IsAuthorized(userID, file, permission) {
		err := errorpkg.NewFilePermissionError(userID, file, permission)
		if g.IsAuthorized(userID, file, model.PermissionViewer) {
			return err
		} else {
			return errorpkg.NewFileNotFoundError(err)
		}
	}
	return nil
}
