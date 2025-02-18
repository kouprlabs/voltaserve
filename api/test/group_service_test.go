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
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type GroupServiceSuite struct {
	suite.Suite
	org   *service.Organization
	users []model.User
}

func TestGroupServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GroupServiceSuite))
}

func (s *GroupServiceSuite) SetupTest() {
	var err error
	s.users, err = test_helper.CreateUsers(3)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.org, err = test_helper.CreateOrganization(s.users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *GroupServiceSuite) TestCreate() {
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group", group.Name)
	s.Equal(s.org.ID, group.Organization.ID)
}

func (s *GroupServiceSuite) TestCreate_NonExistentOrganization() {
	_, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "another group",
		OrganizationID: "non-existent-org-id",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *GroupServiceSuite) TestFind() {
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewGroupService().Find(group.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(group.ID, found.ID)
}

func (s *GroupServiceSuite) TestFind_NonExistentGroup() {
	_, err := service.NewGroupService().Find(helper.NewID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
}

func (s *GroupServiceSuite) TestList() {
	_, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewGroupService().List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.TotalElements)
}

func (s *GroupServiceSuite) TestList_Pagination() {
	_, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewGroupService().List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           1,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group A", list.Data[0].Name)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)

	list, err = service.NewGroupService().List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           2,
		Size:           1,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group B", list.Data[0].Name)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(2), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
}

func (s *GroupServiceSuite) TestList_Sorting() {
	_, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewGroupService().List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           10,
		SortBy:         service.GroupSortByName,
		SortOrder:      service.GroupSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group B", list.Data[0].Name)
	s.Equal("group A", list.Data[1].Name)
}

func (s *GroupServiceSuite) TestProbe() {
	_, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	probe, err := service.NewGroupService().Probe(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), probe.TotalElements)
}

func (s *GroupServiceSuite) TestPatchName() {
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	group, err = service.NewGroupService().PatchName(group.ID, "group (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("group (edit)", group.Name)
}

func (s *GroupServiceSuite) TestPatchName_NonExistentGroup() {
	_, err := service.NewGroupService().PatchName(helper.NewID(), "group", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
}

func (s *GroupServiceSuite) TestDelete() {
	// Create a group to delete
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test deleting the group
	err = service.NewGroupService().Delete(group.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the deleted group
	_, err = service.NewGroupService().Find(group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())

	// Test deleting a non-existent group
	err = service.NewGroupService().Delete("non-existent-group-id", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
}

func (s *GroupServiceSuite) TestAddMember() {
	// Create a group and a user to add as a member
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Add user to organization
	invitationSvc := service.NewInvitationService()
	invitations, err := invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = invitationSvc.Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	// Test adding a member to the group
	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	// Test adding a non-existent member
	err = service.NewGroupService().AddMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}

func (s *GroupServiceSuite) TestRemoveMember() {
	// Create a group and a user to add as a member
	group, err := service.NewGroupService().Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NoError(err)

	// Add user to organization
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	// Test removing the member from the group
	err = service.NewGroupService().RemoveMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
	memberList, err := service.NewUserService().List(service.UserListOptions{
		GroupID: group.ID,
		Page:    1,
		Size:    10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(memberList.Data, 1)
	s.Equal(memberList.Data[0].ID, s.users[0].GetID())

	// Test removing the last owner of the group
	err = service.NewGroupService().RemoveMember(group.ID, s.users[0].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotRemoveSoleOwnerOfGroupError(
			cache.NewGroupCache().GetOrNil(group.ID),
		).Error(),
		err.Error(),
	)

	// Test removing a non-existent member
	err = service.NewGroupService().RemoveMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}
