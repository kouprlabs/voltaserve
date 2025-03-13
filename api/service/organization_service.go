// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"errors"
	"sort"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/guard"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/mapper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
	"github.com/kouprlabs/voltaserve/shared/search"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type OrganizationService struct {
	orgRepo        *repo.OrganizationRepo
	orgCache       *cache.OrganizationCache
	orgGuard       *guard.OrganizationGuard
	orgMapper      *mapper.OrganizationMapper
	orgSearch      *search.OrganizationSearch
	userSearch     *search.UserSearch
	userMapper     *userMapper
	userRepo       *repo.UserRepo
	groupCache     *cache.GroupCache
	groupRepo      *repo.GroupRepo
	groupService   *GroupService
	groupMapper    *mapper.GroupMapper
	workspaceCache *cache.WorkspaceCache
	workspaceRepo  *repo.WorkspaceRepo
	config         *config.Config
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		orgRepo: repo.NewOrganizationRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		orgCache: cache.NewOrganizationCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgGuard: guard.NewOrganizationGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgSearch: search.NewOrganizationSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		orgMapper: mapper.NewOrganizationMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		userSearch: search.NewUserSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		userRepo: repo.NewUserRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		groupCache: cache.NewGroupCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupRepo: repo.NewGroupRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		groupService: NewGroupService(),
		groupMapper: mapper.NewGroupMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		userMapper: newUserMapper(),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		config: config.GetConfig(),
	}
}

func (svc *OrganizationService) Create(opts dto.OrganizationCreateOptions, userID string) (*dto.Organization, error) {
	org, err := svc.orgRepo.Insert(repo.OrganizationInsertOptions{
		ID:   helper.NewID(),
		Name: opts.Name,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.orgRepo.GrantUserPermission(org.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	org, err = svc.orgCache.Refresh(org.GetID())
	if err != nil {
		return nil, err
	}
	if err := svc.orgSearch.Index([]model.Organization{org}); err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.MapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Find(id string, userID string) (*dto.Organization, error) {
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.MapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type OrganizationListOptions struct {
	Query     string
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
}

func (svc *OrganizationService) List(opts OrganizationListOptions, userID string) (*dto.OrganizationList, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.OrganizationSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.OrganizationSortOrderAsc
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.orgMapper.MapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &dto.OrganizationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *OrganizationService) Probe(opts OrganizationListOptions, userID string) (*dto.OrganizationProbe, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &dto.OrganizationProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *OrganizationService) PatchName(id string, name string, userID string) (*dto.Organization, error) {
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionEditor); err != nil {
		return nil, err
	}
	org.SetName(name)
	if err := svc.orgRepo.Save(org); err != nil {
		return nil, err
	}
	if err := svc.sync(org); err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.MapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Delete(id string, userID string) error {
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
		return err
	}
	if err := svc.orgRepo.Delete(id); err != nil {
		return err
	}
	if err := svc.orgCache.Delete(org.GetID()); err != nil {
		return err
	}
	if err := svc.orgSearch.Delete([]string{org.GetID()}); err != nil {
		return err
	}
	return nil
}

func (svc *OrganizationService) RemoveMember(id string, memberID string, userID string) error {
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.checkUserCanRemoveMember(memberID, org, userID); err != nil {
		return err
	}
	if err := svc.revokeGroupPermissions(memberID, org); err != nil {
		return err
	}
	if err := svc.revokeWorkspacePermissions(memberID, org); err != nil {
		return err
	}
	if err := svc.orgRepo.RevokeUserPermission(id, memberID); err != nil {
		return err
	}
	org, err = svc.orgRepo.Find(org.GetID())
	if err != nil {
		return err
	}
	if err := svc.sync(org); err != nil {
		return err
	}
	return nil
}

func (svc *OrganizationService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.OrganizationSortByName ||
		value == dto.OrganizationSortByDateCreated ||
		value == dto.OrganizationSortByDateModified
}

func (svc *OrganizationService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.OrganizationSortOrderAsc || value == dto.OrganizationSortOrderDesc
}

func (svc *OrganizationService) findAll(opts OrganizationListOptions, userID string) ([]model.Organization, error) {
	var res []model.Organization
	var err error
	if opts.Query == "" {
		res, err = svc.load(userID)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = svc.search(opts, userID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *OrganizationService) load(userID string) ([]model.Organization, error) {
	var res []model.Organization
	ids, err := svc.orgRepo.FindIDs()
	if err != nil {
		return nil, err
	}
	res, err = svc.authorizeIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) search(opts OrganizationListOptions, userID string) ([]model.Organization, error) {
	var res []model.Organization
	count, err := svc.orgRepo.Count()
	if err != nil {
		return nil, err
	}
	hits, err := svc.orgSearch.Query(opts.Query, infra.SearchQueryOptions{Limit: count})
	if err != nil {
		return nil, err
	}
	var orgs []model.Organization
	for _, hit := range hits {
		org, err := svc.orgCache.Get(hit.GetID())
		if err != nil {
			var e *errorpkg.ErrorResponse
			// We don't want to break if the search engine contains organizations that shouldn't be there
			if errors.As(err, &e) && e.Code == errorpkg.NewOrganizationNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		orgs = append(orgs, org)
	}
	res, err = svc.authorize(orgs, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) checkUserCanRemoveMember(memberID string, org model.Organization, userID string) error {
	// Ensure the member exists before proceeding
	if _, err := svc.userRepo.Find(memberID); err != nil {
		return err
	}
	// Only organization owners are allowed to remove members
	if memberID != userID {
		if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
			return err
		}
	}
	// Make sure member is not the last remaining owner of the organization
	ownerCount, err := svc.orgRepo.CountOwners(org.GetID())
	if err != nil {
		return err
	}
	if svc.orgGuard.IsAuthorized(memberID, org, model.PermissionOwner) && ownerCount == 1 {
		return errorpkg.NewCannotRemoveSoleOwnerOfOrganizationError(org)
	}
	return nil
}

func (svc *OrganizationService) revokeGroupPermissions(memberID string, org model.Organization) error {
	groupsIDs, err := svc.groupRepo.FindIDsByOrganization(org.GetID())
	if err != nil {
		return err
	}
	for _, groupID := range groupsIDs {
		if err := svc.groupRepo.RevokeUserPermission(groupID, memberID); err != nil {
			logger.GetLogger().Error(err)
		}
		if _, err := svc.groupCache.Refresh(groupID); err != nil {
			logger.GetLogger().Error(err)
		}
	}
	return nil
}

func (svc *OrganizationService) revokeWorkspacePermissions(memberID string, org model.Organization) error {
	workspaceIDs, err := svc.workspaceRepo.FindIDsByOrganization(org.GetID())
	if err != nil {
		return err
	}
	for _, workspaceID := range workspaceIDs {
		if err := svc.workspaceRepo.RevokeUserPermission(workspaceID, memberID); err != nil {
			logger.GetLogger().Error(err)
		}
		if _, err := svc.workspaceCache.Refresh(workspaceID); err != nil {
			logger.GetLogger().Error(err)
		}
	}
	return nil
}

func (svc *OrganizationService) authorize(data []model.Organization, userID string) ([]model.Organization, error) {
	var res []model.Organization
	for _, o := range data {
		if svc.orgGuard.IsAuthorized(userID, o, model.PermissionViewer) {
			res = append(res, o)
		}
	}
	return res, nil
}

func (svc *OrganizationService) authorizeIDs(ids []string, userID string) ([]model.Organization, error) {
	var res []model.Organization
	for _, id := range ids {
		var o model.Organization
		o, err := svc.orgCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewOrganizationNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if svc.orgGuard.IsAuthorized(userID, o, model.PermissionViewer) {
			res = append(res, o)
		}
	}
	return res, nil
}

func (svc *OrganizationService) sort(data []model.Organization, sortBy string, sortOrder string) []model.Organization {
	if sortBy == dto.OrganizationSortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.OrganizationSortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == dto.OrganizationSortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a := helper.StringToTime(data[i].GetCreateTime())
			b := helper.StringToTime(data[j].GetCreateTime())
			if sortOrder == dto.OrganizationSortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == dto.OrganizationSortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a := helper.StringToTime(*data[i].GetUpdateTime())
				b := helper.StringToTime(*data[j].GetUpdateTime())
				if sortOrder == dto.OrganizationSortOrderDesc {
					return a.UnixMilli() > b.UnixMilli()
				} else {
					return a.UnixMilli() < b.UnixMilli()
				}
			} else {
				return false
			}
		})
		return data
	}
	return data
}

func (svc *OrganizationService) paginate(data []model.Organization, page, size uint64) (pageData []model.Organization, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Organization{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *OrganizationService) sync(org model.Organization) error {
	if err := svc.orgCache.Set(org); err != nil {
		return err
	}
	if err := svc.orgSearch.Update([]model.Organization{org}); err != nil {
		return err
	}
	return nil
}
