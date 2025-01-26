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
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileCore struct {
	fileRepo   repo.FileRepo
	fileSearch *search.FileSearch
	fileCache  *cache.FileCache
}

func NewFileCore() *FileCore {
	return &FileCore{
		fileRepo:   repo.NewFileRepo(),
		fileCache:  cache.NewFileCache(),
		fileSearch: search.NewFileSearch(),
	}
}

func (svc *FileCore) GetChildWithName(id string, name string) (model.File, error) {
	children, err := svc.fileRepo.FindChildren(id)
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		if child.GetName() == name {
			return child, nil
		}
	}
	return nil, nil
}

func (svc *FileCore) Sync(file model.File) error {
	if err := svc.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	if err := svc.fileCache.Set(file); err != nil {
		return err
	}
	return nil
}

func (svc *FileCore) SaveAndSync(file model.File) error {
	if err := svc.fileRepo.Save(file); err != nil {
		return err
	}
	if err := svc.Sync(file); err != nil {
		return err
	}
	return nil
}
