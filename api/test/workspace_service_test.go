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
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

type WorkspaceServiceSuite struct {
	suite.Suite
	workspaceSvc *service.WorkspaceService
	org          *service.Organization
	users        []model.User
}

func TestWorkspaceServiceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceSuite))
}

func (s *WorkspaceServiceSuite) SetupTest() {
	users, err := test_helper.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := test_helper.CreateOrganization(users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.workspaceSvc = service.NewWorkspaceService()
	s.users = users
	s.org = org
}

func (s *WorkspaceServiceSuite) TestCreate() {
	// Test successful creation
	workspace, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal("workspace", workspace.Name)
	s.Equal(int64(1*GB), workspace.StorageCapacity)

	// Test invalid organization ID
	_, err = s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  "invalid-org-id",
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestFind() {
	// Create a workspace to find
	createdWorkspace, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful find
	foundWorkspace, err := s.workspaceSvc.Find(createdWorkspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(foundWorkspace)
	s.Equal(createdWorkspace.ID, foundWorkspace.ID)

	// Test non-existent workspace
	_, err = s.workspaceSvc.Find(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test unauthorized user
	_, err = s.workspaceSvc.Find(createdWorkspace.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestList() {
	// Create multiple workspaces
	for i := range 5 {
		_, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
			Name:            fmt.Sprintf("workspace %d", i),
			OrganizationID:  s.org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	// Test list all workspaces
	list, err := s.workspaceSvc.List(service.WorkspaceListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(len(list.Data), 5)

	// Test pagination
	list, err = s.workspaceSvc.List(service.WorkspaceListOptions{Page: 2, Size: 2}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Len(list.Data, 2)

	// Test sorting by name
	list, err = s.workspaceSvc.List(service.WorkspaceListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.WorkspaceSortByName,
		SortOrder: service.WorkspaceSortOrderAsc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Less(list.Data[0].Name, list.Data[1].Name)

	// Test sorting by date created
	list, err = s.workspaceSvc.List(service.WorkspaceListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.WorkspaceSortByDateCreated,
		SortOrder: service.WorkspaceSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	firstCreateTime, _ := time.Parse(time.RFC3339, list.Data[0].CreateTime)
	secondCreateTime, _ := time.Parse(time.RFC3339, list.Data[1].CreateTime)
	s.True(firstCreateTime.After(secondCreateTime))
}

func (s *WorkspaceServiceSuite) TestPatchName() {
	// Create a workspace to patch
	workspace, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful patch
	workspace, err = s.workspaceSvc.PatchName(workspace.ID, "workspace (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal("workspace (edit)", workspace.Name)

	// Test unauthorized user
	_, err = s.workspaceSvc.PatchName(workspace.ID, "workspace", s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test non-existent workspace
	_, err = s.workspaceSvc.PatchName(uuid.NewString(), "workspace", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity() {
	// Create a workspace to patch
	workspace, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful patch
	workspace, err = s.workspaceSvc.PatchStorageCapacity(workspace.ID, int64(2*GB), s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal(int64(2*GB), workspace.StorageCapacity)

	// Test unauthorized user
	_, err = s.workspaceSvc.PatchStorageCapacity(workspace.ID, int64(1*GB), s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test non-existent workspace
	_, err = s.workspaceSvc.PatchStorageCapacity(uuid.NewString(), int64(1*GB), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete() {
	// Create a workspace to delete
	workspace, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful delete
	err = s.workspaceSvc.Delete(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Test non-existent workspace
	err = s.workspaceSvc.Delete(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test unauthorized user
	err = s.workspaceSvc.Delete(workspace.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize() {
	// Create a workspace
	workspace, err := s.workspaceSvc.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test enough space
	hasEnoughSpace, err := s.workspaceSvc.HasEnoughSpaceForByteSize(workspace.ID, 512*MB)
	s.Require().NoError(err)
	s.True(*hasEnoughSpace)

	// Test not enough space
	hasEnoughSpace, err = s.workspaceSvc.HasEnoughSpaceForByteSize(workspace.ID, 2*GB)
	s.Require().NoError(err)
	s.False(*hasEnoughSpace)

	// Test non-existent workspace
	_, err = s.workspaceSvc.HasEnoughSpaceForByteSize(uuid.NewString(), 512*1024*1024)
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}
