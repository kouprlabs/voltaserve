// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FilePatch struct {
	fileCache   *cache.FileCache
	fileRepo    repo.FileRepo
	fileGuard   *guard.FileGuard
	fileCoreSvc *FileCore
	fileMapper  *FileMapper
}

func NewFilePatch() *FilePatch {
	return &FilePatch{
		fileCache:   cache.NewFileCache(),
		fileRepo:    repo.NewFileRepo(),
		fileGuard:   guard.NewFileGuard(),
		fileCoreSvc: NewFileCore(),
		fileMapper:  NewFileMapper(),
	}
}

func (svc *FilePatch) PatchName(id string, name string, userID string) (*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetParentID() != nil {
		existing, err := svc.fileCoreSvc.GetChildWithName(*file.GetParentID(), name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errorpkg.NewFileWithSimilarNameExistsError()
		}
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	file.SetName(name)
	if err = svc.fileRepo.Save(file); err != nil {
		return nil, err
	}
	if err := svc.fileCoreSvc.Sync(file); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
