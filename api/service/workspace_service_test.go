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
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

type WorkspaceServiceSuite struct {
	suite.Suite
	org   *service.Organization
	users []model.User
}

func TestWorkspaceServiceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceSuite))
}

func (s *WorkspaceServiceSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.org, err = test.CreateOrganization(s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *WorkspaceServiceSuite) TestCreate() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("workspace", workspace.Name)
	s.Equal(int64(1*GB), workspace.StorageCapacity)
}

func (s *WorkspaceServiceSuite) TestCreate_NonExistentOrganization() {
	_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  helper.NewID(),
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestFind() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewWorkspaceService().Find(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(workspace.ID, found.ID)
}

func (s *WorkspaceServiceSuite) TestFind_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().Find(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestFind_UnauthorizedUser() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().Find(workspace.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestList() {
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  s.org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewWorkspaceService().List(service.WorkspaceListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("workspace A", list.Data[0].Name)
	s.Equal("workspace B", list.Data[1].Name)
	s.Equal("workspace C", list.Data[2].Name)
}

func (s *WorkspaceServiceSuite) TestList_Paginate() {
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  s.org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewWorkspaceService().List(service.WorkspaceListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("workspace A", list.Data[0].Name)
	s.Equal("workspace B", list.Data[1].Name)

	list, err = service.NewWorkspaceService().List(service.WorkspaceListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("workspace C", list.Data[0].Name)
}

func (s *WorkspaceServiceSuite) TestList_SortByNameDescending() {
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  s.org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewWorkspaceService().List(service.WorkspaceListOptions{
		Page:      1,
		Size:      3,
		SortBy:    service.WorkspaceSortByName,
		SortOrder: service.WorkspaceSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("workspace C", list.Data[0].Name)
	s.Equal("workspace B", list.Data[1].Name)
	s.Equal("workspace A", list.Data[2].Name)
}

func (s *WorkspaceServiceSuite) TestProbe() {
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  s.org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	probe, err := service.NewWorkspaceService().Probe(service.WorkspaceListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *WorkspaceServiceSuite) TestPatchName() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	patched, err := service.NewWorkspaceService().PatchName(workspace.ID, "workspace (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(workspace.ID, patched.ID)
	s.Equal("workspace (edit)", patched.Name)
}

func (s *WorkspaceServiceSuite) TestPatchName_UnauthorisedUser() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().PatchName(workspace.ID, "workspace (edit)", s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchName_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().PatchName(uuid.NewString(), "workspace (edit)", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	workspace, err = service.NewWorkspaceService().PatchStorageCapacity(workspace.ID, int64(2*GB), s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(int64(2*GB), workspace.StorageCapacity)
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity_UnauthorisedUser() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().PatchStorageCapacity(workspace.ID, int64(1*GB), s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().PatchStorageCapacity(uuid.NewString(), int64(1*GB), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *WorkspaceServiceSuite) TestDelete_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete_UnauthorizedUser() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(workspace.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_EnoughSpace() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	hasEnoughSpace, err := service.NewWorkspaceService().HasEnoughSpaceForByteSize(workspace.ID, 512*MB)
	s.Require().NoError(err)
	s.True(*hasEnoughSpace)
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_NotEnoughSpace() {
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  s.org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	hasEnoughSpace, err := service.NewWorkspaceService().HasEnoughSpaceForByteSize(workspace.ID, 2*GB)
	s.Require().NoError(err)
	s.False(*hasEnoughSpace)
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().HasEnoughSpaceForByteSize(uuid.NewString(), 512*1024*1024)
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}
