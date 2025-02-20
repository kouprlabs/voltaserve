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

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

type UserServiceSuite struct {
	suite.Suite
	users []model.User
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(3)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	for _, user := range s.users {
		err := infra.NewSearchManager().Index(infra.UserSearchIndex, []infra.SearchModel{user})
		s.Require().NoError(err)
	}
}

func (s *UserServiceSuite) TestList_GroupMembers() {
	group := s.createGroup()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:    1,
		Size:    10,
		GroupID: group.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(s.users[0].GetID(), list.Data[0].ID)
	s.Equal(s.users[1].GetID(), list.Data[1].ID)
	s.Equal(s.users[2].GetID(), list.Data[2].ID)
}

func (s *UserServiceSuite) TestList_GroupMembers_MissingGroupPermission() {
	group := s.createGroup()

	s.revokeUserPermissionForGroup(group, s.users[0])

	_, err := service.NewUserService().List(service.UserListOptions{
		Page:    1,
		Size:    10,
		GroupID: group.ID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *UserServiceSuite) TestList_GroupMembersPaginate() {
	group := s.createGroup()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:    1,
		Size:    2,
		GroupID: group.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(s.users[0].GetID(), list.Data[0].ID)
	s.Equal(s.users[1].GetID(), list.Data[1].ID)

	list, err = service.NewUserService().List(service.UserListOptions{
		Page:    2,
		Size:    2,
		GroupID: group.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(s.users[2].GetID(), list.Data[0].ID)
}

func (s *UserServiceSuite) TestList_GroupMembersSortByEmailDescending() {
	group := s.createGroup()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:      1,
		Size:      3,
		GroupID:   group.ID,
		SortBy:    service.UserSortByEmail,
		SortOrder: service.UserSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(s.users[2].GetID(), list.Data[0].ID)
	s.Equal(s.users[1].GetID(), list.Data[1].ID)
	s.Equal(s.users[0].GetID(), list.Data[2].ID)
}

func (s *UserServiceSuite) TestList_GroupMembersQuery() {
	group := s.createGroup()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:    1,
		Size:    10,
		GroupID: group.ID,
		Query:   s.users[2].GetID(),
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(1), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(s.users[2].GetID(), list.Data[0].ID)
}

func (s *UserServiceSuite) TestList_OrganizationMembers() {
	org := s.createOrganization()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:           1,
		Size:           10,
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(s.users[0].GetID(), list.Data[0].ID)
	s.Equal(s.users[1].GetID(), list.Data[1].ID)
	s.Equal(s.users[2].GetID(), list.Data[2].ID)
}

func (s *UserServiceSuite) TestList_OrganizationMembers_MissingOrganizationPermission() {
	org := s.createOrganization()

	s.revokeUserPermissionForOrganization(org, s.users[0])

	_, err := service.NewUserService().List(service.UserListOptions{
		Page:           1,
		Size:           10,
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *UserServiceSuite) TestList_OrganizationMembersPaginate() {
	org := s.createOrganization()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:           1,
		Size:           2,
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(s.users[0].GetID(), list.Data[0].ID)
	s.Equal(s.users[1].GetID(), list.Data[1].ID)

	list, err = service.NewUserService().List(service.UserListOptions{
		Page:           2,
		Size:           2,
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(s.users[2].GetID(), list.Data[0].ID)
}

func (s *UserServiceSuite) TestList_OrganizationMembersSortByEmailDescending() {
	org := s.createOrganization()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:           1,
		Size:           3,
		OrganizationID: org.ID,
		SortBy:         service.UserSortByEmail,
		SortOrder:      service.UserSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(s.users[2].GetID(), list.Data[0].ID)
	s.Equal(s.users[1].GetID(), list.Data[1].ID)
	s.Equal(s.users[0].GetID(), list.Data[2].ID)
}

func (s *UserServiceSuite) TestList_OrganizationMembersQuery() {
	org := s.createOrganization()

	list, err := service.NewUserService().List(service.UserListOptions{
		Page:           1,
		Size:           10,
		OrganizationID: org.ID,
		Query:          s.users[2].GetID(),
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(1), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(s.users[2].GetID(), list.Data[0].ID)
}

func (s *UserServiceSuite) TestProbe_GroupMembers() {
	group := s.createGroup()

	probe, err := service.NewUserService().Probe(service.UserListOptions{
		Page:    1,
		Size:    10,
		GroupID: group.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *UserServiceSuite) TestProbe_GroupMembers_MissingGroupPermission() {
	group := s.createGroup()

	s.revokeUserPermissionForGroup(group, s.users[0])

	_, err := service.NewUserService().Probe(service.UserListOptions{
		Page:    1,
		Size:    10,
		GroupID: group.ID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewGroupNotFoundError(err).Error(), err.Error())
}

func (s *UserServiceSuite) TestProbe_OrganizationMembers() {
	org := s.createOrganization()

	probe, err := service.NewUserService().Probe(service.UserListOptions{
		Page:           1,
		Size:           10,
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *UserServiceSuite) TestProbe_OrganizationMembers_MissingOrganizationPermission() {
	org := s.createOrganization()

	s.revokeUserPermissionForOrganization(org, s.users[0])

	_, err := service.NewUserService().Probe(service.UserListOptions{
		Page:           1,
		Size:           10,
		OrganizationID: org.ID,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *UserServiceSuite) createGroup() *service.Group {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	group, err := test.CreateGroup(org.ID, s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails: []string{
			s.users[1].GetEmail(),
			s.users[2].GetEmail(),
		},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)
	err = service.NewInvitationService().Accept(invitations[1].ID, s.users[2].GetID())
	s.Require().NoError(err)

	err = service.NewGroupService().AddMember(group.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
	err = service.NewGroupService().AddMember(group.ID, s.users[2].GetID(), s.users[0].GetID())
	s.Require().NoError(err)

	return group
}

func (s *UserServiceSuite) createOrganization() *service.Organization {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails: []string{
			s.users[1].GetEmail(),
			s.users[2].GetEmail(),
		},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)
	err = service.NewInvitationService().Accept(invitations[1].ID, s.users[2].GetID())
	s.Require().NoError(err)

	return org
}

func (s *UserServiceSuite) revokeUserPermissionForOrganization(org *service.Organization, user model.User) {
	err := repo.NewOrganizationRepo().RevokeUserPermission(org.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewOrganizationCache().Refresh(org.ID)
	s.Require().NoError(err)
}

func (s *UserServiceSuite) revokeUserPermissionForGroup(group *service.Group, user model.User) {
	err := repo.NewGroupRepo().RevokeUserPermission(group.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewGroupCache().Refresh(group.ID)
	s.Require().NoError(err)
}
