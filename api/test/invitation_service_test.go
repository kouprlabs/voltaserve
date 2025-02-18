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
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type InvitationServiceSuite struct {
	suite.Suite
	users []model.User
}

func TestInvitationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InvitationServiceSuite))
}

func (s *InvitationServiceSuite) SetupTest() {
	var err error
	s.users, err = test_helper.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *InvitationServiceSuite) TestCreate() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test-a@voltaserve.com", "test-b@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(invitations, 2)
}

func (s *InvitationServiceSuite) TestCreate_DuplicateEmails() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(invitations)
}

func (s *InvitationServiceSuite) TestCreate_NonExistentOrganization() {
	_, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: helper.NewID(),
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestCreate_UnauthorizedUser() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestListIncoming() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewInvitationService().ListIncoming(service.InvitationListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.InvitationSortByEmail,
		SortOrder: service.InvitationSortOrderAsc,
	}, s.users[1].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)

	list, err = service.NewInvitationService().ListIncoming(service.InvitationListOptions{
		Page:      2,
		Size:      10,
		SortBy:    service.InvitationSortByEmail,
		SortOrder: service.InvitationSortOrderAsc,
	}, s.users[1].GetID())
	s.Require().NoError(err)
	s.Empty(list.Data)
}

func (s *InvitationServiceSuite) TestProbeIncoming() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	probe, err := service.NewInvitationService().ProbeIncoming(service.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[1].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
}

func (s *InvitationServiceSuite) TestCountIncoming() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	count, err := service.NewInvitationService().CountIncoming(s.users[1].GetID())
	s.Require().NoError(err)
	s.Equal(int64(1), *count)
}

func (s *InvitationServiceSuite) TestListOutgoing() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewInvitationService().ListOutgoing(org.ID, service.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(list.Data, 1)

	list, err = service.NewInvitationService().ListOutgoing(org.ID, service.InvitationListOptions{
		Page: 2,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(list.Data)
}

func (s *InvitationServiceSuite) TestProbeOutgoing() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	probe, err := service.NewInvitationService().ProbeOutgoing(org.ID, service.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), probe.TotalElements)
}

func (s *InvitationServiceSuite) TestAccept() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestAccept_AlreadyAccepted() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotAcceptNonPendingInvitationError(
			repo.NewInvitationRepo().FindOrNil(invitations[0].ID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestAccept_UnauthorizedUser() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToAcceptInvitationError(
			s.users[0],
			repo.NewInvitationRepo().FindOrNil(invitations[0].ID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestDecline() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Decline(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDecline_AlreadyDeclined() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Decline(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Decline(invitations[0].ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewCannotDeclineNonPendingInvitationError(
			repo.NewInvitationRepo().FindOrNil(invitations[0].ID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestDecline_UnauthorizedUser() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Decline(invitations[0].ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToDeclineInvitationError(
			s.users[0],
			repo.NewInvitationRepo().FindOrNil(invitations[0].ID),
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestResend() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Resend(invitations[0].ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDelete() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Delete(invitations[0].ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDelete_UnauthorizedUser() {
	org, err := test_helper.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(service.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Delete(invitations[0].ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewUserNotAllowedToDeleteInvitationError(
			s.users[1],
			repo.NewInvitationRepo().FindOrNil(invitations[0].ID),
		).Error(),
		err.Error(),
	)
}
