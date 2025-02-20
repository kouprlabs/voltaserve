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

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

type WorkspaceServiceSuite struct {
	suite.Suite
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
}

func (s *WorkspaceServiceSuite) TestCreate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("workspace", workspace.Name)
	s.Equal(int64(1*GB), workspace.StorageCapacity)
}

func (s *WorkspaceServiceSuite) TestCreate_MissingOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewOrganizationRepo().RevokeUserPermission(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewOrganizationCache().Refresh(org.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewWorkspaceService().Find(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(workspace.ID, found.ID)
}

func (s *WorkspaceServiceSuite) TestFind_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().Find(workspace.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestFind_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().Find(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestList() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
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

func (s *WorkspaceServiceSuite) TestList_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	var workspaces []*service.Workspace
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		w, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		workspaces = append(workspaces, w)
		time.Sleep(1 * time.Second)
	}

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspaces[1].ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspaces[1].ID)
	s.Require().NoError(err)

	list, err := service.NewWorkspaceService().List(service.WorkspaceListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("workspace A", list.Data[0].Name)
	s.Equal("workspace C", list.Data[1].Name)
}

func (s *WorkspaceServiceSuite) TestList_Paginate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
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

func (s *WorkspaceServiceSuite) TestList_Query() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"foo bar", "hello world", "lorem ipsum"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewWorkspaceService().List(service.WorkspaceListOptions{
		Query: "world",
		Page:  1,
		Size:  10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(1), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("hello world", list.Data[0].Name)
}

func (s *WorkspaceServiceSuite) TestProbe() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		_, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
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

func (s *WorkspaceServiceSuite) TestProbe_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	var workspaces []*service.Workspace
	for _, name := range []string{"workspace A", "workspace B", "workspace C"} {
		w, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
			Name:            name,
			OrganizationID:  org.ID,
			StorageCapacity: 1 * GB,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		workspaces = append(workspaces, w)
	}

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspaces[1].ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspaces[1].ID)
	s.Require().NoError(err)

	probe, err := service.NewWorkspaceService().Probe(service.WorkspaceListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *WorkspaceServiceSuite) TestPatchName() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	patched, err := service.NewWorkspaceService().PatchName(workspace.ID, "workspace (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(workspace.ID, patched.ID)
	s.Equal("workspace (edit)", patched.Name)
}

func (s *WorkspaceServiceSuite) TestPatchName_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().PatchName(workspace.ID, "workspace (edit)", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchName_InsufficientPermissions() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().GrantUserPermission(workspace.ID, s.users[0].GetID(), model.PermissionViewer)
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().PatchName(workspace.ID, "workspace (edit)", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewWorkspacePermissionError(
			s.users[0].GetID(),
			cache.NewWorkspaceCache().GetOrNil(workspace.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *WorkspaceServiceSuite) TestPatchName_UnauthorisedUser() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	workspace, err = service.NewWorkspaceService().PatchStorageCapacity(workspace.ID, int64(2*GB), s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(int64(2*GB), workspace.StorageCapacity)
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().PatchStorageCapacity(workspace.ID, int64(2*GB), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().GrantUserPermission(workspace.ID, s.users[0].GetID(), model.PermissionViewer)
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().PatchStorageCapacity(workspace.ID, int64(2*GB), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewWorkspacePermissionError(
			s.users[0].GetID(),
			cache.NewWorkspaceCache().GetOrNil(workspace.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *WorkspaceServiceSuite) TestPatchStorageCapacity_UnauthorisedUser() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().Find(workspace.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(workspace.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestDelete_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().GrantUserPermission(workspace.ID, s.users[0].GetID(), model.PermissionViewer)
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(workspace.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewWorkspacePermissionError(
			s.users[0].GetID(),
			cache.NewWorkspaceCache().GetOrNil(workspace.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *WorkspaceServiceSuite) TestDelete_NonExistentWorkspace() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewWorkspaceService().Delete(uuid.NewString(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_EnoughSpace() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	hasEnoughSpace, err := service.NewWorkspaceService().HasEnoughSpaceForByteSize(workspace.ID, 512*MB, s.users[0].GetID())
	s.Require().NoError(err)
	s.True(*hasEnoughSpace)
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewWorkspaceRepo().RevokeUserPermission(workspace.ID, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = cache.NewWorkspaceCache().Refresh(workspace.ID)
	s.Require().NoError(err)

	_, err = service.NewWorkspaceService().HasEnoughSpaceForByteSize(workspace.ID, 512*MB, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_NotEnoughSpace() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	workspace, err := service.NewWorkspaceService().Create(service.WorkspaceCreateOptions{
		Name:            "workspace",
		OrganizationID:  org.ID,
		StorageCapacity: 1 * GB,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	hasEnoughSpace, err := service.NewWorkspaceService().HasEnoughSpaceForByteSize(workspace.ID, 2*GB, s.users[0].GetID())
	s.Require().NoError(err)
	s.False(*hasEnoughSpace)
}

func (s *WorkspaceServiceSuite) TestHasEnoughSpaceForByteSize_NonExistentWorkspace() {
	_, err := service.NewWorkspaceService().HasEnoughSpaceForByteSize(uuid.NewString(), 512*1024*1024, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewWorkspaceNotFoundError(err).Error(), err.Error())
}
