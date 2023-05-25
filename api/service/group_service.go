package service

import (
	"fmt"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/guard"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type Group struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Image        *string      `json:"image,omitempty"`
	Organization Organization `json:"organization"`
	Permission   string       `json:"permission"`
	CreateTime   string       `json:"createTime,omitempty"`
	UpdateTime   *string      `json:"updateTime"`
}

type GroupSearchOptions struct {
	Text string `json:"text" validate:"required"`
}

type GroupCreateOptions struct {
	Name           string  `json:"name" validate:"required,max=255"`
	Image          *string `json:"image"`
	OrganizationId string  `json:"organizationId" validate:"required"`
}

type GroupUpdateNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type GroupUpdateImageOptions struct {
	Image string `json:"image" validate:"required,base64"`
}

type GroupAddMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

type GroupRemoveMemberOptions struct {
	UserID string `json:"userId" validate:"required"`
}

type GroupService struct {
	groupRepo      repo.CoreGroupRepo
	groupGuard     *guard.GroupGuard
	groupSearch    *search.GroupSearch
	groupMapper    *groupMapper
	groupCache     *cache.GroupCache
	userRepo       repo.CoreUserRepo
	userSearch     *search.UserSearch
	userMapper     *userMapper
	workspaceRepo  repo.CoreWorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	fileRepo       repo.CoreFileRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	orgRepo        repo.CoreOrganizationRepo
	orgCache       *cache.OrganizationCache
	orgGuard       *guard.OrganizationGuard
	imageProc      *infra.ImageProcessor
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
		imageProc:      infra.NewImageProcessor(),
		fileGuard:      guard.NewFileGuard(),
		config:         config.GetConfig(),
	}
}

func (svc *GroupService) Create(req GroupCreateOptions, userId string) (*Group, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	org, err := svc.orgCache.Get(req.OrganizationId)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionEditor); err != nil {
		return nil, err
	}
	group, err := svc.groupRepo.Insert(repo.GroupInsertOptions{
		ID:             helpers.NewId(),
		Name:           req.Name,
		OrganizationId: req.OrganizationId,
		OwnerId:        userId,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.groupRepo.GrantUserPermission(group.GetID(), userId, model.PermissionOwner); err != nil {
		return nil, err
	}
	group, err = svc.groupRepo.Find(group.GetID())
	if err != nil {
		return nil, err
	}
	if err := svc.groupSearch.Index([]model.CoreGroup{group}); err != nil {
		return nil, err
	}
	if err := svc.groupCache.Set(group); err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapGroup(group, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) Find(id string, userId string) (*Group, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapGroup(group, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) FindAll(userId string) ([]*Group, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	ids, err := svc.groupRepo.GetIDs()
	if err != nil {
		return nil, err
	}
	res := make([]*Group, 0)
	for _, id := range ids {
		group, err := svc.groupCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.groupGuard.IsAuthorized(user, group, model.PermissionViewer) {
			dto, err := svc.groupMapper.mapGroup(group, userId)
			if err != nil {
				return nil, err
			}
			res = append(res, dto)
		}
	}
	return res, nil
}

func (svc *GroupService) FindAllForFile(fileId string, userId string) ([]*Group, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileId)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	ids, err := svc.groupRepo.GetIDsForFile(fileId)
	if err != nil {
		return nil, err
	}
	res := make([]*Group, 0)
	for _, id := range ids {
		group, err := svc.groupCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.groupGuard.IsAuthorized(user, group, model.PermissionViewer) {
			dto, err := svc.groupMapper.mapGroup(group, userId)
			if err != nil {
				return nil, err
			}
			res = append(res, dto)
		}
	}
	return res, nil
}

func (svc *GroupService) Search(query string, userId string) ([]*Group, error) {
	groups, err := svc.groupSearch.Query(query)
	if err != nil {
		return nil, err
	}
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	var res []*Group
	for _, g := range groups {
		if svc.groupGuard.IsAuthorized(user, g, model.PermissionViewer) {
			dto, err := svc.groupMapper.mapGroup(g, userId)
			if err != nil {
				return nil, err
			}
			res = append(res, dto)
		}
	}
	return res, nil
}

func (svc *GroupService) UpdateName(id string, name string, userId string) (*Group, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionEditor); err != nil {
		return nil, err
	}
	group.SetName(name)
	if err := svc.groupRepo.Save(group); err != nil {
		return nil, err
	}
	if err := svc.groupSearch.Update([]model.CoreGroup{group}); err != nil {
		return nil, err
	}
	err = svc.groupCache.Set(group)
	if err != nil {
		return nil, err
	}
	res, err := svc.groupMapper.mapGroup(group, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) Delete(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionOwner); err != nil {
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

func (svc *GroupService) AddMember(id string, memberId string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
		return err
	}
	if _, err := svc.userRepo.Find(memberId); err != nil {
		return err
	}
	if err := svc.groupRepo.AddUser(id, memberId); err != nil {
		return err
	}
	if err := svc.groupRepo.GrantUserPermission(group.GetID(), memberId, model.PermissionViewer); err != nil {
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

func (svc *GroupService) RemoveMember(id string, memberId string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
		return err
	}
	if err := svc.RemoveMemberUnauthorized(id, memberId); err != nil {
		return err
	}
	return nil
}

func (svc *GroupService) RemoveMemberUnauthorized(id string, memberId string) error {
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil
	}
	if _, err := svc.userRepo.Find(memberId); err != nil {
		return err
	}
	if err := svc.groupRepo.RemoveMember(id, memberId); err != nil {
		return err
	}
	if err := svc.groupRepo.RevokeUserPermission(id, memberId); err != nil {
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

func (svc *GroupService) refreshCacheForOrganization(organizationId string) error {
	workspaceIds, err := svc.workspaceRepo.GetIdsByOrganization(organizationId)
	if err != nil {
		return err
	}
	for _, workspaceId := range workspaceIds {
		if _, err := svc.workspaceCache.Refresh(workspaceId); err != nil {
			return err
		}
		filesIds, err := svc.fileRepo.GetIdsByWorkspace(workspaceId)
		if err != nil {
			return err
		}
		for _, id := range filesIds {
			if _, err := svc.fileCache.Refresh(id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *GroupService) GetMembers(id string, userId string) ([]*User, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
		return nil, err
	}
	members, err := svc.groupRepo.GetMembers(id)
	if err != nil {
		return nil, err
	}
	res, err := svc.userMapper.mapUsers(members)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GroupService) SearchMembers(id string, query string, userId string) ([]*User, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
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
	groupMembers, err := svc.groupRepo.GetMembers(id)
	if err != nil {
		return nil, err
	}
	var members []model.CoreUser
	for _, m := range groupMembers {
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

func (svc *GroupService) GetAvailableUsers(id string, userId string) ([]*User, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	group, err := svc.groupCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
		return nil, err
	}
	orgMembers, err := svc.orgRepo.GetMembers(group.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	groupMembers, err := svc.groupRepo.GetMembers(id)
	if err != nil {
		return nil, err
	}
	var res []*User
	for _, om := range orgMembers {
		found := false
		for _, tm := range groupMembers {
			if om.GetID() == tm.GetID() {
				found = true
				break
			}
		}
		if !found {
			res = append(res, svc.userMapper.mapUser(om))
		}
	}
	return res, nil
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

func (mp *groupMapper) mapGroup(m model.CoreGroup, userId string) (*Group, error) {
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	v, err := mp.orgMapper.mapOrganization(org, userId)
	if err != nil {
		return nil, err
	}
	res := &Group{
		ID:           m.GetID(),
		Name:         m.GetName(),
		Organization: *v,
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}
	res.Permission = ""
	for _, p := range m.GetUserPermissions() {
		if p.GetUserID() == userId && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
			res.Permission = p.GetValue()
		}
	}
	for _, p := range m.GetGroupPermissions() {
		g, err := mp.groupCache.Get(p.GetGroupID())
		if err != nil {
			return nil, err
		}
		for _, u := range g.GetUsers() {
			if u == userId && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
				res.Permission = p.GetValue()
			}
		}
	}
	return res, nil
}

func (mp *groupMapper) mapGroups(groups []model.CoreGroup, userId string) ([]*Group, error) {
	res := []*Group{}
	for _, g := range groups {
		v, err := mp.mapGroup(g, userId)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
