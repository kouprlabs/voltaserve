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
	"slices"
	"sort"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/guard"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
	"github.com/kouprlabs/voltaserve/shared/search"

	"github.com/kouprlabs/voltaserve/api/config"
)

type GroupService struct {
	groupRepo      *repo.GroupRepo
	groupGuard     *guard.GroupGuard
	groupSearch    *search.GroupSearch
	groupMapper    *groupMapper
	groupCache     *cache.GroupCache
	userRepo       *repo.UserRepo
	userSearch     *search.UserSearch
	userMapper     *userMapper
	workspaceRepo  *repo.WorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	fileRepo       *repo.FileRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	orgRepo        *repo.OrganizationRepo
	orgCache       *cache.OrganizationCache
	orgGuard       *guard.OrganizationGuard
	config         *config.Config
}

func NewGroupService() *GroupService {
	return &GroupService{
		groupRepo: repo.NewGroupRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		groupGuard: guard.NewGroupGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupCache: cache.NewGroupCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupSearch: search.NewGroupSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		groupMapper: newGroupMapper(),
		userRepo: repo.NewUserRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		userSearch: search.NewUserSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		userMapper: newUserMapper(),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgRepo: repo.NewOrganizationRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		orgGuard: guard.NewOrganizationGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgCache: cache.NewOrganizationCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		config: config.GetConfig(),
	}
}

func (svc *GroupService) Create(opts dto.GroupCreateOptions, userID string) (*dto.Group, error) {
	org, err := svc.orgCache.Get(opts.OrganizationID)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionEditor); err != nil {
		return nil, err
	}
	group, err := svc.groupRepo.Insert(repo.GroupInsertOptions{
		ID:             helper.NewID(),
		Name:           opts.Name,
		OrganizationID: opts.OrganizationID,
		OwnerID:        userID,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.groupRepo.GrantUserPermission(group.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	group, err = svc.groupCache.Refresh(group.GetID())
	if err != nil {
		return nil, err
	}
	if err := svc.groupSearch.Index([]model.Group{group}); err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapOne(group, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) Find(id string, userID string) (*dto.Group, error) {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapOne(group, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type GroupListOptions struct {
	Query          string
	OrganizationID string
	Page           uint64
	Size           uint64
	SortBy         string
	SortOrder      string
}

func (svc *GroupService) List(opts GroupListOptions, userID string) (*dto.GroupList, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.GroupSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.GroupSortOrderAsc
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.groupMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &dto.GroupList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *GroupService) Probe(opts GroupListOptions, userID string) (*dto.GroupProbe, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &dto.GroupProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *GroupService) PatchName(id string, name string, userID string) (*dto.Group, error) {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionEditor); err != nil {
		return nil, err
	}
	group.SetName(name)
	if err := svc.groupRepo.Save(group); err != nil {
		return nil, err
	}
	if err := svc.sync(group); err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapOne(group, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) Delete(id string, userID string) error {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionOwner); err != nil {
		return err
	}
	if err := svc.groupRepo.Delete(id); err != nil {
		return err
	}
	if err := svc.groupSearch.Delete([]string{group.GetID()}); err != nil {
		return err
	}
	if err := svc.groupCache.Delete(group.GetID()); err != nil {
		return err
	}
	return nil
}

func (svc *GroupService) AddMember(id string, memberID string, userID string) error {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionOwner); err != nil {
		return err
	}
	if err := svc.checkUserIsMemberOfOrganization(memberID, group.GetOrganizationID()); err != nil {
		return err
	}
	// Ensure that the member doesn't already have a higher permission on the group,
	// if we don't check that, we risk downgrading the existing permission
	if !svc.groupGuard.IsAuthorized(memberID, group, model.PermissionViewer) &&
		!svc.groupGuard.IsAuthorized(memberID, group, model.PermissionEditor) {
		if err := svc.groupRepo.GrantUserPermission(group.GetID(), memberID, model.PermissionViewer); err != nil {
			return err
		}
		if _, err := svc.groupCache.Refresh(group.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *GroupService) RemoveMember(id string, memberID string, userID string) error {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionOwner); err != nil {
		return err
	}
	if err := svc.checkUserIsMemberOfOrganization(memberID, group.GetOrganizationID()); err != nil {
		return err
	}
	// Make sure member is not the last remaining owner of the group
	ownerCount, err := svc.groupRepo.CountOwners(group.GetID())
	if err != nil {
		return err
	}
	if svc.groupGuard.IsAuthorized(memberID, group, model.PermissionOwner) && ownerCount == 1 {
		return errorpkg.NewCannotRemoveSoleOwnerOfGroupError(group)
	}
	if err := svc.groupRepo.RevokeUserPermission(id, memberID); err != nil {
		return err
	}
	if _, err := svc.groupCache.Refresh(group.GetID()); err != nil {
		return err
	}
	return nil
}

func (svc *GroupService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.GroupSortByName ||
		value == dto.GroupSortByDateCreated ||
		value == dto.GroupSortByDateModified
}

func (svc *GroupService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.GroupSortOrderAsc || value == dto.GroupSortOrderDesc
}

func (svc *GroupService) checkUserIsMemberOfOrganization(userID string, organizationID string) error {
	org, err := svc.orgCache.Get(organizationID)
	if err != nil {
		return err
	}
	if !slices.Contains(org.GetMembers(), userID) {
		return errorpkg.NewUserNotMemberOfOrganizationError()
	}
	return nil
}

func (svc *GroupService) findAll(opts GroupListOptions, userID string) ([]model.Group, error) {
	var res []model.Group
	var err error
	if opts.Query == "" {
		res, err = svc.load(opts, userID)
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

func (svc *GroupService) load(opts GroupListOptions, userID string) ([]model.Group, error) {
	var res []model.Group
	if opts.OrganizationID == "" {
		ids, err := svc.groupRepo.FindIDs()
		if err != nil {
			return nil, err
		}
		res, err = svc.authorizeIDs(ids, userID)
		if err != nil {
			return nil, err
		}
	} else {
		groups, err := svc.orgRepo.FindGroups(opts.OrganizationID)
		if err != nil {
			return nil, err
		}
		res, err = svc.authorize(groups, userID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *GroupService) search(opts GroupListOptions, userID string) ([]model.Group, error) {
	var res []model.Group
	count, err := svc.groupRepo.Count()
	if err != nil {
		return nil, err
	}
	hits, err := svc.groupSearch.Query(opts.Query, infra.SearchQueryOptions{Limit: count})
	if err != nil {
		return nil, err
	}
	var groups []model.Group
	for _, hit := range hits {
		group, err := svc.groupCache.Get(hit.GetID())
		if err != nil {
			var e *errorpkg.ErrorResponse
			// We don't want to break if the search engine contains groups that shouldn't be there
			if errors.As(err, &e) && e.Code == errorpkg.NewGroupNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		groups = append(groups, group)
	}
	var filtered []model.Group
	if opts.OrganizationID == "" {
		filtered = groups
	} else {
		for _, g := range groups {
			if g.GetOrganizationID() == opts.OrganizationID {
				filtered = append(filtered, g)
			}
		}
	}
	res, err = svc.authorize(filtered, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) authorize(data []model.Group, userID string) ([]model.Group, error) {
	var res []model.Group
	for _, g := range data {
		if svc.groupGuard.IsAuthorized(userID, g, model.PermissionViewer) {
			res = append(res, g)
		}
	}
	return res, nil
}

func (svc *GroupService) authorizeIDs(ids []string, userID string) ([]model.Group, error) {
	var res []model.Group
	for _, id := range ids {
		var o model.Group
		o, err := svc.groupCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewGroupNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if svc.groupGuard.IsAuthorized(userID, o, model.PermissionViewer) {
			res = append(res, o)
		}
	}
	return res, nil
}

func (svc *GroupService) sort(data []model.Group, sortBy string, sortOrder string) []model.Group {
	if sortBy == dto.GroupSortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.GroupSortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == dto.GroupSortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a := helper.StringToTime(data[i].GetCreateTime())
			b := helper.StringToTime(data[j].GetCreateTime())
			if sortOrder == dto.GroupSortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == dto.GroupSortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a := helper.StringToTime(*data[i].GetUpdateTime())
				b := helper.StringToTime(*data[j].GetUpdateTime())
				if sortOrder == dto.GroupSortOrderDesc {
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

func (svc *GroupService) paginate(data []model.Group, page, size uint64) (pageData []model.Group, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Group{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *GroupService) sync(group model.Group) error {
	if err := svc.groupCache.Set(group); err != nil {
		return err
	}
	if err := svc.groupSearch.Update([]model.Group{group}); err != nil {
		return err
	}
	return nil
}

type groupMapper struct {
	orgCache   *cache.OrganizationCache
	orgMapper  *organizationMapper
	groupCache *cache.GroupCache
}

func newGroupMapper() *groupMapper {
	return &groupMapper{
		orgCache: cache.NewOrganizationCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		orgMapper: newOrganizationMapper(),
		groupCache: cache.NewGroupCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

func (mp *groupMapper) mapOne(m model.Group, userID string) (*dto.Group, error) {
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	o, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	res := &dto.Group{
		ID:           m.GetID(),
		Name:         m.GetName(),
		Organization: *o,
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}
	res.Permission = model.PermissionNone
	for _, p := range m.GetUserPermissions() {
		if p.GetUserID() == userID && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
			res.Permission = p.GetValue()
		}
	}
	for _, p := range m.GetGroupPermissions() {
		g, err := mp.groupCache.Get(p.GetGroupID())
		if err != nil {
			return nil, err
		}
		for _, u := range g.GetMembers() {
			if u == userID && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
				res.Permission = p.GetValue()
			}
		}
	}
	return res, nil
}

func (mp *groupMapper) mapMany(groups []model.Group, userID string) ([]*dto.Group, error) {
	res := make([]*dto.Group, 0)
	for _, group := range groups {
		g, err := mp.mapOne(group, userID)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewGroupNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, g)
	}
	return res, nil
}
