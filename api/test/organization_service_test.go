package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/service"
)

type OrganizationServiceSuite struct {
	suite.Suite
	orgSvc   *service.OrganizationService
	orgRepo  *repo.OrganizationRepo
	orgCache *cache.OrganizationCache
	orgGuard *guard.OrganizationGuard
	userIDs  []string
}

func TestOrganizationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationServiceSuite))
}

func (s *OrganizationServiceSuite) SetupTest() {
	s.orgSvc = service.NewOrganizationService()
	s.orgRepo = repo.NewOrganizationRepo()
	s.orgCache = cache.NewOrganizationCache()
	userIDs, err := s.createUsers()
	if err != nil {
		s.Fail(err.Error())
		return
	}
	s.userIDs = userIDs
}

func (s *OrganizationServiceSuite) TestCreateOrganization() {
	opts := service.OrganizationCreateOptions{Name: "organization"}
	org, err := s.orgSvc.Create(opts, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().NotNil(org)
	s.Equal(opts.Name, org.Name)
	s.Equal(model.PermissionOwner, org.Permission)
}

func (s *OrganizationServiceSuite) TestFindOrganization() {
	// Create a new organization to find
	opts := service.OrganizationCreateOptions{Name: "organization"}
	createdOrg, err := s.orgSvc.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test finding the organization
	foundOrg, err := s.orgSvc.Find(createdOrg.ID, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().NotNil(foundOrg)
	s.Equal(createdOrg.ID, foundOrg.ID)
	s.Equal(createdOrg.Name, foundOrg.Name)

	// Test finding a non-existent organization
	foundOrg, err = s.orgSvc.Find("non-existent-org-id", s.userIDs[0])
	s.Require().Error(err)
	s.Nil(foundOrg)

	// Test finding an organization with insufficient permissions
	foundOrg, err = s.orgSvc.Find(createdOrg.ID, s.userIDs[1])
	s.Require().Error(err)
	s.Nil(foundOrg)
}

func (s *OrganizationServiceSuite) TestListOrganizations() {
	// Create multiple organizations
	names := []string{"organization A", "organization B", "organization C"}
	for _, name := range names {
		opts := service.OrganizationCreateOptions{
			Name:  name,
			Image: nil,
		}
		_, err := s.orgSvc.Create(opts, s.userIDs[0])
		s.Require().NoError(err)
	}

	// Test listing all organizations
	list, err := s.orgSvc.List(service.OrganizationListOptions{
		Page:      1,
		Size:      10,
		SortBy:    service.OrganizationSortByName,
		SortOrder: service.OrganizationSortOrderAsc,
	}, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().NotNil(list)
	s.Equal(uint64(len(names)), list.TotalElements)
	s.Equal(names[0], list.Data[0].Name)
	s.Equal(names[1], list.Data[1].Name)
	s.Equal(names[2], list.Data[2].Name)

	// Test pagination
	list, err = s.orgSvc.List(service.OrganizationListOptions{Page: 1, Size: 2}, s.userIDs[0])
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
	}, s.userIDs[0])
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
	}, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().NotNil(list)
	s.Equal("organization A", list.Data[0].Name)
	s.Equal("organization B", list.Data[1].Name)
}

func (s *OrganizationServiceSuite) TestProbeOrganizations() {
	// Create multiple organizations
	names := []string{"organization A", "organization B", "organization C"}
	for _, name := range names {
		opts := service.OrganizationCreateOptions{
			Name:  name,
			Image: nil,
		}
		_, err := s.orgSvc.Create(opts, s.userIDs[0])
		s.Require().NoError(err)
	}

	// Test probing organizations
	probe, err := s.orgSvc.Probe(service.OrganizationListOptions{Page: 1, Size: 10}, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().NotNil(probe)
	s.Equal(uint64(len(names)), probe.TotalElements)
	s.Equal(uint64(1), probe.TotalPages)
}

func (s *OrganizationServiceSuite) TestPatchOrganizationName() {
	// Create a new organization
	createdOrg, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.userIDs[0])
	s.Require().NoError(err)

	// Test patching the organization name
	newName := "organization (edit)"
	updatedOrg, err := s.orgSvc.PatchName(createdOrg.ID, newName, s.userIDs[0])
	s.Require().NoError(err)
	s.Require().NotNil(updatedOrg)
	s.Equal(newName, updatedOrg.Name)

	// Test patching with insufficient permissions
	updatedOrg, err = s.orgSvc.PatchName(createdOrg.ID, newName, s.userIDs[1])
	s.Require().Error(err)
	s.Nil(updatedOrg)

	// Test patching a non-existent organization
	updatedOrg, err = s.orgSvc.PatchName("non-existent-org-id", newName, s.userIDs[0])
	s.Require().Error(err)
	s.Nil(updatedOrg)
}

func (s *OrganizationServiceSuite) TestDeleteOrganization() {
	// Create a new organization
	opts := service.OrganizationCreateOptions{Name: "organization"}
	createdOrg, err := s.orgSvc.Create(opts, s.userIDs[0])
	s.Require().NoError(err)

	// Test deleting the organization
	err = s.orgSvc.Delete(createdOrg.ID, s.userIDs[0])
	s.Require().NoError(err)

	// Verify the organization is deleted
	foundOrg, err := s.orgSvc.Find(createdOrg.ID, s.userIDs[0])
	s.Require().Error(err)
	s.Nil(foundOrg)

	// Test deleting a non-existent organization
	err = s.orgSvc.Delete("non-existent-org-id", s.userIDs[0])
	s.Require().Error(err)

	// Test deleting with insufficient permissions
	createdOrg, err = s.orgSvc.Create(opts, s.userIDs[0])
	s.Require().NoError(err)
	err = s.orgSvc.Delete(createdOrg.ID, s.userIDs[1])
	s.Require().Error(err)
}

func (s *OrganizationServiceSuite) TestRemoveMember() {
	// Create a new organization
	createdOrg, err := s.orgSvc.Create(service.OrganizationCreateOptions{Name: "organization"}, s.userIDs[0])
	s.Require().NoError(err)

	// Add another user to the organization
	err = s.orgRepo.GrantUserPermission(createdOrg.ID, s.userIDs[1], model.PermissionEditor)
	s.Require().NoError(err)

	// Test removing the member
	err = s.orgSvc.RemoveMember(createdOrg.ID, s.userIDs[1], s.userIDs[0])
	s.Require().NoError(err)

	// Verify the member is removed
	org, err := s.orgCache.Get(createdOrg.ID)
	s.Require().NoError(err)
	s.False(s.orgGuard.IsAuthorized(s.userIDs[1], org, model.PermissionEditor))

	// Test removing a non-existent member
	err = s.orgSvc.RemoveMember(createdOrg.ID, "non-existent-user-id", s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewUserNotFoundError(err).Error(), err.Error())

	// Test removing the last owner
	err = s.orgSvc.RemoveMember(createdOrg.ID, s.userIDs[0], s.userIDs[0])
	s.Require().Error(err)
	s.Equal(errorpkg.NewCannotRemoveSoleOwnerOfOrganizationError(org).Error(), err.Error())
}

func (s *OrganizationServiceSuite) createUsers() ([]string, error) {
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
	return ids, nil
}
