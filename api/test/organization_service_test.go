package test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
	"github.com/kouprlabs/voltaserve/api/test/test_helper"
)

type OrganizationServiceSuite struct {
	suite.Suite
	orgSvc   *service.OrganizationService
	orgRepo  *repo.OrganizationRepo
	orgCache *cache.OrganizationCache
	orgGuard *guard.OrganizationGuard
	users    []model.User
}

func TestOrganizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationServiceSuite))
}

func (s *OrganizationServiceSuite) SetupTest() {
	s.orgSvc = service.NewOrganizationService()
	s.orgRepo = repo.NewOrganizationRepo()
	s.orgCache = cache.NewOrganizationCache()
	users, err := test_helper.CreateUsers(2)
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.users = users
}

func (s *OrganizationServiceSuite) TestCreateOrganization() {
	opts := service.OrganizationCreateOptions{Name: "organization"}
	org, err := s.orgSvc.Create(opts, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(org)
	s.Equal(opts.Name, org.Name)
	s.Equal(model.PermissionOwner, org.Permission)
}

func (s *OrganizationServiceSuite) TestFindOrganization() {
	// Create a new organization to find
	org, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test finding the organization
	foundOrg, err := s.orgSvc.Find(org.ID, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(foundOrg)
	s.Equal(org.ID, foundOrg.ID)
	s.Equal(org.Name, foundOrg.Name)

	// Test finding a non-existent organization
	foundOrg, err = s.orgSvc.Find("non-existent-org-id", s.users[0].GetID())
	s.Require().Error(err)
	s.Nil(foundOrg)

	// Test finding an organization with insufficient permissions
	foundOrg, err = s.orgSvc.Find(org.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Nil(foundOrg)
}

func (s *OrganizationServiceSuite) TestListOrganizations() {
	// Create multiple organizations
	names := []string{"organization A", "organization B", "organization C"}
	for _, name := range names {
		_, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: name}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	// Test listing all organizations
	list, err := s.orgSvc.List(service.OrganizationListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.OrganizationSortByName,
		SortOrder: service.OrganizationSortOrderAsc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(list)
	s.Equal(uint64(len(names)), list.TotalElements)
	s.Equal(names[0], list.Data[0].Name)
	s.Equal(names[1], list.Data[1].Name)
	s.Equal(names[2], list.Data[2].Name)

	// Test pagination
	list, err = s.orgSvc.List(service.OrganizationListOptions{Page: 1, Size: 2}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(list)
	s.Equal(uint64(2), list.Size)
	s.Equal(uint64(3), list.TotalElements)
	s.Equal(uint64(2), list.TotalPages)

	// Test sorting by name in descending order
	list, err = s.orgSvc.List(service.OrganizationListOptions{
		Page:      1,
		Size:      2,
		SortBy:    service.OrganizationSortByName,
		SortOrder: service.OrganizationSortOrderDesc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(list)
	s.Equal("organization C", list.Data[0].Name)
	s.Equal("organization B", list.Data[1].Name)

	// Test sorting by date created
	list, err = s.orgSvc.List(service.OrganizationListOptions{
		Page:      1,
		Size:      2,
		SortBy:    service.OrganizationSortByName,
		SortOrder: service.OrganizationSortOrderAsc,
	}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(list)
	s.Equal("organization A", list.Data[0].Name)
	s.Equal("organization B", list.Data[1].Name)
}

func (s *OrganizationServiceSuite) TestProbeOrganizations() {
	// Create multiple organizations
	names := []string{"organization A", "organization B", "organization C"}
	for _, name := range names {
		_, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: name}, s.users[0].GetID())
		s.Require().NoError(err)
	}

	// Test probing organizations
	probe, err := s.orgSvc.Probe(service.OrganizationListOptions{Page: 1, Size: 10}, s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(probe)
	s.Equal(uint64(len(names)), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *OrganizationServiceSuite) TestPatchOrganizationName() {
	// Create a new organization
	org, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test patching the organization name
	org, err = s.orgSvc.PatchName(org.ID, "organization (edit)", s.users[0].GetID())
	s.Require().NoError(err)
	s.Require().NotNil(org)
	s.Equal("organization (edit)", org.Name)

	// Test patching with insufficient permissions
	org, err = s.orgSvc.PatchName(org.ID, "organization", s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
	s.Nil(org)

	// Test patching a non-existent organization
	org, err = s.orgSvc.PatchName("non-existent-org-id", "organization", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
	s.Nil(org)
}

func (s *OrganizationServiceSuite) TestDeleteOrganization() {
	// Create a new organization
	org, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.users[0].GetID())
	s.Require().NoError(err)

	// Test deleting the organization
	err = s.orgSvc.Delete(org.ID, s.users[0].GetID())
	s.Require().NoError(err)

	// Verify the organization is deleted
	org, err = s.orgSvc.Find(org.ID, s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
	s.Nil(org)

	// Test deleting a non-existent organization
	err = s.orgSvc.Delete("non-existent-org-id", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())

	// Test deleting with insufficient permissions
	org, err = s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.users[0].GetID())
	s.Require().NoError(err)
	err = s.orgSvc.Delete(org.ID, s.users[1].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewOrganizationNotFoundError(err).Error(), err.Error())
}

func (s *OrganizationServiceSuite) TestRemoveMember() {
	// Create a new organization
	org, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.users[0].GetID())
	s.Require().NoError(err)

	// Add another user to the organization
	err = s.orgRepo.GrantUserPermission(org.ID, s.users[1].GetID(), model.PermissionEditor)
	s.Require().NoError(err)

	// Test removing the member
	err = s.orgSvc.RemoveMember(org.ID, s.users[1].GetID(), s.users[0].GetID())
	s.Require().NoError(err)
	s.False(s.orgGuard.IsAuthorized(s.users[1].GetID(), s.orgCache.GetOrNil(org.ID), model.PermissionEditor))

	// Test removing a non-existent member
	err = s.orgSvc.RemoveMember(org.ID, "non-existent-user-id", s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotFoundError(err).Error(), err.Error())

	// Test removing the last owner
	err = s.orgSvc.RemoveMember(org.ID, s.users[0].GetID(), s.users[0].GetID())
	s.Require().Error(err)
	s.Equal(errorpkg.NewCannotRemoveSoleOwnerOfOrganizationError(s.orgCache.GetOrNil(org.ID)).Error(), err.Error())
}
