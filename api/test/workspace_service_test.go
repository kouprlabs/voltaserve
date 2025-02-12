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

type WorkspaceServiceSuite struct {
	suite.Suite
	service *service.WorkspaceService
	userIDs []string
	org     *service.Organization
}

func TestWorkspaceServiceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceSuite))
}

func (s *WorkspaceServiceSuite) SetupTest() {
	userIDs, err := s.createUsers()
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := s.createOrganization(userIDs[0])
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.service = service.NewWorkspaceService()
	s.userIDs = userIDs
	s.org = org
}

func (s *WorkspaceServiceSuite) TestCreate() {
	// Test successful creation
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1024 * 1024 * 1024, // 1GB
	}
	workspace, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal(opts.Name, workspace.Name)
	s.Equal(opts.StorageCapacity, workspace.StorageCapacity)

	// Test invalid organization ID
	opts.Name = "workspace"
	opts.OrganizationID = "invalid-org-id"
	_, err = s.service.Create(opts, s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestFind() {
	// Create a workspace to find
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test successful find
	foundWorkspace, err := s.service.Find(createdWorkspace.ID, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(foundWorkspace)
	s.Equal(createdWorkspace.ID, foundWorkspace.ID)

	// Test non-existent workspace
	_, err = s.service.Find(uuid.NewString(), s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test unauthorized user
	_, err = s.service.Find(createdWorkspace.ID, s.userIDs[1])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestList() {
	// Create multiple workspaces
	for i := range 5 {
		opts := service.WorkspaceCreateOptions{
			Name:            fmt.Sprintf("workspace %d", i),
			OrganizationID:  s.org.ID,
			StorageCapacity: 1024 * 1024 * 1024,
		}
		_, err := s.service.Create(opts, s.userIDs[0])
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	// Test list all workspaces
	listOpts := service.WorkspaceListOptions{Page: 1, Size: 10}
	workspaceList, err := s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(workspaceList)
	s.GreaterOrEqual(len(workspaceList.Data), 5)

	// Test pagination
	listOpts.Page = 2
	listOpts.Size = 2
	workspaceList, err = s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(workspaceList)
	s.Len(workspaceList.Data, 2)

	// Test sorting by name
	listOpts.SortBy = service.WorkspaceSortByName
	listOpts.SortOrder = service.WorkspaceSortOrderAsc
	workspaceList, err = s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(workspaceList)
	s.Less(workspaceList.Data[0].Name, workspaceList.Data[1].Name)

	// Test sorting by date created
	listOpts.SortBy = service.WorkspaceSortByDateCreated
	listOpts.SortOrder = service.WorkspaceSortOrderDesc
	workspaceList, err = s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(workspaceList)
	firstCreateTime, _ := time.Parse(time.RFC3339, workspaceList.Data[0].CreateTime)
	secondCreateTime, _ := time.Parse(time.RFC3339, workspaceList.Data[1].CreateTime)
	s.True(firstCreateTime.After(secondCreateTime))
}

func (s *WorkspaceServiceSuite) TestPatchName() {
	// Create a workspace to patch
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test successful patch
	newName := "workspace (edit)"
	updatedWorkspace, err := s.service.PatchName(createdWorkspace.ID, newName, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(updatedWorkspace)
	s.Equal(newName, updatedWorkspace.Name)

	// Test unauthorized user
	_, err = s.service.PatchName(createdWorkspace.ID, newName, s.userIDs[1])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test non-existent workspace
	_, err = s.service.PatchName(uuid.NewString(), newName, s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity() {
	// Create a workspace to patch
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test successful patch
	newStorageCapacity := int64(2 * 1024 * 1024 * 1024) // 2GB
	updatedWorkspace, err := s.service.PatchStorageCapacity(createdWorkspace.ID, newStorageCapacity, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(updatedWorkspace)
	s.Equal(newStorageCapacity, updatedWorkspace.StorageCapacity)

	// Test unauthorized user
	_, err = s.service.PatchStorageCapacity(createdWorkspace.ID, newStorageCapacity, s.userIDs[1])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test non-existent workspace
	_, err = s.service.PatchStorageCapacity(uuid.NewString(), newStorageCapacity, s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete() {
	// Create a workspace to delete
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1024 * 1024 * 1024,
	}
	createdWorkspace, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test successful delete
	err = s.service.Delete(createdWorkspace.ID, s.userIDs[0])
	s.Require().NoError(err)

	// Test non-existent workspace
	err = s.service.Delete(uuid.NewString(), s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test unauthorized user
	err = s.service.Delete(createdWorkspace.ID, s.userIDs[1])
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize() {
	// Create a workspace to test
	opts := service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1024 * 1024 * 1024, // 1GB
	}
	createdWorkspace, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test enough space
	hasEnoughSpace, err := s.service.HasEnoughSpaceForByteSize(createdWorkspace.ID, 512*1024*1024) // 512MB
	s.Require().NoError(err)
	s.True(*hasEnoughSpace)

	// Test not enough space
	hasEnoughSpace, err = s.service.HasEnoughSpaceForByteSize(createdWorkspace.ID, 2*1024*1024*1024) // 2GB
	s.Require().NoError(err)
	s.False(*hasEnoughSpace)

	// Test non-existent workspace
	_, err = s.service.HasEnoughSpaceForByteSize(uuid.NewString(), 512*1024*1024)
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) createUsers() ([]string, error) {
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

func (s *WorkspaceServiceSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}
