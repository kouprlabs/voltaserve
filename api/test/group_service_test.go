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
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type GroupServiceSuite struct {
	suite.Suite
	groupSvc *service.GroupService
	org      *service.Organization
	users    []model.User
}

func TestGroupServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GroupServiceSuite))
}

func (s *GroupServiceSuite) SetupTest() {
	users, err := test_helper.CreateUsers(3)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	org, err := test_helper.CreateOrganization(users[0].GetID())
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.groupSvc = service.NewGroupService()
	s.users = users
	s.org = org
}

func (s *GroupServiceSuite) TestCreateGroup() {
	// Test creating a group with valid options
	group, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(group)
	s.Equal("group", group.Name)
	s.Equal(s.org.ID, group.Organization.ID)

	// Test creating a group with a non-existent organization
	group, err = s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "another group",
		OrganizationID: "non-existent-org-id",
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
	s.Nil(group)
}

func (s *GroupServiceSuite) TestFindGroup() {
	// Create a group to find
	group, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the created group
	foundGroup, err := s.groupSvc.Find(group.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(foundGroup)
	s.Equal(group.ID, foundGroup.ID)

	// Test finding a non-existent group
	foundGroup, err = s.groupSvc.Find("non-existent-group-id", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
	s.Nil(foundGroup)
}

func (s *GroupServiceSuite) TestListGroups() {
	// Create multiple groups to list
	_, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test listing groups with default options
	list, err := s.groupSvc.List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Equal(uint64(2), list.TotalElements)

	// Test listing groups with pagination
	list, err = s.groupSvc.List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           1,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(2), list.TotalElements)

	// Test listing groups with sorting
	list, err = s.groupSvc.List(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           10,
		SortBy:         service.GroupSortByName,
		SortOrder:      service.GroupSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(list)
	s.Equal("group B", list.Data[0].Name)
}

func (s *GroupServiceSuite) TestProbeGroups() {
	// Create multiple groups to probe
	_, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	_, err = s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probing groups
	probe, err := s.groupSvc.Probe(service.GroupListOptions{
		OrganizationID: s.org.ID,
		Page:           1,
		Size:           10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(probe)
	s.Equal(uint64(2), probe.TotalElements)
}

func (s *GroupServiceSuite) TestPatchGroupName() {
	// Create a group to patch
	group, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test patching the group name
	group, err = s.groupSvc.PatchName(group.ID, "group (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.NotNil(group)
	s.Equal("group (edit)", group.Name)

	// Test patching a non-existent group
	group, err = s.groupSvc.PatchName("non-existent-group-id", "group", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
	s.Nil(group)
}

func (s *GroupServiceSuite) TestDeleteGroup() {
	// Create a group to delete
	group, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test deleting the group
	err = s.groupSvc.Delete(group.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the deleted group
	group, err = s.groupSvc.Find(group.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
	s.Nil(group)

	// Test deleting a non-existent group
	err = s.groupSvc.Delete("non-existent-group-id", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
}

func (s *GroupServiceSuite) TestAddMember() {
	// Create a group and a user to add as a member
	group, err := s.groupSvc.Create(service.GroupCreateOptions{
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
	err = s.groupSvc.AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	// Test adding a non-existent member
	err = s.groupSvc.AddMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}

func (s *GroupServiceSuite) TestRemoveMember() {
	// Create a group and a user to add as a member
	group, err := s.groupSvc.Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
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

	err = s.groupSvc.AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	// Test removing the member from the group
	err = s.groupSvc.RemoveMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
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
	err = s.groupSvc.RemoveMember(group.ID, s.users[0].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotRemoveSoleOwnerOfGroupError(
			cache.NewGroupCache().GetOrNil(group.ID),
		).Error(),
		err.Error(),
	)

	// Test removing a non-existent member
	err = s.groupSvc.RemoveMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}
