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

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

type GroupServiceSuite struct {
	suite.Suite
	users []model.User
}

func TestGroupServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GroupServiceSuite))
}

func (s *GroupServiceSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(3)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *GroupServiceSuite) TestCreate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group", group.Name)
	s.Equal(org.ID, group.Organization.ID)
}

func (s *GroupServiceSuite) TestCreate_MissingOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForOrganization(org, s.users[0])

	_, err = service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestCreate_NonExistentOrganization() {
	_, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "another group",
		OrganizationID: "non-existent-org-id",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestFind() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewGroupService().Find(group.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(group.ID, found.ID)
}

func (s *GroupServiceSuite) TestFind_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	_, err = service.NewGroupService().Find(group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestFind_NonExistentGroup() {
	_, err := service.NewGroupService().Find(helper.NewID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestList() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"group A", "group B", "group C"} {
		_, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewGroupService().List(dto.GroupListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("group A", list.Data[0].Name)
	s.Equal("group B", list.Data[1].Name)
	s.Equal("group C", list.Data[2].Name)
}

func (s *GroupServiceSuite) TestList_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	var groups []*dto.Group
	for _, name := range []string{"group A", "group B", "group C"} {
		g, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		groups = append(groups, g)
		time.Sleep(1 * time.Second)
	}

	s.revokeUserPermissionForGroup(groups[1], s.users[0])

	list, err := service.NewGroupService().List(dto.GroupListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("group A", list.Data[0].Name)
	s.Equal("group C", list.Data[1].Name)
}

func (s *GroupServiceSuite) TestList_Paginate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"group A", "group B", "group C"} {
		_, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewGroupService().List(dto.GroupListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("group A", list.Data[0].Name)
	s.Equal("group B", list.Data[1].Name)

	list, err = service.NewGroupService().List(dto.GroupListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("group C", list.Data[0].Name)
}

func (s *GroupServiceSuite) TestList_SortByNameDescending() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"group A", "group B", "group C"} {
		_, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewGroupService().List(dto.GroupListOptions{
		Page:      1,
		Size:      3,
		SortBy:    dto.WorkspaceSortByName,
		SortOrder: dto.WorkspaceSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group C", list.Data[0].Name)
	s.Equal("group B", list.Data[1].Name)
	s.Equal("group A", list.Data[2].Name)
}

func (s *GroupServiceSuite) TestList_Query() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"foo bar", "hello world", "lorem ipsum"} {
		_, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewGroupService().List(dto.GroupListOptions{
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

func (s *GroupServiceSuite) TestProbe() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	for _, name := range []string{"group A", "group B", "group C"} {
		_, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	probe, err := service.NewGroupService().Probe(dto.GroupListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *GroupServiceSuite) TestProbe_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	var groups []*dto.Group
	for _, name := range []string{"group A", "group B", "group C"} {
		g, err := service.NewGroupService().Create(dto.GroupCreateOptions{
			Name:           name,
			OrganizationID: org.ID,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		groups = append(groups, g)
	}

	s.revokeUserPermissionForGroup(groups[1], s.users[0])

	probe, err := service.NewGroupService().Probe(dto.GroupListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *GroupServiceSuite) TestPatchName() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	group, err = service.NewGroupService().PatchName(group.ID, "group (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group (edit)", group.Name)
}

func (s *GroupServiceSuite) TestPatchName_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	_, err = service.NewGroupService().PatchName(group.ID, "group (edit)", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestPatchName_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForGroup(group, s.users[0], model.PermissionViewer)

	_, err = service.NewGroupService().PatchName(group.ID, "group (edit)", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewGroupPermissionError(
			s.users[0].GetID(),
			cache.NewGroupCache().GetOrNil(group.ID),
			model.PermissionEditor,
		).Error(),
		err.Error(),
	)
}

func (s *GroupServiceSuite) TestPatchName_NonExistentGroup() {
	_, err := service.NewGroupService().PatchName(helper.NewID(), "group", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestDelete() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().Delete(group.ID, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewGroupService().Find(group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestDelete_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	err = service.NewGroupService().Delete(group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestDelete_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForGroup(group, s.users[0], model.PermissionViewer)

	err = service.NewGroupService().Delete(group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewGroupPermissionError(
			s.users[0].GetID(),
			cache.NewGroupCache().GetOrNil(group.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *GroupServiceSuite) TestDelete_NonExistentGroup() {
	err := service.NewGroupService().Delete(helper.NewID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestAddMember() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *GroupServiceSuite) TestAddMember_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestAddMember_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForGroup(group, s.users[0], model.PermissionViewer)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewGroupPermissionError(
			s.users[0].GetID(),
			cache.NewGroupCache().GetOrNil(group.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *GroupServiceSuite) TestAddMember_NonMemberOfOrganization() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}

func (s *GroupServiceSuite) TestRemoveMember() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().RemoveMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	memberList, err := service.NewUserService().List(dto.UserListOptions{
		GroupID: group.ID,
		Page:    1,
		Size:    10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(memberList.Data, 1)
	s.Equal(memberList.Data[0].ID, s.users[0].GetID())
}

func (s *GroupServiceSuite) TestRemoveMember_MissingPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForGroup(group, s.users[0])

	err = service.NewGroupService().RemoveMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *GroupServiceSuite) TestRemoveMember_InsufficientPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForGroup(group, s.users[0], model.PermissionViewer)

	err = service.NewGroupService().RemoveMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewGroupPermissionError(
			s.users[0].GetID(),
			cache.NewGroupCache().GetOrNil(group.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *GroupServiceSuite) TestRemoveMember_LastOwnerOfGroup() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().RemoveMember(group.ID, s.users[0].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotRemoveSoleOwnerOfGroupError(
			cache.NewGroupCache().GetOrNil(group.ID),
		).Error(),
		err.Error(),
	)
}

func (s *GroupServiceSuite) TestRemoveMember_NonMemberOfOrganization() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := service.NewGroupService().Create(dto.GroupCreateOptions{
		Name:           "group",
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NoError(err)

	err = service.NewGroupService().RemoveMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}

func (s *GroupServiceSuite) grantUserPermissionForGroup(group *dto.Group, user model.User, permission string) {
	err := repo.NewGroupRepo().GrantUserPermission(group.ID, user.GetID(), permission)
	s.Require().NoError(err)
	_, err = cache.NewGroupCache().Refresh(group.ID)
	s.Require().NoError(err)
}

func (s *GroupServiceSuite) revokeUserPermissionForGroup(group *dto.Group, user model.User) {
	err := repo.NewGroupRepo().RevokeUserPermission(group.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewGroupCache().Refresh(group.ID)
	s.Require().NoError(err)
}

func (s *GroupServiceSuite) revokeUserPermissionForOrganization(org *dto.Organization, user model.User) {
	err := repo.NewOrganizationRepo().RevokeUserPermission(org.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewOrganizationCache().Refresh(org.ID)
	s.Require().NoError(err)
}
