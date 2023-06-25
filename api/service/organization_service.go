package service

import (
	"fmt"
	"sort"
	"time"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	log "github.com/sirupsen/logrus"
)

type Organization struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Image      *string `json:"image,omitempty"`
	Permission string  `json:"permission"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
}

type OrganizationList struct {
	Data          []*Organization `json:"data"`
	TotalPages    uint            `json:"totalPages"`
	TotalElements uint            `json:"totalElements"`
	Page          uint            `json:"page"`
	Size          uint            `json:"size"`
}

type OrganizationSearchOptions struct {
	Text string `json:"text" validate:"required"`
}

type OrganizationCreateOptions struct {
	Name  string  `json:"name" validate:"required,max=255"`
	Image *string `json:"image"`
}

type OrganizationUpdateNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type OrganizationUpdateImageOptions struct {
	Image string `json:"image" validate:"required,base64"`
}

type OrganizationRemoveMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

type OrganizationService struct {
	orgRepo      repo.OrganizationRepo
	orgCache     *cache.OrganizationCache
	orgGuard     *guard.OrganizationGuard
	orgMapper    *organizationMapper
	orgSearch    *search.OrganizationSearch
	userRepo     repo.UserRepo
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
		userRepo:     repo.NewUserRepo(),
		userSearch:   search.NewUserSearch(),
		groupRepo:    repo.NewGroupRepo(),
		groupService: NewGroupService(),
		groupMapper:  newGroupMapper(),
		userMapper:   newUserMapper(),
		config:       config.GetConfig(),
	}
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
	org, err = svc.orgRepo.Find(org.GetID())
	if err != nil {
		return nil, err
	}
	if err := svc.orgSearch.Index([]model.Organization{org}); err != nil {
		return nil, err
	}
	if err := svc.orgCache.Set(org); err != nil {
		return nil, nil
	}
	res, err := svc.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Find(id string, userID string) (*Organization, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) List(page uint, size uint, sortBy string, sortOrder string, userID string) (*OrganizationList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	ids, err := svc.orgRepo.GetIDs()
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorizationByIDs(ids, user)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(authorized, sortBy, sortOrder, userID)
	paged, totalElements, totalPages := svc.doPaging(sorted, page, size)
	mapped, err := svc.orgMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &OrganizationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          page,
		Size:          size,
	}, nil
}

func (svc *OrganizationService) Search(query string, page uint, size uint, userID string) (*OrganizationList, error) {
	workspaces, err := svc.orgSearch.Query(query)
	if err != nil {
		return nil, err
	}
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorization(workspaces, user)
	if err != nil {
		return nil, err
	}
	paged, totalElements, totalPages := svc.doPaging(authorized, page, size)
	mapped, err := svc.orgMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &OrganizationList{
		Data:          mapped,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          page,
		Size:          size,
	}, nil
}

func (svc *OrganizationService) SearchMembers(id string, query string, userID string) ([]*User, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionViewer); err != nil {
		return nil, err
	}
	users, err := svc.userSearch.Query(fmt.Sprintf(`
	{
		"query": {
			"query_string": {
				"fields": ["email"],
				"query": "%s",
				"fuzziness": "AUTO"
			}
		}
	}
	`, query))
	if err != nil {
		return nil, err
	}
	orgMembers, err := svc.orgRepo.GetMembers(id)
	if err != nil {
		return nil, err
	}
	var members []model.User
	for _, m := range orgMembers {
		for _, u := range users {
			if u.GetID() == m.GetID() {
				members = append(members, m)
			}
		}
	}
	res, err := svc.userMapper.mapUsers(members)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) UpdateName(id string, name string, userID string) (*Organization, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionEditor); err != nil {
		return nil, err
	}
	org.SetName(name)
	if err := svc.orgRepo.Save(org); err != nil {
		return nil, err
	}
	if err := svc.orgSearch.Update([]model.Organization{org}); err != nil {
		return nil, err
	}
	err = svc.orgCache.Set(org)
	if err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Delete(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionOwner); err != nil {
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
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	member, err := svc.userRepo.Find(memberID)
	if err != nil {
		return err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return err
	}

	/* Make sure member is not the last remaining owner of the organization */
	ownerCount, err := svc.orgRepo.GetOwnerCount(org.GetID())
	if err != nil {
		return err
	}
	if svc.orgGuard.IsAuthorized(member, org, model.PermissionOwner) && ownerCount == 1 {
		return errorpkg.NewCannotRemoveLastRemainingOwnerOfOrganizationError(org.GetID())
	}

	if userID != member.GetID() {
		if err := svc.orgGuard.Authorize(user, org, model.PermissionOwner); err != nil {
			return err
		}
	}

	/* Remove member from all groups belonging to this organization */
	groupsIDs, err := svc.groupRepo.GetIDsForOrganization(org.GetID())
	if err != nil {
		return err
	}
	for _, groupID := range groupsIDs {
		if err := svc.groupService.RemoveMemberUnauthorized(groupID, member.GetID()); err != nil {
			log.Error(err)
		}
	}

	if err := svc.orgRepo.RevokeUserPermission(id, member.GetID()); err != nil {
		return err
	}
	if err := svc.orgRepo.RemoveMember(id, member.GetID()); err != nil {
		return err
	}
	if _, err := svc.orgCache.Refresh(org.GetID()); err != nil {
		return err
	}
	return nil
}

func (svc *OrganizationService) GetMembers(id string, userID string) ([]*User, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionViewer); err != nil {
		return nil, err
	}
	members, err := svc.orgRepo.GetMembers(id)
	if err != nil {
		return nil, err
	}
	res, err := svc.userMapper.mapUsers(members)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) GetGroups(id string, userID string) ([]*Group, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionViewer); err != nil {
		return nil, err
	}
	groups, err := svc.orgRepo.GetGroups(id)
	if err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapGroups(groups, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) doAuthorization(data []model.Organization, user model.User) ([]model.Organization, error) {
	var res []model.Organization
	for _, w := range data {
		if svc.orgGuard.IsAuthorized(user, w, model.PermissionViewer) {
			res = append(res, w)
		}
	}
	return res, nil
}

func (svc *OrganizationService) doAuthorizationByIDs(ids []string, user model.User) ([]model.Organization, error) {
	var res []model.Organization
	for _, id := range ids {
		var org model.Organization
		org, err := svc.orgCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.orgGuard.IsAuthorized(user, org, model.PermissionViewer) {
			res = append(res, org)
		}
	}
	return res, nil
}

func (svc *OrganizationService) doSorting(data []model.Organization, sortBy string, sortOrder string, userID string) []model.Organization {
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

func (svc *OrganizationService) doPaging(data []model.Organization, page uint, size uint) (pageData []model.Organization, totalElements uint, totalPages uint) {
	page = page - 1
	low := size * page
	high := low + size
	if low >= uint(len(data)) {
		pageData = []model.Organization{}
	} else if high >= uint(len(data)) {
		high = uint(len(data))
		pageData = data[low:high]
	} else {
		pageData = data[low:high]
	}
	totalElements = uint(len(data))
	if totalElements == 0 {
		totalPages = 1
	} else {
		if size > uint(len(data)) {
			size = uint(len(data))
		}
		totalPages = totalElements / size
		if totalPages == 0 {
			totalPages = 1
		}
		if totalElements%size > 0 {
			totalPages = totalPages + 1
		}
	}
	return pageData, totalElements, totalPages
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
	for _, f := range orgs {
		v, err := mp.mapOne(f, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
