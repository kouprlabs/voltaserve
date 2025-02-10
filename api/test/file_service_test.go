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

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
)

type FileServiceTestSuite struct {
	suite.Suite
	fileSvc   *service.FileService
	workspace *service.Workspace
	userIDs   []string
}

func (suite *FileServiceTestSuite) SetupTest() {
	userIDs, err := suite.createUsers()
	if err != nil {
		suite.Fail(err.Error())
	}
	org, err := suite.createOrganization(userIDs[0])
	if err != nil {
		suite.Fail(err.Error())
	}
	workspace, err := suite.createWorkspace(org.ID, userIDs[0])
	if err != nil {
		suite.Fail(err.Error())
	}
	suite.fileSvc = service.NewFileService()
	suite.workspace = workspace
	suite.userIDs = userIDs
}

func TestFileServiceSuite(t *testing.T) {
	suite.Run(t, new(FileServiceTestSuite))
}

func (suite *FileServiceTestSuite) TestCreate() {
	// Test creating a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal("test-file.txt", file.Name)
	suite.Equal(model.FileTypeFile, file.Type)

	// Test creating a folder
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal("test-folder", folder.Name)
	suite.Equal(model.FileTypeFolder, folder.Type)

	// Test creating a file with an invalid parent ID
	_, err = suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "invalid-parent-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    "invalid-parent-id",
	}, suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFileNotFoundError(err).Error())

	// Test creating a file with a duplicate name
	_, err = suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFileWithSimilarNameExistsError().Error())

	// Test creating a file with a duplicate name using a path name
	file, err = suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "a/b/c/test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	pathFiles, err := suite.fileSvc.FindPath(file.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	var pathComponents []string
	for _, path := range pathFiles {
		pathComponents = append(pathComponents, path.Name)
	}
	path := strings.Join(pathComponents[1:], "/")
	suite.Equal("a/b/c/test-file.txt", path)
}

func (suite *FileServiceTestSuite) TestFind() {
	// Create a file to find
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "find-me.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test finding the file
	foundFiles, err := suite.fileSvc.Find([]string{file.ID}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Len(foundFiles, 1)
	suite.Equal(file.ID, foundFiles[0].ID)

	// Test finding a non-existent file
	foundFiles, err = suite.fileSvc.Find([]string{"non-existent-id"}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Empty(foundFiles)
}

func (suite *FileServiceTestSuite) TestFindByPath() {
	// Create a folder and a file inside it
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test finding the file by path
	foundFile, err := suite.fileSvc.FindByPath(fmt.Sprintf("/%s/test-folder/test-file.txt", suite.workspace.ID), suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal(file.ID, foundFile.ID)

	// Test finding a non-existent path
	_, err = suite.fileSvc.FindByPath(fmt.Sprintf("/%s/non-existent-path", suite.workspace.ID), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFileNotFoundError(err).Error())

	// Test finding the file without a leading slash
	_, err = suite.fileSvc.FindByPath(fmt.Sprintf("%s/test-folder/test-file.txt/", suite.workspace.ID), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFilePathMissingLeadingSlash().Error())

	// Test finding the file with a trailing slash
	_, err = suite.fileSvc.FindByPath(fmt.Sprintf("/%s/test-folder/test-file.txt/", suite.workspace.ID), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error())
}

func (suite *FileServiceTestSuite) TestListByPath() {
	// Create a folder and a file inside it
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test listing files in the folder
	files, err := suite.fileSvc.ListByPath(fmt.Sprintf("/%s/test-folder", suite.workspace.ID), suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Len(files, 1)
	suite.Equal(file.ID, files[0].ID)

	// Test listing the file
	_, err = suite.fileSvc.ListByPath(fmt.Sprintf("/%s/test-folder/test-file.txt", suite.workspace.ID), suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Len(files, 1)
	suite.Equal(file.ID, files[0].ID)

	// Test listing files in the folder without a leading slash
	_, err = suite.fileSvc.ListByPath(fmt.Sprintf("%s/test-folder", suite.workspace.ID), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFilePathMissingLeadingSlash().Error())

	// Test listing the file with a trailing slash
	_, err = suite.fileSvc.ListByPath(fmt.Sprintf("/%s/test-folder/test-file.txt/", suite.workspace.ID), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error())
}

func (suite *FileServiceTestSuite) TestFindPath() {
	// Create a folder and a file inside it
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test finding the path of the file
	path, err := suite.fileSvc.FindPath(file.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Len(path, 3)
	suite.Equal(path[0].ID, suite.workspace.RootID)
	suite.Equal(path[1].ID, folder.ID)
	suite.Equal(path[2].ID, file.ID)
}

func (suite *FileServiceTestSuite) TestProbe() {
	// Create a folder and a file inside it
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	_, err = suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test probing the folder
	probe, err := suite.fileSvc.Probe(folder.ID, service.FileListOptions{Page: 1, Size: 10}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal(uint64(1), probe.TotalElements)
	suite.Equal(uint64(1), probe.TotalPages)
}

func (suite *FileServiceTestSuite) TestList() {
	// Create a folder and a file inside it
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	_, err = suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test listing files in the folder
	list, err := suite.fileSvc.List(folder.ID, service.FileListOptions{Page: 1, Size: 10}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Len(list.Data, 1)
	suite.Equal("test-file.txt", list.Data[0].Name)
}

func (suite *FileServiceTestSuite) TestComputeSize() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test computing the size of the file
	size, err := suite.fileSvc.ComputeSize(file.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(size)
	suite.GreaterOrEqual(*size, int64(0))
}

func (suite *FileServiceTestSuite) TestCount() {
	// Create a folder and a file inside it
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	_, err = suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test counting items in the folder
	count, err := suite.fileSvc.Count(folder.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(count)
	suite.Equal(int64(1), *count)
}

func (suite *FileServiceTestSuite) TestCopy() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Create a destination folder
	folder, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test copying the file
	copiedFile, err := suite.fileSvc.Copy(file.ID, folder.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal("test-file.txt", copiedFile.Name)
	suite.Equal(model.FileTypeFile, copiedFile.Type)
	suite.Equal(folder.ID, *copiedFile.ParentID)
}

func (suite *FileServiceTestSuite) TestDelete() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test deleting the file
	err = suite.fileSvc.Delete(file.ID, suite.userIDs[0])
	suite.Require().NoError(err)

	// Verify the file is deleted
	foundFiles, err := suite.fileSvc.Find([]string{file.ID}, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Empty(foundFiles)
}

func (suite *FileServiceTestSuite) TestDownloadOriginalBuffer() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Store the file
	path := filepath.Join("assets", "file.txt")
	content, err := os.ReadFile(path) //nolint:gosec // Used for tests only
	suite.Require().NoError(err)
	_, err = suite.fileSvc.Store(file.ID, service.FileStoreOptions{Path: &path}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test downloading the file
	buf := new(bytes.Buffer)
	_, err = suite.fileSvc.DownloadOriginalBuffer(file.ID, "", buf, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal(string(content), buf.String())
}

func (suite *FileServiceTestSuite) TestMove() {
	// Create two folders
	folder1, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "folder1",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)
	folder2, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "folder2",
		Type:        model.FileTypeFolder,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Create a file in folder1
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder1.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test moving the file to folder2
	movedFile, err := suite.fileSvc.Move(file.ID, folder2.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal(folder2.ID, *movedFile.ParentID)
}

func (suite *FileServiceTestSuite) TestPatchName() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test patching the file name
	patchedFile, err := suite.fileSvc.PatchName(file.ID, "new-name.txt", suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Equal("new-name.txt", patchedFile.Name)
}

func (suite *FileServiceTestSuite) TestGrantUserPermission() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test granting user permission
	err = suite.fileSvc.GrantUserPermission([]string{file.ID}, suite.userIDs[1], model.PermissionViewer, suite.userIDs[0])
	suite.Require().NoError(err)
}

func (suite *FileServiceTestSuite) TestRevokeUserPermission() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Grant user permission first
	err = suite.fileSvc.GrantUserPermission([]string{file.ID}, suite.userIDs[1], model.PermissionViewer, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test revoking user permission
	err = suite.fileSvc.RevokeUserPermission([]string{file.ID}, suite.userIDs[1], suite.userIDs[0])
	suite.Require().NoError(err)
}

func (suite *FileServiceTestSuite) TestGrantGroupPermission() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Create a group
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: suite.workspace.Organization.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test granting group permission
	err = suite.fileSvc.GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, suite.userIDs[0])
	suite.Require().NoError(err)
}

func (suite *FileServiceTestSuite) TestRevokeGroupPermission() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Create a group
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: suite.workspace.Organization.ID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Grant group permission first
	err = suite.fileSvc.GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test revoking group permission
	err = suite.fileSvc.RevokeGroupPermission([]string{file.ID}, group.ID, suite.userIDs[0])
	suite.Require().NoError(err)
}

func (suite *FileServiceTestSuite) TestReprocess() {
	// Create a file
	file, err := suite.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: suite.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    suite.workspace.RootID,
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Store the file
	_, err = suite.fileSvc.Store(file.ID, service.FileStoreOptions{
		Path: helper.ToPtr(filepath.Join("assets", "file.txt")),
	}, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test reprocessing the file
	reprocessResult, err := suite.fileSvc.Reprocess(file.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.Empty(reprocessResult.Accepted)
}

func (suite *FileServiceTestSuite) createUsers() ([]string, error) {
	db, err := infra.NewPostgresManager().GetDB()
	if err != nil {
		return nil, nil
	}
	var ids []string
	for i := range 2 {
		id := helper.NewID()
		db = db.Exec("INSERT INTO \"user\" (id, full_name, username, email, password_hash, create_time) VALUES (?, ?, ?, ?, ?, ?)",
			id, fmt.Sprintf("user %d", i), id+"@voltaserve.com", id+"@voltaserve.com", "", helper.NewTimestamp())
		if db.Error != nil {
			return nil, db.Error
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (suite *FileServiceTestSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (suite *FileServiceTestSuite) createWorkspace(orgID string, userID string) (*service.Workspace, error) {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  orgID,
		StorageCapacity: int64(config.GetConfig().Defaults.WorkspaceStorageCapacityMB),
	}, userID)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}
