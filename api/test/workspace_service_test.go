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
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

type WorkspaceServiceSuite struct {
	suite.Suite
	service *service.WorkspaceService
	org     *service.Organization
	users   []model.User
}

func TestWorkspaceServiceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceSuite))
}

func (s *WorkspaceServiceSuite) SetupTest() {
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
	s.service = service.NewWorkspaceService()
	s.users = users
	s.org = org
}

func (s *WorkspaceServiceSuite) TestCreate() {
	// Test successful creation
	workspace, err := s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal("workspace", workspace.Name)
	s.Equal(int64(1*GB), workspace.StorageCapacity)

	// Test invalid organization ID
	_, err = s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  "invalid-org-id",
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestFind() {
	// Create a workspace to find
	createdWorkspace, err := s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful find
	foundWorkspace, err := s.service.Find(createdWorkspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(foundWorkspace)
	s.Equal(createdWorkspace.ID, foundWorkspace.ID)

	// Test non-existent workspace
	_, err = s.service.Find(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test unauthorized user
	_, err = s.service.Find(createdWorkspace.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestList() {
	// Create multiple workspaces
	for i := range 5 {
		_, err := s.service.Create(service.WorkspaceCreateOptions{
			Name:            fmt.Sprintf("workspace %d", i),
			OrganizationID:  s.org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	// Test list all workspaces
	list, err := s.service.List(service.WorkspaceListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.GreaterOrEqual(len(list.Data), 5)

	// Test pagination
	list, err = s.service.List(service.WorkspaceListOptions{Page: 2, Size: 2}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Len(list.Data, 2)

	// Test sorting by name
	list, err = s.service.List(service.WorkspaceListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.WorkspaceSortByName,
		SortOrder: service.WorkspaceSortOrderAsc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Less(list.Data[0].Name, list.Data[1].Name)

	// Test sorting by date created
	list, err = s.service.List(service.WorkspaceListOptions{
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
	workspace, err := s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful patch
	workspace, err = s.service.PatchName(workspace.ID, "workspace (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal("workspace (edit)", workspace.Name)

	// Test unauthorized user
	_, err = s.service.PatchName(workspace.ID, "workspace", s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test non-existent workspace
	_, err = s.service.PatchName(uuid.NewString(), "workspace", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity() {
	// Create a workspace to patch
	workspace, err := s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful patch
	workspace, err = s.service.PatchStorageCapacity(workspace.ID, int64(2*GB), s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(workspace)
	s.Equal(int64(2*GB), workspace.StorageCapacity)

	// Test unauthorized user
	_, err = s.service.PatchStorageCapacity(workspace.ID, int64(1*GB), s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test non-existent workspace
	_, err = s.service.PatchStorageCapacity(uuid.NewString(), int64(1*GB), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete() {
	// Create a workspace to delete
	workspace, err := s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test successful delete
	err = s.service.Delete(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Test non-existent workspace
	err = s.service.Delete(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())

	// Test unauthorized user
	err = s.service.Delete(workspace.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize() {
	// Create a workspace to test
	workspace, err := s.service.Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test enough space
	hasEnoughSpace, err := s.service.HasEnoughSpaceForByteSize(workspace.ID, 512*MB)
	s.Require().NoError(err)
	s.True(*hasEnoughSpace)

	// Test not enough space
	hasEnoughSpace, err = s.service.HasEnoughSpaceForByteSize(workspace.ID, 2*GB)
	s.Require().NoError(err)
	s.False(*hasEnoughSpace)

	// Test non-existent workspace
	_, err = s.service.HasEnoughSpaceForByteSize(uuid.NewString(), 512*1024*1024)
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) createUsers() ([]model.User, error) {
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

func (s *WorkspaceServiceSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}
