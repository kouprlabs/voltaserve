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
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
)

type FileServiceTestSuite struct {
	suite.Suite
	fileSvc   *service.FileService
	workspace *service.Workspace
	users     []model.User
}

func (s *FileServiceTestSuite) SetupTest() {
	users, err := s.createUsers()
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := s.createOrganization(users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	workspace, err := s.createWorkspace(org.ID, users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
	}
	s.fileSvc = service.NewFileService()
	s.workspace = workspace
	s.users = users
}

func TestFileServiceSuite(t *testing.T) {
	suite.Run(t, new(FileServiceTestSuite))
}

func (s *FileServiceTestSuite) TestCreate() {
	// Test creating a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("test-file.txt", file.Name)
	s.Equal(model.FileTypeFile, file.Type)

	// Test creating a folder
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("test-folder", folder.Name)
	s.Equal(model.FileTypeFolder, folder.Type)

	// Test creating a file with an invalid parent ID
	_, err = s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "invalid-parent-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    "invalid-parent-id",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())

	// Test creating a file with a duplicate name
	_, err = s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileWithSimilarNameExistsError().Error(), err.Error())

	// Test creating a file with a duplicate name using a path name
	file, err = s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "a/b/c/test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	pathFiles, err := s.fileSvc.FindPath(file.ID, s.users[0].GetID())
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
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "find-me.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the file
	foundFiles, err := s.fileSvc.Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(foundFiles, 1)
	s.Equal(file.ID, foundFiles[0].ID)

	// Test finding a non-existent file
	foundFiles, err = s.fileSvc.Find([]string{"non-existent-id"}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(foundFiles)
}

func (s *FileServiceTestSuite) TestFindByPath() {
	// Create a folder and a file inside it
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the file by path
	foundFile, err := s.fileSvc.FindByPath(fmt.Sprintf("/%s/test-folder/test-file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(file.ID, foundFile.ID)

	// Test finding a non-existent path
	_, err = s.fileSvc.FindByPath(fmt.Sprintf("/%s/non-existent-path", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFileNotFoundError(err).Error(), err.Error())

	// Test finding the file without a leading slash
	_, err = s.fileSvc.FindByPath(fmt.Sprintf("%s/test-folder/test-file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())

	// Test finding the file with a trailing slash
	_, err = s.fileSvc.FindByPath(fmt.Sprintf("/%s/test-folder/test-file.txt/", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestListByPath() {
	// Create a folder and a file inside it
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test listing files in the folder
	files, err := s.fileSvc.ListByPath(fmt.Sprintf("/%s/test-folder", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)

	// Test listing the file
	_, err = s.fileSvc.ListByPath(fmt.Sprintf("/%s/test-folder/test-file.txt", s.workspace.ID), s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(files, 1)
	s.Equal(file.ID, files[0].ID)

	// Test listing files in the folder without a leading slash
	_, err = s.fileSvc.ListByPath(fmt.Sprintf("%s/test-folder", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathMissingLeadingSlash().Error(), err.Error())

	// Test listing the file with a trailing slash
	_, err = s.fileSvc.ListByPath(fmt.Sprintf("/%s/test-folder/test-file.txt/", s.workspace.ID), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewFilePathOfTypeFileHasTrailingSlash().Error(), err.Error())
}

func (s *FileServiceTestSuite) TestFindPath() {
	// Create a folder and a file inside it
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the path of the file
	path, err := s.fileSvc.FindPath(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(path, 3)
	s.Equal(s.workspace.RootID, path[0].ID)
	s.Equal(folder.ID, path[1].ID)
	s.Equal(file.ID, path[2].ID)
}

func (s *FileServiceTestSuite) TestProbe() {
	// Create a folder and a file inside it
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probing the folder
	probe, err := s.fileSvc.Probe(folder.ID, service.FileListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *FileServiceTestSuite) TestList() {
	// Create a folder and a file inside it
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test listing files in the folder
	list, err := s.fileSvc.List(folder.ID, service.FileListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)
	s.Equal(file.ID, list.Data[0].ID)
}

func (s *FileServiceTestSuite) TestComputeSize() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test computing the size of the file
	size, err := s.fileSvc.ComputeSize(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(size)
	s.GreaterOrEqual(*size, int64(0))
}

func (s *FileServiceTestSuite) TestCount() {
	// Create a folder and a file inside it
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folder.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test counting items in the folder
	count, err := s.fileSvc.Count(folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(count)
	s.Equal(int64(1), *count)
}

func (s *FileServiceTestSuite) TestCopy() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Create a destination folder
	folder, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-folder",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test copying the file
	copiedFile, err := s.fileSvc.Copy(file.ID, folder.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("test-file.txt", copiedFile.Name)
	s.Equal(model.FileTypeFile, copiedFile.Type)
	s.Equal(folder.ID, *copiedFile.ParentID)
}

func (s *FileServiceTestSuite) TestDelete() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test deleting the file
	err = s.fileSvc.Delete(file.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Verify the file is deleted
	foundFiles, err := s.fileSvc.Find([]string{file.ID}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(foundFiles)
}

func (s *FileServiceTestSuite) TestDownloadOriginalBuffer() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
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
	_, err = s.fileSvc.Store(file.ID, service.FileStoreOptions{Path: &path}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test downloading the file
	buf := new(bytes.Buffer)
	_, err = s.fileSvc.DownloadOriginalBuffer(file.ID, "", buf, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(string(content), buf.String())
}

func (s *FileServiceTestSuite) TestMove() {
	// Create two folders
	folderA, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder A",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	folderB, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "folder B",
		Type:        model.FileTypeFolder,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Create a file in folder1
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    folderA.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test moving the file to folder2
	movedFile, err := s.fileSvc.Move(file.ID, folderB.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(folderB.ID, *movedFile.ParentID)
}

func (s *FileServiceTestSuite) TestPatchName() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test patching the file name
	patchedFile, err := s.fileSvc.PatchName(file.ID, "new-name.txt", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("new-name.txt", patchedFile.Name)
}

func (s *FileServiceTestSuite) TestGrantUserPermission() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test granting user permission
	err = s.fileSvc.GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeUserPermission() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Grant user permission first
	err = s.fileSvc.GrantUserPermission([]string{file.ID}, s.users[1].GetID(), model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	// Test revoking user permission
	err = s.fileSvc.RevokeUserPermission([]string{file.ID}, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestGrantGroupPermission() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
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
	err = s.fileSvc.GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestRevokeGroupPermission() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
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
	err = s.fileSvc.GrantGroupPermission([]string{file.ID}, group.ID, model.PermissionViewer, s.users[0].GetID())
	s.Require().NoError(err)

	// Test revoking group permission
	err = s.fileSvc.RevokeGroupPermission([]string{file.ID}, group.ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *FileServiceTestSuite) TestReprocess() {
	// Create a file
	file, err := s.fileSvc.Create(service.FileCreateOptions{
		WorkspaceID: s.workspace.ID,
		Name:        "test-file.txt",
		Type:        model.FileTypeFile,
		ParentID:    s.workspace.RootID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Store the file
	file, err = s.fileSvc.Store(file.ID, service.FileStoreOptions{
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
	reprocessResult, err := s.fileSvc.Reprocess(file.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(reprocessResult.Accepted, 1)
}

func (s *FileServiceTestSuite) createUsers() ([]model.User, error) {
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
	var res []model.User
	userRepo := repo.NewUserRepo()
	for _, id := range ids {
		user, err := userRepo.Find(id)
		if err != nil {
			continue
		}
		res = append(res, user)
	}
	return res, nil
}

func (s *FileServiceTestSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *FileServiceTestSuite) createWorkspace(orgID string, userID string) (*service.Workspace, error) {
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
