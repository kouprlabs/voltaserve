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
	"fmt"
	"strings"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileFetchService struct {
	fileCache      cache.FileCache
	fileRepo       repo.FileRepo
	fileSearch     search.FileSearch
	fileGuard      guard.FileGuard
	fileMapper     *fileMapper
	fileCoreSvc    *fileCoreService
	fileIdent      *infra.FileIdentifier
	userRepo       repo.UserRepo
	workspaceRepo  repo.WorkspaceRepo
	workspaceSvc   *WorkspaceService
	workspaceGuard guard.WorkspaceGuard
}

func NewFileFetchService() *FileFetchService {
	return &FileFetchService{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileSearch:     search.NewFileSearch(),
		fileGuard:      guard.NewFileGuard(),
		fileMapper:     newFileMapper(),
		fileCoreSvc:    newFileCoreService(),
		fileIdent:      infra.NewFileIdentifier(),
		userRepo:       repo.NewUserRepo(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceSvc:   NewWorkspaceService(),
		workspaceGuard: guard.NewWorkspaceGuard(),
	}
}

func (svc *FileFetchService) Find(ids []string, userID string) ([]*File, error) {
	var res []*File
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			continue
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		mapped, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, mapped)
	}
	return res, nil
}

func (svc *FileFetchService) FindByPath(path string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		return svc.getUserAsFile(user), nil
	}
	components, err := svc.getComponentsFromPath(path)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceSvc.Find(helper.WorkspaceIDFromSlug(components[0]), userID)
	if err != nil {
		return nil, err
	}
	if len(components) == 1 {
		return svc.getWorkspaceAsFile(workspace), nil
	}
	file, err := svc.getFileFromComponents(components, userID)
	if err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileFetchService) ListByPath(path string, userID string) ([]*File, error) {
	if path == "/" {
		return svc.getWorkspacesAsFiles(userID)
	}
	components, err := svc.getComponentsFromPath(path)
	if err != nil {
		return nil, err
	}
	file, err := svc.getFileFromComponents(components, userID)
	if err != nil {
		return nil, err
	}
	if file.GetType() == model.FileTypeFolder {
		children, err := svc.getAuthorizedChildren(file.GetID(), userID)
		if err != nil {
			return nil, err
		}
		res, err := svc.fileMapper.mapMany(children, userID)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if file.GetType() == model.FileTypeFile {
		res, err := svc.Find([]string{file.GetID()}, userID)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else {
		// This should never happen
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("invalid file type %s", file.GetType()))
	}
}

func (svc *FileFetchService) getWorkspacesAsFiles(userID string) ([]*File, error) {
	workspaces, err := svc.workspaceSvc.findAll(userID)
	if err != nil {
		return nil, err
	}
	res := make([]*File, 0)
	for _, w := range workspaces {
		res = append(res, svc.getWorkspaceAsFile(w))
	}
	return res, nil
}

func (svc *FileFetchService) getWorkspaceAsFile(workspace *Workspace) *File {
	return &File{
		ID:          workspace.RootID,
		WorkspaceID: workspace.ID,
		Name:        helper.SlugFromWorkspace(workspace.ID, workspace.Name),
		Type:        model.FileTypeFolder,
		Permission:  workspace.Permission,
		CreateTime:  workspace.CreateTime,
		UpdateTime:  workspace.UpdateTime,
	}
}

func (svc *FileFetchService) getUserAsFile(user model.User) *File {
	return &File{
		ID:          user.GetID(),
		WorkspaceID: "",
		Name:        "/",
		Type:        model.FileTypeFolder,
		Permission:  "owner",
		CreateTime:  user.GetCreateTime(),
		UpdateTime:  nil,
	}
}

func (svc *FileFetchService) getFileFromComponents(components []string, userID string) (model.File, error) {
	workspace, err := svc.workspaceRepo.Find(helper.WorkspaceIDFromSlug(components[0]))
	if err != nil {
		return nil, err
	}
	res, err := svc.fileCache.Get(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	for _, component := range components[1:] {
		file, err := svc.findComponentInFolder(component, res.GetID(), userID)
		if err != nil {
			return nil, err
		}
		res = file
		if file.GetType() == model.FileTypeFolder {
			continue
		} else if file.GetType() == model.FileTypeFile {
			break
		}
	}
	return res, err
}

func (svc *FileFetchService) findComponentInFolder(component string, id string, userID string) (model.File, error) {
	children, err := svc.getAuthorizedChildren(id, userID)
	if err != nil {
		return nil, err
	}
	var filtered []model.File
	for _, f := range children {
		if f.GetName() == component {
			filtered = append(filtered, f)
		}
	}
	if len(filtered) > 0 {
		return filtered[0], nil
	} else {
		return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("item not found '%s'", component))
	}
}

func (svc *FileFetchService) getAuthorizedChildren(id string, userID string) ([]model.File, error) {
	childrenIDs, err := svc.fileRepo.FindChildrenIDs(id)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.fileCoreSvc.authorizeIDs(childrenIDs, userID)
	if err != nil {
		return nil, err
	}
	return authorized, nil
}

func (svc *FileFetchService) getComponentsFromPath(path string) ([]string, error) {
	components := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	return components, nil
}

func (svc *FileFetchService) FindPath(id string, userID string) ([]*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	path, err := svc.fileRepo.FindPath(id)
	if err != nil {
		return nil, err
	}
	res := make([]*File, 0)
	for _, file := range path {
		f, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append([]*File{f}, res...)
	}
	return res, nil
}
