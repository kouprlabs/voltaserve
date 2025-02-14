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

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/service"
)

type GroupServiceSuite struct {
	suite.Suite
	service *service.GroupService
	userIDs []string
	orgID   string
}

func TestGroupServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GroupServiceSuite))
}

func (s *GroupServiceSuite) SetupTest() {
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
	s.service = service.NewGroupService()
	s.userIDs = userIDs
	s.orgID = org.ID
}

func (s *GroupServiceSuite) TestCreateGroup() {
	// Test creating a group with valid options
	opts := service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.orgID,
	}
	group, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(group)
	s.Equal(opts.Name, group.Name)
	s.Equal(s.orgID, group.Organization.ID)

	// Test creating a group with a non-existent organization
	group, err = s.service.Create(service.GroupCreateOptions{
		Name:           "another group",
		OrganizationID: "non-existent-org-id",
	}, s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
	s.Nil(group)
}

func (s *GroupServiceSuite) TestFindGroup() {
	// Create a group to find
	createdGroup, err := s.service.Create(service.GroupCreateOptions{
		Name:           "group",
		OrganizationID: s.orgID,
	}, s.userIDs[0])
	s.Require().NoError(err)

	// Test finding the created group
	foundGroup, err := s.service.Find(createdGroup.ID, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(foundGroup)
	s.Equal(createdGroup.ID, foundGroup.ID)

	// Test finding a non-existent group
	foundGroup, err = s.service.Find("non-existent-group-id", s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
	s.Nil(foundGroup)
}

func (s *GroupServiceSuite) TestListGroups() {
	// Create multiple groups to list
	_, err := s.service.Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.orgID,
	}, s.userIDs[0])
	s.Require().NoError(err)
	_, err = s.service.Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.orgID,
	}, s.userIDs[0])
	s.Require().NoError(err)

	// Test listing groups with default options
	listOpts := service.GroupListOptions{
		OrganizationID: s.orgID,
		Page:           1,
		Size:           10,
	}
	groupList, err := s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(groupList)
	s.Equal(uint64(2), groupList.TotalElements)

	// Test listing groups with pagination
	listOpts.Page = 1
	listOpts.Size = 1
	groupList, err = s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(groupList)
	s.Equal(uint64(1), groupList.Size)
	s.Equal(uint64(2), groupList.TotalElements)

	// Test listing groups with sorting
	listOpts.SortBy = service.GroupSortByName
	listOpts.SortOrder = service.GroupSortOrderDesc
	groupList, err = s.service.List(listOpts, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(groupList)
	s.Equal("group B", groupList.Data[0].Name)
}

func (s *GroupServiceSuite) TestProbeGroups() {
	// Create multiple groups to probe
	_, err := s.service.Create(service.GroupCreateOptions{
		Name:           "group A",
		OrganizationID: s.orgID,
	}, s.userIDs[0])
	s.Require().NoError(err)
	_, err = s.service.Create(service.GroupCreateOptions{
		Name:           "group B",
		OrganizationID: s.orgID,
	}, s.userIDs[0])
	s.Require().NoError(err)

	// Test probing groups
	groupProbe, err := s.service.Probe(service.GroupListOptions{
		OrganizationID: s.orgID,
		Page:           1,
		Size:           10,
	}, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(groupProbe)
	s.Equal(uint64(2), groupProbe.TotalElements)
}

func (s *GroupServiceSuite) TestPatchGroupName() {
	// Create a group to patch
	opts := service.GroupCreateOptions{Name: "group", OrganizationID: s.orgID}
	createdGroup, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test patching the group name
	newName := "group (edit)"
	updatedGroup, err := s.service.PatchName(createdGroup.ID, newName, s.userIDs[0])
	s.Require().NoError(err)
	s.NotNil(updatedGroup)
	s.Equal(newName, updatedGroup.Name)

	// Test patching a non-existent group
	updatedGroup, err = s.service.PatchName("non-existent-group-id", newName, s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
	s.Nil(updatedGroup)
}

func (s *GroupServiceSuite) TestDeleteGroup() {
	// Create a group to delete
	opts := service.GroupCreateOptions{Name: "group", OrganizationID: s.orgID}
	createdGroup, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test deleting the group
	err = s.service.Delete(createdGroup.ID, s.userIDs[0])
	s.Require().NoError(err)

	// Test finding the deleted group
	foundGroup, err := s.service.Find(createdGroup.ID, s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
	s.Nil(foundGroup)

	// Test deleting a non-existent group
	err = s.service.Delete("non-existent-group-id", s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(nil).Error(), err.Error())
}

func (s *GroupServiceSuite) TestAddMember() {
	// Create a group and a user to add as a member
	opts := service.GroupCreateOptions{Name: "group", OrganizationID: s.orgID}
	createdGroup, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Add user to organization
	invitationSvc := service.NewInvitationService()
	invitations, err := invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.orgID,
		Emails:         []string{fmt.Sprintf("%s@voltaserve.com", s.userIDs[1])},
	}, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = invitationSvc.Accept(invitations[0].ID, s.userIDs[1])
	s.Require().NoError(err)

	// Test adding a member to the group
	err = s.service.AddMember(createdGroup.ID, s.userIDs[1], s.userIDs[0])
	s.Require().NoError(err)

	// Test adding a non-existent member
	err = s.service.AddMember(createdGroup.ID, s.userIDs[2], s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}

func (s *GroupServiceSuite) TestRemoveMember() {
	// Create a group and a user to add as a member
	opts := service.GroupCreateOptions{Name: "group", OrganizationID: s.orgID}
	createdGroup, err := s.service.Create(opts, s.userIDs[0])
	s.Require().NoError(err)
	createGroupModel, err := cache.NewGroupCache().Get(createdGroup.ID)
	s.Require().NoError(err)

	// Add user to organization
	invitationSvc := service.NewInvitationService()
	invitations, err := invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.orgID,
		Emails:         []string{fmt.Sprintf("%s@voltaserve.com", s.userIDs[1])},
	}, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().Len(invitations, 1)
	err = invitationSvc.Accept(invitations[0].ID, s.userIDs[1])
	s.Require().NoError(err)

	err = s.service.AddMember(createdGroup.ID, s.userIDs[1], s.userIDs[0])
	s.Require().NoError(err)

	// Test removing the member from the group
	err = s.service.RemoveMember(createdGroup.ID, s.userIDs[1], s.userIDs[0])
	s.Require().NoError(err)
	memberList, err := service.NewUserService().List(service.UserListOptions{
		GroupID: createdGroup.ID,
		Page:    1,
		Size:    10,
	}, s.userIDs[0])
	s.Require().NoError(err)
	s.Len(memberList.Data, 1)
	s.Equal(memberList.Data[0].ID, s.userIDs[0])

	// Test removing the last owner of the group
	err = s.service.RemoveMember(createdGroup.ID, s.userIDs[0], s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewCannotRemoveSoleOwnerOfGroupError(createGroupModel).Error(), err.Error())

	// Test removing a non-existent member
	err = s.service.RemoveMember(createdGroup.ID, s.userIDs[2], s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotMemberOfOrganizationError().Error(), err.Error())
}

func (s *GroupServiceSuite) createUsers() ([]string, error) {
	db, err := infra.NewPostgresManager().GetDB()
	if err != nil {
		return nil, nil
	}
	var ids []string
	for i := range 3 {
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

func (s *GroupServiceSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}
