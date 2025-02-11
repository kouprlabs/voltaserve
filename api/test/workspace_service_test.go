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
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/service"
)

type WorkspaceServiceTestSuite struct {
	suite.Suite
	service *service.WorkspaceService
	userIDs []string
	org     *service.Organization
}

func TestWorkspaceServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceTestSuite))
}

func (suite *WorkspaceServiceTestSuite) SetupTest() {
	userIDs, err := suite.createUsers()
	if err != nil {
		suite.Fail(err.Error())
	}
	org, err := suite.createOrganization(userIDs[0])
	if err != nil {
		suite.Fail(err.Error())
	}
	suite.service = service.NewWorkspaceService()
	suite.userIDs = userIDs
	suite.org = org
}

func (suite *WorkspaceServiceTestSuite) TestCreateWorkspace() {
	// Test successful creation
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  suite.org.ID,
		StorageCapacity: 1024 * 1024 * 1024, // 1GB
	}
	workspace, err := suite.service.Create(opts, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(workspace)
	suite.Equal(opts.Name, workspace.Name)
	suite.Equal(opts.StorageCapacity, workspace.StorageCapacity)

	// Test invalid organization ID
	opts.Name = "workspace"
	opts.OrganizationID = "invalid-org-id"
	_, err = suite.service.Create(opts, suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (suite *WorkspaceServiceTestSuite) TestFindWorkspace() {
	// Create a workspace to find
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  suite.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := suite.service.Create(opts, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test successful find
	foundWorkspace, err := suite.service.Find(createdWorkspace.ID, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(foundWorkspace)
	suite.Equal(createdWorkspace.ID, foundWorkspace.ID)

	// Test non-existent workspace
	_, err = suite.service.Find(uuid.NewString(), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())

	// Test unauthorized user
	_, err = suite.service.Find(createdWorkspace.ID, suite.userIDs[1])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())
}

func (suite *WorkspaceServiceTestSuite) TestListWorkspaces() {
	// Create multiple workspaces
	for i := range 5 {
		opts := service.WorkspaceCreateOptions{
			Name:            fmt.Sprintf("workspace %d", i),
			OrganizationID:  suite.org.ID,
			StorageCapacity: 1024 * 1024 * 1024,
		}
		_, err := suite.service.Create(opts, suite.userIDs[0])
		suite.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	// Test list all workspaces
	listOpts := service.WorkspaceListOptions{Page: 1, Size: 10}
	workspaceList, err := suite.service.List(listOpts, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(workspaceList)
	suite.GreaterOrEqual(len(workspaceList.Data), 5)

	// Test pagination
	listOpts.Page = 2
	listOpts.Size = 2
	workspaceList, err = suite.service.List(listOpts, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(workspaceList)
	suite.Len(workspaceList.Data, 2)

	// Test sorting by name
	listOpts.SortBy = service.WorkspaceSortByName
	listOpts.SortOrder = service.WorkspaceSortOrderAsc
	workspaceList, err = suite.service.List(listOpts, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(workspaceList)
	suite.Less(workspaceList.Data[0].Name, workspaceList.Data[1].Name)

	// Test sorting by date created
	listOpts.SortBy = service.WorkspaceSortByDateCreated
	listOpts.SortOrder = service.WorkspaceSortOrderDesc
	workspaceList, err = suite.service.List(listOpts, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(workspaceList)
	firstCreateTime, _ := time.Parse(time.RFC3339, workspaceList.Data[0].CreateTime)
	secondCreateTime, _ := time.Parse(time.RFC3339, workspaceList.Data[1].CreateTime)
	suite.True(firstCreateTime.After(secondCreateTime))
}

func (suite *WorkspaceServiceTestSuite) TestPatchWorkspaceName() {
	// Create a workspace to patch
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  suite.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := suite.service.Create(opts, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test successful patch
	newName := "workspace (edit)"
	updatedWorkspace, err := suite.service.PatchName(createdWorkspace.ID, newName, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(updatedWorkspace)
	suite.Equal(newName, updatedWorkspace.Name)

	// Test unauthorized user
	_, err = suite.service.PatchName(createdWorkspace.ID, newName, suite.userIDs[1])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())

	// Test non-existent workspace
	_, err = suite.service.PatchName(uuid.NewString(), newName, suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())
}

func (suite *WorkspaceServiceTestSuite) TestPatchWorkspaceStorageCapacity() {
	// Create a workspace to patch
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  suite.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := suite.service.Create(opts, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test successful patch
	newStorageCapacity := int64(2 * 1024 * 1024 * 1024) // 2GB
	updatedWorkspace, err := suite.service.PatchStorageCapacity(createdWorkspace.ID, newStorageCapacity, suite.userIDs[0])
	suite.Require().NoError(err)
	suite.NotNil(updatedWorkspace)
	suite.Equal(newStorageCapacity, updatedWorkspace.StorageCapacity)

	// Test unauthorized user
	_, err = suite.service.PatchStorageCapacity(createdWorkspace.ID, newStorageCapacity, suite.userIDs[1])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())

	// Test non-existent workspace
	_, err = suite.service.PatchStorageCapacity(uuid.NewString(), newStorageCapacity, suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())
}

func (suite *WorkspaceServiceTestSuite) TestDeleteWorkspace() {
	// Create a workspace to delete
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  suite.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := suite.service.Create(opts, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test successful delete
	err = suite.service.Delete(createdWorkspace.ID, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test non-existent workspace
	err = suite.service.Delete(uuid.NewString(), suite.userIDs[0])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())

	// Test unauthorized user
	err = suite.service.Delete(createdWorkspace.ID, suite.userIDs[1])
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())
}

func (suite *WorkspaceServiceTestSuite) TestHasEnoughSpaceForByteSize() {
	// Create a workspace to test
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  suite.org.ID,
		StorageCapacity: 1024 * 1024 * 1024, // 1GB
	}
	createdWorkspace, err := suite.service.Create(opts, suite.userIDs[0])
	suite.Require().NoError(err)

	// Test enough space
	hasEnoughSpace, err := suite.service.HasEnoughSpaceForByteSize(createdWorkspace.ID, 512*1024*1024) // 512MB
	suite.Require().NoError(err)
	suite.True(*hasEnoughSpace)

	// Test not enough space
	hasEnoughSpace, err = suite.service.HasEnoughSpaceForByteSize(createdWorkspace.ID, 2*1024*1024*1024) // 2GB
	suite.Require().NoError(err)
	suite.False(*hasEnoughSpace)

	// Test non-existent workspace
	_, err = suite.service.HasEnoughSpaceForByteSize(uuid.NewString(), 512*1024*1024)
	suite.Require().Error(err)
	suite.Equal(err.Error(), errorpkg.NewWorkspaceNotFoundError(err).Error())
}

func (suite *WorkspaceServiceTestSuite) createUsers() ([]string, error) {
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

func (suite *WorkspaceServiceTestSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}
