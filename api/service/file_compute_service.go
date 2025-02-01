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
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FileComputeService struct {
	fileCache cache.FileCache
	fileRepo  repo.FileRepo
	fileGuard guard.FileGuard
}

func NewFileComputeService() *FileComputeService {
	return &FileComputeService{
		fileCache: cache.NewFileCache(),
		fileRepo:  repo.NewFileRepo(),
		fileGuard: guard.NewFileGuard(),
	}
}

func (svc *FileComputeService) ComputeSize(id string, userID string) (*int64, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.fileRepo.ComputeSize(id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *FileComputeService) Count(id string, userID string) (*int64, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.fileRepo.CountItems(id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
