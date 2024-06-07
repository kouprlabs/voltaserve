package service

import (
	"sort"
	"time"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/log"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type OrganizationService struct {
	orgRepo      repo.OrganizationRepo
	orgCache     *cache.OrganizationCache
	orgGuard     *guard.OrganizationGuard
	orgMapper    *organizationMapper
	orgSearch    *search.OrganizationSearch
	userSearch   *search.UserSearch
	userMapper   *userMapper
	groupRepo    repo.GroupRepo
	groupService *GroupService
	groupMapper  *groupMapper
	config       config.Config
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{
		orgRepo:      repo.NewOrganizationRepo(),
		orgCache:     cache.NewOrganizationCache(),
		orgGuard:     guard.NewOrganizationGuard(),
		orgSearch:    search.NewOrganizationSearch(),
		orgMapper:    newOrganizationMapper(),
		userSearch:   search.NewUserSearch(),
		groupRepo:    repo.NewGroupRepo(),
		groupService: NewGroupService(),
		groupMapper:  newGroupMapper(),
		userMapper:   newUserMapper(),
		config:       config.GetConfig(),
	}
}

type OrganizationCreateOptions struct {
	Name  string  `json:"name" validate:"required,max=255"`
	Image *string `json:"image"`
}

type Organization struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Image      *string `json:"image,omitempty"`
	Permission string  `json:"permission"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
}

func (svc *OrganizationService) Create(opts OrganizationCreateOptions, userID string) (*Organization, error) {
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
	res, err := svc.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Find(id string, userID string) (*Organization, error) {
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type OrganizationListOptions struct {
	Query     string
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type OrganizationList struct {
	Data          []*Organization `json:"data"`
	TotalPages    uint            `json:"totalPages"`
	TotalElements uint            `json:"totalElements"`
	Page          uint            `json:"page"`
	Size          uint            `json:"size"`
}

func (svc *OrganizationService) List(opts OrganizationListOptions, userID string) (*OrganizationList, error) {
	var authorized []model.Organization
	if opts.Query == "" {
		ids, err := svc.orgRepo.GetIDs()
		if err != nil {
			return nil, err
		}
		authorized, err = svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
	} else {
		orgs, err := svc.orgSearch.Query(opts.Query)
		if err != nil {
			return nil, err
		}
		authorized, err = svc.doAuthorization(orgs, userID)
		if err != nil {
			return nil, err
		}
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(authorized, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.orgMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &OrganizationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *OrganizationService) PatchName(id string, name string, userID string) (*Organization, error) {
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
	res, err := svc.orgMapper.mapOne(org, userID)
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

	/* Make sure member is not the last remaining owner of the organization */
	ownerCount, err := svc.orgRepo.GetOwnerCount(org.GetID())
	if err != nil {
		return err
	}
	if svc.orgGuard.IsAuthorized(memberID, org, model.PermissionOwner) && ownerCount == 1 {
		return errorpkg.NewCannotRemoveLastRemainingOwnerOfOrganizationError(org.GetID())
	}

	if userID != memberID {
		if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
			return err
		}
	}

	/* Remove member from all groups belonging to this organization */
	groupsIDs, err := svc.groupRepo.GetIDsForOrganization(org.GetID())
	if err != nil {
		return err
	}
	for _, groupID := range groupsIDs {
		if err := svc.groupService.RemoveMemberUnauthorized(groupID, memberID); err != nil {
			log.GetLogger().Error(err)
		}
	}

	if err := svc.orgRepo.RevokeUserPermission(id, memberID); err != nil {
		return err
	}
	if err := svc.orgRepo.RemoveMember(id, memberID); err != nil {
		return err
	}
	if _, err := svc.orgCache.Refresh(org.GetID()); err != nil {
		return err
	}
	return nil
}

func (svc *OrganizationService) doAuthorization(data []model.Organization, userID string) ([]model.Organization, error) {
	var res []model.Organization
	for _, o := range data {
		if svc.orgGuard.IsAuthorized(userID, o, model.PermissionViewer) {
			res = append(res, o)
		}
	}
	return res, nil
}

func (svc *OrganizationService) doAuthorizationByIDs(ids []string, userID string) ([]model.Organization, error) {
	var res []model.Organization
	for _, id := range ids {
		var o model.Organization
		o, err := svc.orgCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.orgGuard.IsAuthorized(userID, o, model.PermissionViewer) {
			res = append(res, o)
		}
	}
	return res, nil
}

func (svc *OrganizationService) doSorting(data []model.Organization, sortBy string, sortOrder string) []model.Organization {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == SortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a, _ := time.Parse(time.RFC3339, data[i].GetCreateTime())
			b, _ := time.Parse(time.RFC3339, data[j].GetCreateTime())
			if sortOrder == SortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == SortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a, _ := time.Parse(time.RFC3339, *data[i].GetUpdateTime())
				b, _ := time.Parse(time.RFC3339, *data[j].GetUpdateTime())
				if sortOrder == SortOrderDesc {
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

func (svc *OrganizationService) doPagination(data []model.Organization, page, size uint) (pageData []model.Organization, totalElements uint, totalPages uint) {
	totalElements = uint(len(data))
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

type organizationMapper struct {
	groupCache *cache.GroupCache
}

func newOrganizationMapper() *organizationMapper {
	return &organizationMapper{
		groupCache: cache.NewGroupCache(),
	}
}

func (mp *organizationMapper) mapOne(m model.Organization, userID string) (*Organization, error) {
	res := &Organization{
		ID:         m.GetID(),
		Name:       m.GetName(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	res.Permission = ""
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
		for _, u := range g.GetUsers() {
			if u == userID && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
				res.Permission = p.GetValue()
			}
		}
	}
	return res, nil
}

func (mp *organizationMapper) mapMany(orgs []model.Organization, userID string) ([]*Organization, error) {
	res := make([]*Organization, 0)
	for _, org := range orgs {
		o, err := mp.mapOne(org, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}
