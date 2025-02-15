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
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"testing"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/stretchr/testify/suite"
)

type InvitationServiceSuite struct {
	suite.Suite
	svc   *service.InvitationService
	users []model.User
	org   *service.Organization
}

func TestInvitationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InvitationServiceSuite))
}

func (s *InvitationServiceSuite) SetupTest() {
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
	s.svc = service.NewInvitationService()
	s.users = users
	s.org = org
}

func (s *InvitationServiceSuite) TestCreateInvitation() {
	// Test successful creation of invitations
	opts := service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{"test1@example.com", "test2@example.com"},
	}
	invitations, err := s.svc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(invitations, 2)

	// Test duplicate emails should not create new invitations
	invitations, err = s.svc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(invitations, 0)

	// Test invalid organization ID
	opts.OrganizationID = "invalid-org"
	_, err = s.svc.Create(opts, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())

	// Test user not authorized to create invitations
	opts.OrganizationID = s.org.ID
	_, err = s.svc.Create(opts, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestListIncomingInvitations() {
	// Create some invitations
	opts := service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}
	_, err := s.svc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)

	// Test list incoming invitations
	listOpts := service.InvitationListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.InvitationSortByEmail,
		SortOrder: service.InvitationSortOrderAsc,
	}
	list, err := s.svc.ListIncoming(listOpts, s.users[1].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)

	// Test pagination
	listOpts.Page = 2
	list, err = s.svc.ListIncoming(listOpts, s.users[1].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 0)
}

func (s *InvitationServiceSuite) TestProbeIncomingInvitations() {
	// Create some invitations
	_, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probe incoming invitations
	probe, err := s.svc.ProbeIncoming(service.InvitationListOptions{Page: 1, Size: 10}, s.users[1].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
}

func (s *InvitationServiceSuite) TestCountIncomingInvitations() {
	// Create some invitations
	_, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test count incoming invitations
	count, err := s.svc.CountIncoming(s.users[1].GetID())
	s.Require().NoError(err)
	s.Equal(int64(1), *count)
}

func (s *InvitationServiceSuite) TestListOutgoingInvitations() {
	// Create some invitations
	_, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// List outgoing invitations
	list, err := s.svc.ListOutgoing(s.org.ID, service.InvitationListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)

	// Test pagination
	list, err = s.svc.ListOutgoing(s.org.ID, service.InvitationListOptions{Page: 2, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 0)
}

func (s *InvitationServiceSuite) TestProbeOutgoingInvitations() {
	// Create some invitations
	_, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test probe outgoing invitations
	probe, err := s.svc.ProbeOutgoing(s.org.ID, service.InvitationListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
}

func (s *InvitationServiceSuite) TestAcceptInvitation() {
	// Create an invitation
	invitations, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test accept invitation
	err = s.svc.Accept(invitationID, s.users[1].GetID())
	s.Require().NoError(err)

	// Test cannot accept already accepted invitation
	err = s.svc.Accept(invitationID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotAcceptNonPendingInvitationError(
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)

	// Create another invitation
	invitations, err = s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	invitationID = invitations[0].ID
	s.Require().NoError(err)

	// Test user not allowed to accept invitation
	err = s.svc.Accept(invitationID, s.users[0].GetID())
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
	invitations, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test decline invitation
	err = s.svc.Decline(invitationID, s.users[1].GetID())
	s.Require().NoError(err)

	// Test cannot decline already declined invitation
	err = s.svc.Decline(invitationID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotDeclineNonPendingInvitationError(
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)

	// Create another invitation
	invitations, err = s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID = invitations[0].ID

	// Test user not allowed to decline invitation
	err = s.svc.Decline(invitationID, s.users[0].GetID())
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
	invitations, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test resend invitation
	err = s.svc.Resend(invitationID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDeleteInvitation() {
	// Create an invitation
	invitations, err := s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID := invitations[0].ID

	// Test delete invitation
	err = s.svc.Delete(invitationID, s.users[0].GetID())
	s.Require().NoError(err)

	// Create another invitation
	invitations, err = s.svc.Create(service.InvitationCreateOptions{
		OrganizationID: s.org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitationID = invitations[0].ID

	// Test user not allowed to delete invitation
	err = s.svc.Delete(invitationID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToDeleteInvitationError(
			s.users[1],
			repo.NewInvitationRepo().FindOrNil(invitationID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) createUsers() ([]model.User, error) {
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

func (s *InvitationServiceSuite) createOrganization(userID string) (*service.Organization, error) {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{Name: "organization"}, userID)
	if err != nil {
		return nil, err
	}
	return org, nil
}
