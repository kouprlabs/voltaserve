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

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

type FileServiceTestSuite struct {
	suite.Suite
	workspace *service.Workspace
	users     []model.User
}

func (s *FileServiceTestSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := test.CreateOrganization(s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.workspace, err = test.CreateWorkspace(org.ID, s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
	}
}

func TestFileServiceSuite(t *testing.T) {
	suite.Run(t, new(FileServiceTestSuite))
}

func (s *FileServiceTestSuite) TestCreate_File() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("file.txt", file.Name)
	s.Equal(model.FileTypeFile, file.Type)
}

func (s *FileServiceTestSuite) TestCreate_Folder() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("folder", folder.Name)
	s.Equal(model.FileTypeFolder, folder.Type)
}

func (s *FileServiceTestSuite) TestCreate_NonExistentParent() {
	_, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    "non-existent-parent",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_DuplicateName() {
	_, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileWithSimilarNameExistsError().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestCreate_DuplicateNameUsingPath() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "a/b/c/file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	files, err := service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("a/b/c/file.txt", service.NewFileService().GetPathStringWithoutWorkspace(files))
}

func (s *FileServiceTestSuite) TestFind() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	foundFiles, err := service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(foundFiles, 1)
	s.Equal(file.ID, foundFiles[0].ID)
}

func (s *FileServiceTestSuite) TestFind_NonExistentFile() {
	files, err := service.NewFileService().Find([]string{helper.NewID()}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(files)
}

func (s *FileServiceTestSuite) TestFindByPath() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewFileService().FindByPath(fmt.Sprintf("/%s/folder/file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(file.ID, found.ID)
}

func (s *FileServiceTestSuite) TestFindByPath_NonExistentPath() {
	_, err := service.NewFileService().FindByPath(fmt.Sprintf("/%s/non-existent-path", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindByPath_FileWithoutLeadingSlash() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().FindByPath(fmt.Sprintf("%s/folder/file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindByPath_FileWithoutTrailingSlash() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/folder/file.txt/", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath_Folder() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)
}

func (s *FileServiceTestSuite) TestListByPath_File() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder/file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)
}

func (s *FileServiceTestSuite) TestListByPath_ListFolderWithoutLeadingSlash() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().ListByPath(fmt.Sprintf("%s/folder", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath_ListFileWithoutTrailingSlash() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewFileService().ListByPath(fmt.Sprintf("/%s/folder/file.txt/", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindPath() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	path, err := service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(path, 3)
	s.Equal(s.workspace.RootID, path[0].ID)
	s.Equal(folder.ID, path[1].ID)
	s.Equal(file.ID, path[2].ID)
}

func (s *FileServiceTestSuite) TestList() {
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: s.workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    s.workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(s.workspace.RootID, service.FileListOptions{
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

func (s *FileServiceTestSuite) TestList_Paginate() {
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: s.workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    s.workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(s.workspace.RootID, service.FileListOptions{
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

	list, err = service.NewFileService().List(s.workspace.RootID, service.FileListOptions{
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
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: s.workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    s.workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(s.workspace.RootID, service.FileListOptions{
		Page:      1,
		Size:      3,
		SortBy:    service.FileSortByName,
		SortOrder: service.FileSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("file C", list.Data[0].Name)
	s.Equal("file B", list.Data[1].Name)
	s.Equal("file A", list.Data[2].Name)
}

func (s *FileServiceTestSuite) TestList_Query() {
	for _, name := range []string{"foo bar", "hello world", "lorem ipsum"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: s.workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    s.workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewFileService().List(s.workspace.RootID, service.FileListOptions{
		Query: &service.FileQuery{
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

func (s *FileServiceTestSuite) TestProbe() {
	for _, name := range []string{"file A", "file B", "file C"} {
		_, err := service.NewFileService().Create(service.FileCreateOptions{
			WorkspaceID: s.workspace.ID,
			Name:        name,
			Type:        model.FileTypeFile,
			ParentID:    s.workspace.RootID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	probe, err := service.NewFileService().Probe(s.workspace.RootID, service.FileListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *FileServiceTestSuite) TestComputeSize() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	size, err := service.NewFileService().ComputeSize(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(size)
	s.GreaterOrEqual(*size, int64(0))
}

func (s *FileServiceTestSuite) TestCount() {
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	count, err := service.NewFileService().Count(folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(count)
	s.Equal(int64(1), *count)
}

func (s *FileServiceTestSuite) TestCopy() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	clone, err := service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("file.txt", clone.Name)
	s.Equal(model.FileTypeFile, clone.Type)
	s.Equal(folder.ID, *clone.ParentID)
}

func (s *FileServiceTestSuite) TestDelete() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().Delete(file.ID, s.users[0].GetID())
	s.Require().NoError(err)

	files, err := service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(files)
}

func (s *FileServiceTestSuite) TestDownloadOriginalBuffer() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
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

func (s *FileServiceTestSuite) TestMove() {
	folderA, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err = service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(folderB.ID, *file.ParentID)
}

func (s *FileServiceTestSuite) TestPatchName() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	patched, err := service.NewFileService().PatchName(file.ID, "file (edit).txt", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(file.ID, patched.ID)
	s.Equal("file (edit).txt", patched.Name)
}

func (s *FileServiceTestSuite) TestGrantUserPermission() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeUserPermission() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().RevokeUserPermission([]string{file.ID}, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestGrantGroupPermission() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewFileService().RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestReprocess() {
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	file, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("fixtures", "files", "file.txt")),
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewTaskService().Patch(file.Snapshot.Task.ID, service.TaskPatchOptions{
		Fields: []string{service.TaskFieldStatus},
		Status: helper.ToPtr(model.TaskStatusError),
	})
	s.Require().NoError(err)

	reprocessResult, err := service.NewFileService().Reprocess(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(reprocessResult.Accepted, 1)
}
