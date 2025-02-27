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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
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
	s.users, err = test.CreateUsers(4)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *InvitationServiceSuite) TestCreate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test-a@voltaserve.com", "test-b@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Len(invitations, 2)
}

func (s *InvitationServiceSuite) TestCreate_MissingOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForOrganization(org, s.users[0])

	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test-a@voltaserve.com", "test-b@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestCreate_DuplicateEmails() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Empty(invitations)
}

func (s *InvitationServiceSuite) TestCreate_NonExistentOrganization() {
	_, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: helper.NewID(),
		Emails:         []string{"test@voltaserve.com"},
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestListIncoming() {
	for _, userID := range []string{s.users[1].GetID(), s.users[2].GetID(), s.users[3].GetID()} {
		org, err := test.CreateOrganization(userID)
		s.Require().NoError(err)
		_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
			OrganizationID: org.ID,
			Emails:         []string{s.users[0].GetEmail()},
		}, userID)
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewInvitationService().ListIncoming(dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(s.users[1].GetEmail(), list.Data[0].Owner.Email)
	s.Equal(s.users[2].GetEmail(), list.Data[1].Owner.Email)
	s.Equal(s.users[3].GetEmail(), list.Data[2].Owner.Email)
}

func (s *InvitationServiceSuite) TestListIncoming_SortByEmailDescending() {
	for _, userID := range []string{s.users[1].GetID(), s.users[2].GetID(), s.users[3].GetID()} {
		org, err := test.CreateOrganization(userID)
		s.Require().NoError(err)
		_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
			OrganizationID: org.ID,
			Emails:         []string{s.users[0].GetEmail()},
		}, userID)
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewInvitationService().ListIncoming(dto.InvitationListOptions{
		Page:      1,
		Size:      3,
		SortBy:    dto.InvitationSortByEmail,
		SortOrder: dto.InvitationSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(s.users[3].GetEmail(), list.Data[0].Owner.Email)
	s.Equal(s.users[2].GetEmail(), list.Data[1].Owner.Email)
	s.Equal(s.users[1].GetEmail(), list.Data[2].Owner.Email)
}

func (s *InvitationServiceSuite) TestListIncoming_Paginate() {
	for _, userID := range []string{s.users[1].GetID(), s.users[2].GetID(), s.users[3].GetID()} {
		org, err := test.CreateOrganization(userID)
		s.Require().NoError(err)
		_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
			OrganizationID: org.ID,
			Emails:         []string{s.users[0].GetEmail()},
		}, userID)
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewInvitationService().ListIncoming(dto.InvitationListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(s.users[1].GetEmail(), list.Data[0].Owner.Email)
	s.Equal(s.users[2].GetEmail(), list.Data[1].Owner.Email)

	list, err = service.NewInvitationService().ListIncoming(dto.InvitationListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(s.users[3].GetEmail(), list.Data[0].Owner.Email)
}

func (s *InvitationServiceSuite) TestProbeIncoming() {
	for _, userID := range []string{s.users[1].GetID(), s.users[2].GetID(), s.users[3].GetID()} {
		org, err := test.CreateOrganization(userID)
		s.Require().NoError(err)
		_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
			OrganizationID: org.ID,
			Emails:         []string{s.users[0].GetEmail()},
		}, userID)
		s.Require().NoError(err)
	}

	probe, err := service.NewInvitationService().ProbeIncoming(dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *InvitationServiceSuite) TestCountIncoming() {
	for _, userID := range []string{s.users[1].GetID(), s.users[2].GetID(), s.users[3].GetID()} {
		org, err := test.CreateOrganization(userID)
		s.Require().NoError(err)
		_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
			OrganizationID: org.ID,
			Emails:         []string{s.users[0].GetEmail()},
		}, userID)
		s.Require().NoError(err)
	}

	count, err := service.NewInvitationService().GetCountIncoming(s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(int64(3), *count)
}

func (s *InvitationServiceSuite) TestListOutgoing() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewInvitationService().ListOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal(strings.ToLower(s.users[1].GetEmail()), list.Data[0].Email)
	s.Equal(strings.ToLower(s.users[2].GetEmail()), list.Data[1].Email)
	s.Equal(strings.ToLower(s.users[3].GetEmail()), list.Data[2].Email)
}

func (s *InvitationServiceSuite) TestListOutgoing_MissingOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForOrganization(org, s.users[0])

	_, err = service.NewInvitationService().ListOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestListOutgoing_InsufficientOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForOrganization(org, s.users[0], model.PermissionViewer)

	_, err = service.NewInvitationService().ListOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewOrganizationPermissionError(
			s.users[0].GetID(),
			cache.NewOrganizationCache().GetOrNil(org.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestListOutgoing_SortByEmailDescending() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewInvitationService().ListOutgoing(org.ID, dto.InvitationListOptions{
		Page:      1,
		Size:      3,
		SortBy:    dto.InvitationSortByEmail,
		SortOrder: dto.InvitationSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(strings.ToLower(s.users[3].GetEmail()), list.Data[0].Email)
	s.Equal(strings.ToLower(s.users[2].GetEmail()), list.Data[1].Email)
	s.Equal(strings.ToLower(s.users[1].GetEmail()), list.Data[2].Email)
}

func (s *InvitationServiceSuite) TestListOutgoing_Paginate() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	list, err := service.NewInvitationService().ListOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(strings.ToLower(s.users[1].GetEmail()), list.Data[0].Email)
	s.Equal(strings.ToLower(s.users[2].GetEmail()), list.Data[1].Email)

	list, err = service.NewInvitationService().ListOutgoing(org.ID, dto.InvitationListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal(strings.ToLower(s.users[3].GetEmail()), list.Data[0].Email)
}

func (s *InvitationServiceSuite) TestProbeOutgoing() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	probe, err := service.NewInvitationService().ProbeOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *InvitationServiceSuite) TestProbeOutgoing_MissingOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.revokeUserPermissionForOrganization(org, s.users[0])

	_, err = service.NewInvitationService().ProbeOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *InvitationServiceSuite) TestProbeOutgoing_InsufficientOrganizationPermission() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	_, err = service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail(), s.users[2].GetEmail(), s.users[3].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	s.grantUserPermissionForOrganization(org, s.users[0], model.PermissionViewer)

	_, err = service.NewInvitationService().ProbeOutgoing(org.ID, dto.InvitationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(
		errorpkg.NewOrganizationPermissionError(
			s.users[0].GetID(),
			cache.NewOrganizationCache().GetOrNil(org.ID),
			model.PermissionOwner,
		).Error(),
		err.Error(),
	)
}

func (s *InvitationServiceSuite) TestAccept() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Accept(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestAccept_AlreadyAccepted() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Decline(invitations[0].ID, s.users[1].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDecline_AlreadyDeclined() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
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
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Resend(invitations[0].ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDelete() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
		OrganizationID: org.ID,
		Emails:         []string{s.users[1].GetEmail()},
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewInvitationService().Delete(invitations[0].ID, s.users[0].GetID())
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) TestDelete_UnauthorizedUser() {
	org, err := test.CreateOrganization(s.users[0].GetID())
	s.Require().NoError(err)
	invitations, err := service.NewInvitationService().Create(dto.InvitationCreateOptions{
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

func (s *InvitationServiceSuite) grantUserPermissionForOrganization(org *dto.Organization, user model.User, permission string) {
	err := repo.NewOrganizationRepo().GrantUserPermission(org.ID, user.GetID(), permission)
	s.Require().NoError(err)
	_, err = cache.NewOrganizationCache().Refresh(org.ID)
	s.Require().NoError(err)
}

func (s *InvitationServiceSuite) revokeUserPermissionForOrganization(org *dto.Organization, user model.User) {
	err := repo.NewOrganizationRepo().RevokeUserPermission(org.ID, user.GetID())
	s.Require().NoError(err)
	_, err = cache.NewOrganizationCache().Refresh(org.ID)
	s.Require().NoError(err)
}
