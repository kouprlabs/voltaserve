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

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type InvitationServiceSuite struct {
	suite.Suite
	invitationSvc *service.InvitationService
	users         []model.User
	org           *service.Organization
}

func TestInvitationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InvitationServiceSuite))
}

func (s *InvitationServiceSuite) SetupTest() {
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
	s.invitationSvc = service.NewInvitationService()
	s.users = users
	s.org = org
}

func (s *InvitationServiceSuite) TestCreateInvitation() {
	// Test successful creation of invitations
	opts := service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{"test1@example.com", "test2@example.com"},
	}
	invitations, err := s.invitationSvc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(invitations, 2)

	// Test duplicate emails should not create new invitations
	invitations, err = s.invitationSvc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(invitations)

	// Test invalid organization ID
	opts.OrganizationID = "invalid-org"
	_, err = s.invitationSvc.Create(opts, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())

	// Test user not authorized to create invitations
	opts.OrganizationID = s.org.ID
	_, err = s.invitationSvc.Create(opts, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestListIncomingInvitations() {
	// Create some invitations
	opts := service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}
	_, err := s.invitationSvc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)

	// Test listing incoming invitations
	listOpts := service.InvitationListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.InvitationSortByEmail,
		SortOrder: service.InvitationSortOrderAsc,
	}
	list, err := s.invitationSvc.ListIncoming(listOpts, s.users[1].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)

	// Test pagination
	listOpts.Page = 2
	list, err = s.invitationSvc.ListIncoming(listOpts, s.users[1].GetID())
	s.Require().NoError(err)
	s.Empty(list.Data)
}

func (s *InvitationServiceSuite) TestProbeIncomingInvitations() {
	// Create some invitations
	_, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probe incoming invitations
	probe, err := s.invitationSvc.ProbeIncoming(service.InvitationListOptions{Page: 1, Size: 10}, s.users[1].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
}

func (s *InvitationServiceSuite) TestCountIncomingInvitations() {
	// Create some invitations
	_, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test count incoming invitations
	count, err := s.invitationSvc.CountIncoming(s.users[1].GetID())
	s.Require().NoError(err)
	s.Equal(int64(1), *count)
}

func (s *InvitationServiceSuite) TestListOutgoingInvitations() {
	// Create some invitations
	_, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// List outgoing invitations
	list, err := s.invitationSvc.ListOutgoing(s.org.ID, service.InvitationListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)

	// Test pagination
	list, err = s.invitationSvc.ListOutgoing(s.org.ID, service.InvitationListOptions{Page: 2, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(list.Data)
}

func (s *InvitationServiceSuite) TestProbeOutgoingInvitations() {
	// Create some invitations
	_, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probe outgoing invitations
	probe, err := s.invitationSvc.ProbeOutgoing(s.org.ID, service.InvitationListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
}

func (s *InvitationServiceSuite) TestAcceptInvitation() {
	// Create an invitation
	invitations, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test accept invitation
	err = s.invitationSvc.Accept(invitationID, s.users[1].GetID())
	s.Require().NoError(err)

	// Test cannot accept already accepted invitation
	err = s.invitationSvc.Accept(invitationID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotAcceptNonPendingInvitationError(
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)

	// Create another invitation
	invitations, err = s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	invitationID = invitations[0].ID
	s.Require().NoError(err)

	// Test user not allowed to accept invitation
	err = s.invitationSvc.Accept(invitationID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToAcceptInvitationError(
			s.users[0],
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestDeclineInvitation() {
	// Create an invitation
	invitations, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test decline invitation
	err = s.invitationSvc.Decline(invitationID, s.users[1].GetID())
	s.Require().NoError(err)

	// Test cannot decline already declined invitation
	err = s.invitationSvc.Decline(invitationID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotDeclineNonPendingInvitationError(
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)

	// Create another invitation
	invitations, err = s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID = invitations[0].ID

	// Test user not allowed to decline invitation
	err = s.invitationSvc.Decline(invitationID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToDeclineInvitationError(
			s.users[0],
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestResendInvitation() {
	// Create an invitation
	invitations, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test resend invitation
	err = s.invitationSvc.Resend(invitationID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDeleteInvitation() {
	// Create an invitation
	invitations, err := s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test delete invitation
	err = s.invitationSvc.Delete(invitationID, s.users[0].GetID())
	s.Require().NoError(err)

	// Create another invitation
	invitations, err = s.invitationSvc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID = invitations[0].ID

	// Test user not allowed to delete invitation
	err = s.invitationSvc.Delete(invitationID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToDeleteInvitationError(
			s.users[1],
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)
}
