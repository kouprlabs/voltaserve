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
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileCreateService struct {
	fileRepo    repo.FileRepo
	fileSearch  search.FileSearch
	fileCache   cache.FileCache
	fileGuard   guard.FileGuard
	fileMapper  *fileMapper
	fileCoreSvc *fileCoreService
}

func NewFileCreateService() *FileCreateService {
	return &FileCreateService{
		fileRepo:    repo.NewFileRepo(),
		fileSearch:  search.NewFileSearch(),
		fileCache:   cache.NewFileCache(),
		fileGuard:   guard.NewFileGuard(),
		fileMapper:  newFileMapper(),
		fileCoreSvc: newFileCoreService(),
	}
}

type FileCreateOptions struct {
	WorkspaceID string `json:"workspaceId" validate:"required"`
	Name        string `json:"name"        validate:"required,max=255"`
	Type        string `json:"type"        validate:"required,oneof=file folder"`
	ParentID    string `json:"parentId"    validate:"required"`
}

func (svc *FileCreateService) Create(opts FileCreateOptions, userID string) (*File, error) {
	path := helper.PathFromFilename(opts.Name)
	parentID := opts.ParentID
	if len(path) > 1 {
		newParentID, err := svc.createDirectoriesForPath(path, parentID, opts.WorkspaceID, userID)
		if err != nil {
			return nil, err
		}
		parentID = *newParentID
	}
	return svc.create(FileCreateOptions{
		WorkspaceID: opts.WorkspaceID,
		Name:        helper.FilenameFromPath(path),
		Type:        opts.Type,
		ParentID:    parentID,
	}, userID)
}

func (svc *FileCreateService) createDirectoriesForPath(path []string, parentID string, workspaceID string, userID string) (*string, error) {
	for _, component := range path[:len(path)-1] {
		existing, err := svc.fileCoreSvc.getChildWithName(parentID, component)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			parentID = existing.GetID()
		} else {
			folder, err := svc.create(FileCreateOptions{
				Name:        component,
				Type:        model.FileTypeFolder,
				ParentID:    parentID,
				WorkspaceID: workspaceID,
			}, userID)
			if err != nil {
				return nil, err
			}
			parentID = folder.ID
		}
	}
	return &parentID, nil
}

func (svc *FileCreateService) create(opts FileCreateOptions, userID string) (*File, error) {
	if len(opts.ParentID) > 0 {
		if err := svc.validateParent(opts.ParentID, userID); err != nil {
			return nil, err
		}
		existing, err := svc.fileCoreSvc.getChildWithName(opts.ParentID, opts.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			res, err := svc.fileMapper.MapOne(existing, userID)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	file, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        opts.Name,
		WorkspaceID: opts.WorkspaceID,
		ParentID:    opts.ParentID,
		Type:        opts.Type,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(file.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	if err = svc.fileSearch.Index([]model.File{file}); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.MapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileCreateService) validateParent(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFolder {
		return errorpkg.NewFileIsNotAFolderError(file)
	}
	return nil
}
