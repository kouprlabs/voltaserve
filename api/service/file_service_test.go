// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

type FileServiceTestSuite struct {
	suite.Suite
	users []model.User
}

func (s *FileServiceTestSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func TestFileServiceSuite(t *testing.T) {
	suite.Run(t, new(FileServiceTestSuite))
}

func (s *FileServiceTestSuite) TestCreate_File() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("file.txt", file.Name)
	s.Equal(model.FileTypeFile, file.Type)
}

func (s *FileServiceTestSuite) TestCreate_Folder() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("folder", folder.Name)
	s.Equal(model.FileTypeFolder, folder.Type)
}

func (s *FileServiceTestSuite) TestCreate_NonExistentParent() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    "non-existent-parent",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_NonExistentWorkspace() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: helper.NewID(),
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Require().Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_MissingWorkspacePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForWorkspace(workspace, s.users[0])

	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Require().Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_MissingParentPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	root, err := service.NewFileService().Find([]string{workspace.RootID}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(root[0], s.users[0])

	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Require().Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_InsufficientParentPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	err = repo.NewFileRepo(
		config.GetConfig().Postgres,
		config.GetConfig().Environment,
	).GrantUserPermission(workspace.RootID, s.users[0].GetID(), model.PermissionViewer)
	s.Require().NoError(err)
	root, err := cache.NewFileCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Refresh(workspace.RootID)
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Require().Equal(errorpkg.NewFilePermissionError(s.users[0].GetID(), root, model.PermissionEditor).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_DuplicateName() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileWithSimilarNameExistsError().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_DuplicateNameUsingPath() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "a/b/c/file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	files, err := service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("a/b/c/file.txt", service.NewFileService().GetPathStringWithoutWorkspace(files))
}

func (s *FileServiceTestSuite) TestFind() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)
}

func (s *FileServiceTestSuite) TestFind_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFind_NonExistentFile() {
	files, err := service.NewFileService().Find([]string{helper.NewID()}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(files)
}

func (s *FileServiceTestSuite) TestFindByPath() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewFileService().FindByPath(fmt.Sprintf("/%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(file.ID, found.ID)
}

func (s *FileServiceTestSuite) TestFindByPath_MissingFolderPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(folder, s.users[0])

	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindByPath_MissingLeafPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindByPath_NonExistentPath() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/non-existent-path", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindByPath_FileWithoutLeadingSlash() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().FindByPath(fmt.Sprintf("%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindByPath_FileWithoutTrailingSlash() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/folder/file.txt/", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath_Folder() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder", workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)
}

func (s *FileServiceTestSuite) TestListByPath_File() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)
}

func (s *FileServiceTestSuite) TestListByPath_ListFolderWithoutLeadingSlash() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().ListByPath(fmt.Sprintf("%s/folder", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath_ListFileWithoutTrailingSlash() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder/file.txt/", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath_MissingFolderPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(folder, s.users[0])

	_, err = service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath_MissingLeafPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder/file.txt", workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindPath() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	path, err := service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(path, 3)
	s.Equal(workspace.RootID, path[0].ID)
	s.Equal(folder.ID, path[1].ID)
	s.Equal(file.ID, path[2].ID)
}

func (s *FileServiceTestSuite) TestFindPath_MissingFolderPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(folder, s.users[0])

	_, err = service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindPath_MissingLeafPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestList() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("file A", list.Data[0].Name)
	s.Equal("file B", list.Data[1].Name)
	s.Equal("file C", list.Data[2].Name)
}

func (s *FileServiceTestSuite) TestList_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	var files []*dto.File
	for _, name := range []string{"file A", "file B", "file C"} {
		f, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		files = append(files, f)
	}

	s.revokeUserPermissionForFile(files[1], s.users[0])

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("file A", list.Data[0].Name)
	s.Equal("file C", list.Data[1].Name)
}

func (s *FileServiceTestSuite) TestList_Paginate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("file A", list.Data[0].Name)
	s.Equal("file B", list.Data[1].Name)

	list, err = service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("file C", list.Data[0].Name)
}

func (s *FileServiceTestSuite) TestList_SortByNameDescending() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Page:      1,
		Size:      3,
		SortBy:    dto.FileSortByName,
		SortOrder: dto.FileSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("file C", list.Data[0].Name)
	s.Equal("file B", list.Data[1].Name)
	s.Equal("file A", list.Data[2].Name)
}

func (s *FileServiceTestSuite) TestList_QueryText() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"foo bar", "hello world", "lorem ipsum"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Query: &dto.FileQuery{
			Text: helper.ToPtr("world"),
		},
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(1), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("hello world", list.Data[0].Name)
}

func (s *FileServiceTestSuite) TestList_QueryType() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	types := []string{model.FileTypeFile, model.FileTypeFile, model.FileTypeFolder}
	for i, name := range []string{"file A", "file B", "folder"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        types[i],
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Query: &dto.FileQuery{
			Type: helper.ToPtr(model.FileTypeFolder),
		},
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(1), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("folder", list.Data[0].Name)

	list, err = service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Query: &dto.FileQuery{
			Type: helper.ToPtr(model.FileTypeFile),
		},
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("file A", list.Data[0].Name)
	s.Equal("file B", list.Data[1].Name)
}

func (s *FileServiceTestSuite) TestList_QueryCreateTimeAfterBefore() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	var checkpoints []time.Time
	for _, name := range []string{"file A", "file B", "file C"} {
		checkpoints = append(checkpoints, time.Now())
		time.Sleep(1 * time.Second)
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}
	time.Sleep(1 * time.Second)
	checkpoints = append(checkpoints, time.Now())

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Query: &dto.FileQuery{
			CreateTimeAfter:  helper.ToPtr(helper.TimeToTimestamp(checkpoints[1])),
			CreateTimeBefore: helper.ToPtr(helper.TimeToTimestamp(checkpoints[3])),
		},
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("file B", list.Data[0].Name)
	s.Equal("file C", list.Data[1].Name)
}

func (s *FileServiceTestSuite) TestList_QueryUpdateTimeAfterBefore() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	var checkpoints []time.Time
	for _, name := range []string{"file A", "file B", "file C"} {
		checkpoints = append(checkpoints, time.Now())
		time.Sleep(1 * time.Second)
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}
	time.Sleep(1 * time.Second)
	checkpoints = append(checkpoints, time.Now())

	list, err := service.NewFileService().List(workspace.RootID, service.FileListOptions{
		Query: &dto.FileQuery{
			UpdateTimeAfter:  helper.ToPtr(helper.TimeToTimestamp(checkpoints[1])),
			UpdateTimeBefore: helper.ToPtr(helper.TimeToTimestamp(checkpoints[3])),
		},
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("file B", list.Data[0].Name)
	s.Equal("file C", list.Data[1].Name)
}

func (s *FileServiceTestSuite) TestProbe() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	probe, err := service.NewFileService().Probe(workspace.RootID, service.FileListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *FileServiceTestSuite) TestProbe_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	var files []*dto.File
	for _, name := range []string{"file A", "file B", "file C"} {
		f, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		files = append(files, f)
	}

	s.revokeUserPermissionForFile(files[1], s.users[0])

	probe, err := service.NewFileService().Probe(workspace.RootID, service.FileListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *FileServiceTestSuite) TestComputeSize() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	size, err := service.NewFileService().GetSize(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.GreaterOrEqual(*size, int64(0))
}

func (s *FileServiceTestSuite) TestComputeSize_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().GetSize(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCount() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	count, err := service.NewFileService().GetCount(folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(int64(1), *count)
}

func (s *FileServiceTestSuite) TestCount_NotAFolder() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().GetCount(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileIsNotAFolderError(cache.NewFileCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).GetOrNil(file.ID)).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCount_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(folder, s.users[0])

	_, err = service.NewFileService().GetCount(folder.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCopy() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	clone, err := service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("file.txt", clone.Name)
	s.Equal(model.FileTypeFile, clone.Type)
	s.Equal(folder.ID, *clone.ParentID)
}

func (s *FileServiceTestSuite) TestCopy_MissingSourcePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCopy_InsufficientSourcePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	_, err = service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestCopy_MissingTargetPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(folder, s.users[0])

	_, err = service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCopy_InsufficientTargetPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(folder, s.users[0], model.PermissionViewer)

	_, err = service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(folder.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestDelete() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().Delete(file.ID, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(files)
}

func (s *FileServiceTestSuite) TestDelete_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	err = service.NewFileService().Delete(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestDelete_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	err = service.NewFileService().Delete(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestDownloadOriginalBuffer() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	path := filepath.Join("fixtures", "files", "file.txt")
	content, err := os.ReadFile(path) //nolint:gosec // Used for tests only
	s.Require().NoError(err)
	_, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{Path: &path}, s.users[0].GetID())
	s.Require().NoError(err)

	buf := new(bytes.Buffer)
	_, err = service.NewFileService().DownloadOriginalBuffer(file.ID, "", buf, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(string(content), buf.String())
}

func (s *FileServiceTestSuite) TestDownloadOriginalBuffer_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("fixtures", "files", "file.txt")),
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().DownloadOriginalBuffer(file.ID, "", new(bytes.Buffer), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestMove() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folderA, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err = service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(folderB.ID, *file.ParentID)
}

func (s *FileServiceTestSuite) TestMove_MissingSourcePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folderA, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestMove_InsufficientSourcePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folderA, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	_, err = service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestMove_MissingTargetPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folderA, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(folderB, s.users[0])

	_, err = service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestMove_InsufficientTargetPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	folderA, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(folderB, s.users[0], model.PermissionViewer)

	_, err = service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(folderB.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestPatchName() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	patched, err := service.NewFileService().PatchName(file.ID, "file (edit).txt", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(file.ID, patched.ID)
	s.Equal("file (edit).txt", patched.Name)
}

func (s *FileServiceTestSuite) TestPatchName_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().PatchName(file.ID, "file (edit).txt", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestPatchName_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	_, err = service.NewFileService().PatchName(file.ID, "file (edit).txt", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestGrantUserPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestGrantUserPermission_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestGrantUserPermission_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestRevokeUserPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().RevokeUserPermission([]string{file.ID}, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeUserPermission_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	err = service.NewFileService().RevokeUserPermission([]string{file.ID}, s.users[1].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestRevokeUserPermission_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	err = service.NewFileService().RevokeUserPermission([]string{file.ID}, s.users[1].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestGrantGroupPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestGrantGroupPermission_MissingGroupPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestGrantGroupPermission_MissingFilePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestGrantGroupPermission_InsufficientFilePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission_MissingGroupPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	err = service.NewFileService().RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission_MissingFilePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	err = service.NewFileService().RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission_InsufficientFilePermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	err = service.NewFileService().RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) TestReprocess() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("fixtures", "files", "file.txt")),
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewTaskService().Patch(file.Snapshot.Task.ID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldStatus},
		Status: helper.ToPtr(model.TaskStatusError),
	})
	s.Require().NoError(err)

	reprocessResult, err := service.NewFileService().Reprocess(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(reprocessResult.Accepted, 1)
}

func (s *FileServiceTestSuite) TestReprocess_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("fixtures", "files", "file.txt")),
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewTaskService().Patch(file.Snapshot.Task.ID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldStatus},
		Status: helper.ToPtr(model.TaskStatusError),
	})
	s.Require().NoError(err)

	s.revokeUserPermissionForFile(file, s.users[0])

	_, err = service.NewFileService().Reprocess(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestReprocess_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := test.CreateWorkspace(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("fixtures", "files", "file.txt")),
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewTaskService().Patch(file.Snapshot.Task.ID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldStatus},
		Status: helper.ToPtr(model.TaskStatusError),
	})
	s.Require().NoError(err)

	s.grantUserPermissionForFile(file, s.users[0], model.PermissionViewer)

	_, err = service.NewFileService().Reprocess(file.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewFilePermissionError(
			s.users[0].GetID(),
			cache.NewFileCache(
				config.GetConfig().Postgres,
				config.GetConfig().Redis,
				config.GetConfig().Environment,
			).GetOrNil(file.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *FileServiceTestSuite) grantUserPermissionForFile(file *dto.File, user model.User, permission string) {
	err := repo.NewFileRepo(
		config.GetConfig().Postgres,
		config.GetConfig().Environment,
	).GrantUserPermission(file.ID, user.GetID(), permission)
	s.Require().NoError(err)
	_, err = cache.NewFileCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Refresh(file.ID)
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) revokeUserPermissionForFile(file *dto.File, user model.User) {
	err := repo.NewFileRepo(
		config.GetConfig().Postgres,
		config.GetConfig().Environment,
	).RevokeUserPermission(
		[]model.File{cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		).GetOrNil(file.ID)},
		user.GetID(),
	)
	s.Require().NoError(err)
	_, err = cache.NewFileCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Refresh(file.ID)
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) revokeUserPermissionForGroup(group *dto.Group, user model.User) {
	err := repo.NewGroupRepo(
		config.GetConfig().Postgres,
		config.GetConfig().Environment,
	).RevokeUserPermission(group.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewGroupCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Refresh(group.ID)
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) revokeUserPermissionForWorkspace(workspace *dto.Workspace, user model.User) {
	err := repo.NewWorkspaceRepo(
		config.GetConfig().Postgres,
		config.GetConfig().Environment,
	).RevokeUserPermission(workspace.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Refresh(workspace.ID)
	s.Require().NoError(err)
}
