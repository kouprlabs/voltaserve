package service

import (
	"sort"
	"time"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type GroupService struct {
	groupRepo      repo.GroupRepo
	groupGuard     *guard.GroupGuard
	groupSearch    *search.GroupSearch
	groupMapper    *groupMapper
	groupCache     *cache.GroupCache
	userRepo       repo.UserRepo
	userSearch     *search.UserSearch
	userMapper     *userMapper
	workspaceRepo  repo.WorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	fileRepo       repo.FileRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	orgRepo        repo.OrganizationRepo
	orgCache       *cache.OrganizationCache
	orgGuard       *guard.OrganizationGuard
	config         config.Config
}

func NewGroupService() *GroupService {
	return &GroupService{
		groupRepo:      repo.NewGroupRepo(),
		groupGuard:     guard.NewGroupGuard(),
		groupCache:     cache.NewGroupCache(),
		groupSearch:    search.NewGroupSearch(),
		groupMapper:    newGroupMapper(),
		userRepo:       repo.NewUserRepo(),
		userSearch:     search.NewUserSearch(),
		userMapper:     newUserMapper(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceCache: cache.NewWorkspaceCache(),
		fileRepo:       repo.NewFileRepo(),
		fileCache:      cache.NewFileCache(),
		orgRepo:        repo.NewOrganizationRepo(),
		orgGuard:       guard.NewOrganizationGuard(),
		orgCache:       cache.NewOrganizationCache(),
		fileGuard:      guard.NewFileGuard(),
		config:         config.GetConfig(),
	}
}

type GroupCreateOptions struct {
	Name           string  `json:"name" validate:"required,max=255"`
	Image          *string `json:"image"`
	OrganizationID string  `json:"organizationId" validate:"required"`
}

type Group struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Image        *string      `json:"image,omitempty"`
	Organization Organization `json:"organization"`
	Permission   string       `json:"permission"`
	CreateTime   string       `json:"createTime,omitempty"`
	UpdateTime   *string      `json:"updateTime"`
}

func (svc *GroupService) Create(opts GroupCreateOptions, userID string) (*Group, error) {
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
	group, err = svc.groupRepo.Find(group.GetID())
	if err != nil {
		return nil, err
	}
	if err := svc.groupSearch.Index([]model.Group{group}); err != nil {
		return nil, err
	}
	if err := svc.groupCache.Set(group); err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapOne(group, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) Find(id string, userID string) (*Group, error) {
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
	Page           uint
	Size           uint
	SortBy         string
	SortOrder      string
}

type GroupList struct {
	Data          []*Group `json:"data"`
	TotalPages    uint     `json:"totalPages"`
	TotalElements uint     `json:"totalElements"`
	Page          uint     `json:"page"`
	Size          uint     `json:"size"`
}

func (svc *GroupService) List(opts GroupListOptions, userID string) (*GroupList, error) {
	var authorized []model.Group
	if opts.Query == "" {
		if opts.OrganizationID == "" {
			ids, err := svc.groupRepo.GetIDs()
			if err != nil {
				return nil, err
			}
			authorized, err = svc.doAuthorizationByIDs(ids, userID)
			if err != nil {
				return nil, err
			}
		} else {
			groups, err := svc.orgRepo.GetGroups(opts.OrganizationID)
			if err != nil {
				return nil, err
			}
			authorized, err = svc.doAuthorization(groups, userID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		groups, err := svc.groupSearch.Query(opts.Query)
		if err != nil {
			return nil, err
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
		authorized, err = svc.doAuthorization(filtered, userID)
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
	mapped, err := svc.groupMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &GroupList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *GroupService) PatchName(id string, name string, userID string) (*Group, error) {
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
	if err := svc.groupSearch.Update([]model.Group{group}); err != nil {
		return nil, err
	}
	err = svc.groupCache.Set(group)
	if err != nil {
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
		return nil
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
	if err := svc.refreshCacheForOrganization(group.GetOrganizationID()); err != nil {
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
	if _, err := svc.userRepo.Find(memberID); err != nil {
		return err
	}
	if err := svc.groupRepo.AddUser(id, memberID); err != nil {
		return err
	}
	if err := svc.groupRepo.GrantUserPermission(group.GetID(), memberID, model.PermissionViewer); err != nil {
		return err
	}
	if _, err := svc.groupCache.Refresh(group.GetID()); err != nil {
		return err
	}
	if err := svc.refreshCacheForOrganization(group.GetOrganizationID()); err != nil {
		return err
	}
	return nil
}

func (svc *GroupService) RemoveMember(id string, memberID string, userID string) error {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionOwner); err != nil {
		return err
	}
	if err := svc.RemoveMemberUnauthorized(id, memberID); err != nil {
		return err
	}
	return nil
}

func (svc *GroupService) RemoveMemberUnauthorized(id string, memberID string) error {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if _, err := svc.userRepo.Find(memberID); err != nil {
		return err
	}
	if err := svc.groupRepo.RemoveMember(id, memberID); err != nil {
		return err
	}
	if err := svc.groupRepo.RevokeUserPermission(id, memberID); err != nil {
		return err
	}
	if _, err := svc.groupCache.Refresh(group.GetID()); err != nil {
		return err
	}
	if err := svc.refreshCacheForOrganization(group.GetOrganizationID()); err != nil {
		return err
	}
	return nil
}

func (svc *GroupService) refreshCacheForOrganization(orgID string) error {
	workspaceIDs, err := svc.workspaceRepo.GetIDsByOrganization(orgID)
	if err != nil {
		return err
	}
	for _, workspaceID := range workspaceIDs {
		if _, err := svc.workspaceCache.Refresh(workspaceID); err != nil {
			return err
		}
		filesIDs, err := svc.fileRepo.GetIDsByWorkspace(workspaceID)
		if err != nil {
			return err
		}
		for _, id := range filesIDs {
			if _, err := svc.fileCache.Refresh(id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *GroupService) doAuthorization(data []model.Group, userID string) ([]model.Group, error) {
	var res []model.Group
	for _, g := range data {
		if svc.groupGuard.IsAuthorized(userID, g, model.PermissionViewer) {
			res = append(res, g)
		}
	}
	return res, nil
}

func (svc *GroupService) doAuthorizationByIDs(ids []string, userID string) ([]model.Group, error) {
	var res []model.Group
	for _, id := range ids {
		var o model.Group
		o, err := svc.groupCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.groupGuard.IsAuthorized(userID, o, model.PermissionViewer) {
			res = append(res, o)
		}
	}
	return res, nil
}

func (svc *GroupService) doSorting(data []model.Group, sortBy string, sortOrder string) []model.Group {
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

func (svc *GroupService) doPagination(data []model.Group, page, size uint) ([]model.Group, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Group{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

type groupMapper struct {
	orgCache   *cache.OrganizationCache
	orgMapper  *organizationMapper
	groupCache *cache.GroupCache
}

func newGroupMapper() *groupMapper {
	return &groupMapper{
		orgCache:   cache.NewOrganizationCache(),
		orgMapper:  newOrganizationMapper(),
		groupCache: cache.NewGroupCache(),
	}
}

func (mp *groupMapper) mapOne(m model.Group, userID string) (*Group, error) {
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	o, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	res := &Group{
		ID:           m.GetID(),
		Name:         m.GetName(),
		Organization: *o,
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
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

func (mp *groupMapper) mapMany(groups []model.Group, userID string) ([]*Group, error) {
	res := []*Group{}
	for _, group := range groups {
		g, err := mp.mapOne(group, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, g)
	}
	return res, nil
}
