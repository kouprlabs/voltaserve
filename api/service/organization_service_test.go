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

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test"
)

type OrganizationServiceSuite struct {
	suite.Suite
	users []model.User
}

func TestOrganizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationServiceSuite))
}

func (s *OrganizationServiceSuite) SetupTest() {
	var err error
	s.users, err = test.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
}

func (s *OrganizationServiceSuite) TestCreate() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("organization", org.Name)
	s.Equal(model.PermissionOwner, org.Permission)
}

func (s *OrganizationServiceSuite) TestFind() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	found, err := service.NewOrganizationService().Find(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(org.ID, found.ID)
	s.Equal(org.Name, found.Name)
}

func (s *OrganizationServiceSuite) TestFind_NonExistentOrganization() {
	_, err := service.NewOrganizationService().Find(helper.NewID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestFind_UnauthorizedUser() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewOrganizationService().Find(org.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(nil).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestList() {
	for _, name := range []string{"organization A", "organization B", "organization C"} {
		_, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
			Name: name,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewOrganizationService().List(service.OrganizationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(3), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(1), list.TotalPages)
	s.Equal("organization A", list.Data[0].Name)
	s.Equal("organization B", list.Data[1].Name)
	s.Equal("organization C", list.Data[2].Name)
}

func (s *OrganizationServiceSuite) TestList_Paginate() {
	for _, name := range []string{"organization A", "organization B", "organization C"} {
		_, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
			Name: name,
		}, s.users[0].GetID())
		s.Require().NoError(err)
		time.Sleep(1 * time.Second)
	}

	list, err := service.NewOrganizationService().List(service.OrganizationListOptions{
		Page: 1,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(1), list.Page)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("organization A", list.Data[0].Name)
	s.Equal("organization B", list.Data[1].Name)

	list, err = service.NewOrganizationService().List(service.OrganizationListOptions{
		Page: 2,
		Size: 2,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(2), list.Page)
	s.Equal(uint64(1), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)
	s.Equal("organization C", list.Data[0].Name)
}

func (s *OrganizationServiceSuite) TestList_SortByNameDescending() {
	for _, name := range []string{"organization A", "organization B", "organization C"} {
		_, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
			Name: name,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	list, err := service.NewOrganizationService().List(service.OrganizationListOptions{
		Page:      1,
		Size:      3,
		SortBy:    service.OrganizationSortByName,
		SortOrder: service.OrganizationSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("organization C", list.Data[0].Name)
	s.Equal("organization B", list.Data[1].Name)
	s.Equal("organization A", list.Data[2].Name)
}

func (s *OrganizationServiceSuite) TestProbe() {
	for _, name := range []string{"organization A", "organization B", "organization C"} {
		_, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
			Name: name,
		}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	probe, err := service.NewOrganizationService().Probe(service.OrganizationListOptions{
		Page: 1,
		Size: 10,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal(uint64(3), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *OrganizationServiceSuite) TestPatchName() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	org, err = service.NewOrganizationService().PatchName(org.ID, "organization (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.Equal("organization (edit)", org.Name)
}

func (s *OrganizationServiceSuite) TestPatchName_NonExistentOrganization() {
	_, err := service.NewOrganizationService().PatchName(helper.NewID(), "organization (edit)", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestPatchName_UnauthorizedUser() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	org, err = service.NewOrganizationService().PatchName(org.ID, "organization (edit)", s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
	s.Nil(org)
}

func (s *OrganizationServiceSuite) TestDelete() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewOrganizationService().Delete(org.ID, s.users[0].GetID())
	s.Require().NoError(err)

	_, err = service.NewOrganizationService().Find(org.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestDelete_NonExistentOrganization() {
	err := service.NewOrganizationService().Delete(helper.NewID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestDelete_UnauthorizedUser() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewOrganizationService().Delete(org.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestRemoveMember() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = repo.NewOrganizationRepo().GrantUserPermission(org.ID, s.users[1].GetID(), model.PermissionEditor)
	s.Require().NoError(err)

	err = service.NewOrganizationService().RemoveMember(org.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
	s.False(guard.NewOrganizationGuard().IsAuthorized(s.users[1].GetID(), cache.NewOrganizationCache().GetOrNil(org.ID), model.PermissionEditor))
}

func (s *OrganizationServiceSuite) TestRemoveMember_NonExistentMember() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewOrganizationService().RemoveMember(org.ID, helper.NewID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotFoundError(err).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestRemoveMember_LastOwnerOfOrganization() {
	org, err := service.NewOrganizationService().Create(service.OrganizationCreateOptions{
		Name: "organization",
	}, s.users[0].GetID())
	s.Require().NoError(err)

	err = service.NewOrganizationService().RemoveMember(org.ID, s.users[0].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewCannotRemoveSoleOwnerOfOrganizationError(cache.NewOrganizationCache().GetOrNil(org.ID)).Error(), err.Error())
}
