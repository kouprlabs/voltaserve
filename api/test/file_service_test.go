// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type FileServiceTestSuite struct {
	suite.Suite
	workspace *service.Workspace
	users     []model.User
}

func (s *FileServiceTestSuite) SetupTest() {
	var err error
	s.users, err = test_helper.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.workspace, err = test_helper.CreateWorkspace(org.ID, s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
	}
}

func TestFileServiceSuite(t *testing.T) {
	suite.Run(t, new(FileServiceTestSuite))
}

func (s *FileServiceTestSuite) TestCreate() {
	// Test creating a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("test-file.txt", file.Name)
	s.Equal(model.FileTypeFile, file.Type)

	// Test creating a folder
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("test-folder", folder.Name)
	s.Equal(model.FileTypeFolder, folder.Type)

	// Test creating a file with an invalid parent ID
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "invalid-parent-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    "invalid-parent-id",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())

	// Test creating a file with a duplicate name
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileWithSimilarNameExistsError().Error(), err.Error())

	// Test creating a file with a duplicate name using a path name
	file, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "a/b/c/test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	pathFiles, err := service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	var pathComponents []string
	for _, path := range pathFiles {
		pathComponents = append(pathComponents, path.Name)
	}
	path := strings.Join(pathComponents[1:], "/")
	s.Equal("a/b/c/test-file.txt", path)
}

func (s *FileServiceTestSuite) TestFind() {
	// Create a file to find
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "find-me.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the file
	foundFiles, err := service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(foundFiles, 1)
	s.Equal(file.ID, foundFiles[0].ID)

	// Test finding a non-existent file
	foundFiles, err = service.NewFileService().Find([]string{"non-existent-id"}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(foundFiles)
}

func (s *FileServiceTestSuite) TestFindByPath() {
	// Create a folder and a file inside it
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the file by path
	foundFile, err := service.NewFileService().FindByPath(fmt.Sprintf("/%s/test-folder/test-file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(file.ID, foundFile.ID)

	// Test finding a non-existent path
	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/non-existent-path", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())

	// Test finding the file without a leading slash
	_, err = service.NewFileService().FindByPath(fmt.Sprintf("%s/test-folder/test-file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())

	// Test finding the file with a trailing slash
	_, err = service.NewFileService().FindByPath(fmt.Sprintf("/%s/test-folder/test-file.txt/", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath() {
	// Create a folder and a file inside it
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test listing files in the folder
	files, err := service.NewFileService().ListByPath(fmt.Sprintf("/%s/test-folder", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)

	// Test listing the file
	_, err = service.NewFileService().ListByPath(fmt.Sprintf("/%s/test-folder/test-file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)

	// Test listing files in the folder without a leading slash
	_, err = service.NewFileService().ListByPath(fmt.Sprintf("%s/test-folder", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())

	// Test listing the file with a trailing slash
	_, err = service.NewFileService().ListByPath(fmt.Sprintf("/%s/test-folder/test-file.txt/", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindPath() {
	// Create a folder and a file inside it
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the path of the file
	path, err := service.NewFileService().FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(path, 3)
	s.Equal(s.workspace.RootID, path[0].ID)
	s.Equal(folder.ID, path[1].ID)
	s.Equal(file.ID, path[2].ID)
}

func (s *FileServiceTestSuite) TestProbe() {
	// Create a folder and a file inside it
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probing the folder
	probe, err := service.NewFileService().Probe(folder.ID, service.FileListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *FileServiceTestSuite) TestList() {
	// Create a folder and a file inside it
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test listing files in the folder
	list, err := service.NewFileService().List(folder.ID, service.FileListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)
	s.Equal(file.ID, list.Data[0].ID)
}

func (s *FileServiceTestSuite) TestComputeSize() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test computing the size of the file
	size, err := service.NewFileService().ComputeSize(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(size)
	s.GreaterOrEqual(*size, int64(0))
}

func (s *FileServiceTestSuite) TestCount() {
	// Create a folder and a file inside it
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test counting items in the folder
	count, err := service.NewFileService().Count(folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(count)
	s.Equal(int64(1), *count)
}

func (s *FileServiceTestSuite) TestCopy() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Create a destination folder
	folder, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test copying the file
	copiedFile, err := service.NewFileService().Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("test-file.txt", copiedFile.Name)
	s.Equal(model.FileTypeFile, copiedFile.Type)
	s.Equal(folder.ID, *copiedFile.ParentID)
}

func (s *FileServiceTestSuite) TestDelete() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test deleting the file
	err = service.NewFileService().Delete(file.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Verify the file is deleted
	foundFiles, err := service.NewFileService().Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(foundFiles)
}

func (s *FileServiceTestSuite) TestDownloadOriginalBuffer() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Store the file
	path := filepath.Join("fixtures", "files", "file.txt")
	content, err := os.ReadFile(path) //nolint:gosec // Used for tests only
	s.Require().NoError(err)
	_, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{Path: &path}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test downloading the file
	buf := new(bytes.Buffer)
	_, err = service.NewFileService().DownloadOriginalBuffer(file.ID, "", buf, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(string(content), buf.String())
}

func (s *FileServiceTestSuite) TestMove() {
	// Create two folders
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

	// Create a file in folder1
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test moving the file to folder2
	movedFile, err := service.NewFileService().Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(folderB.ID, *movedFile.ParentID)
}

func (s *FileServiceTestSuite) TestPatchName() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test patching the file name
	patchedFile, err := service.NewFileService().PatchName(file.ID, "new-name.txt", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("new-name.txt", patchedFile.Name)
}

func (s *FileServiceTestSuite) TestGrantUserPermission() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test granting user permission
	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeUserPermission() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Grant user permission first
	err = service.NewFileService().GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	// Test revoking user permission
	err = service.NewFileService().RevokeUserPermission([]string{file.ID}, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestGrantGroupPermission() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Create a group
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test granting group permission
	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Create a group
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.workspace.Organization.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Grant group permission first
	err = service.NewFileService().GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	// Test revoking group permission
	err = service.NewFileService().RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestReprocess() {
	// Create a file
	file, err := service.NewFileService().Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Store the file
	file, err = service.NewFileService().Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("fixtures", "files", "file.txt")),
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Alter task ID status to allow reprocessing
	_, err = service.NewTaskService().Patch(file.Snapshot.Task.ID, service.TaskPatchOptions{
		Fields: []string{service.TaskFieldStatus},
		Status: helper.ToPtr(string(model.TaskStatusError)),
	})
	s.Require().NoError(err)

	// Test reprocessing the file
	reprocessResult, err := service.NewFileService().Reprocess(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(reprocessResult.Accepted, 1)
}
