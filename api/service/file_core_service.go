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
	"errors"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileCoreService interface {
	GetChildWithName(id string, name string) (model.File, error)
	Sync(file model.File) error
	SaveAndSync(file model.File) error
	Authorize(userID string, files []model.File, permission string) ([]model.File, error)
	AuthorizeIDs(userID string, ids []string, permission string) ([]model.File, error)
}

func NewFileCoreService() FileCoreService {
	return newFileCoreService()
}

type fileCoreService struct {
	fileRepo   repo.FileRepo
	fileSearch search.FileSearch
	fileCache  cache.FileCache
	fileGuard  guard.FileGuard
}

func newFileCoreService() *fileCoreService {
	return &fileCoreService{
		fileRepo:   repo.NewFileRepo(),
		fileCache:  cache.NewFileCache(),
		fileSearch: search.NewFileSearch(),
		fileGuard:  guard.NewFileGuard(),
	}
}

func (svc *fileCoreService) GetChildWithName(id string, name string) (model.File, error) {
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

func (svc *fileCoreService) Sync(file model.File) error {
	if err := svc.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	if err := svc.fileCache.Set(file); err != nil {
		return err
	}
	return nil
}

func (svc *fileCoreService) SaveAndSync(file model.File) error {
	if err := svc.fileRepo.Save(file); err != nil {
		return err
	}
	if err := svc.Sync(file); err != nil {
		return err
	}
	return nil
}

func (svc *fileCoreService) Authorize(userID string, files []model.File, permission string) ([]model.File, error) {
	var res []model.File
	for _, f := range files {
		if svc.fileGuard.IsAuthorized(userID, f, permission) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *fileCoreService) AuthorizeIDs(userID string, ids []string, permission string) ([]model.File, error) {
	var res []model.File
	for _, id := range ids {
		var f model.File
		f, err := svc.fileCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if svc.fileGuard.IsAuthorized(userID, f, permission) {
			res = append(res, f)
		}
	}
	return res, nil
}
