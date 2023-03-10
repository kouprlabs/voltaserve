package core

import (
	"fmt"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	log "github.com/sirupsen/logrus"
)

type Organization struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Image      *string `json:"image,omitempty"`
	Permission string  `json:"permission"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
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
	UserId string `json:"userId" validate:"required"`
}

type OrganizationService struct {
	orgRepo      *repo.OrganizationRepo
	orgCache     *cache.OrganizationCache
	orgGuard     *guard.OrganizationGuard
	orgMapper    *organizationMapper
	orgSearch    *search.OrganizationSearch
	userRepo     *repo.UserRepo
	userSearch   *search.UserSearch
	userMapper   *userMapper
	groupRepo    *repo.GroupRepo
	groupService *GroupService
	groupMapper  *groupMapper
	imageProc    *infra.ImageProcessor
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
		imageProc:    infra.NewImageProcessor(),
		userMapper:   newUserMapper(),
		config:       config.GetConfig(),
	}
}

func (svc *OrganizationService) Create(req OrganizationCreateOptions, userId string) (*Organization, error) {
	org, err := svc.orgRepo.Insert(repo.OrganizationInsertOptions{
		Id:   helpers.NewId(),
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.orgRepo.GrantUserPermission(org.GetId(), userId, model.PermissionOwner); err != nil {
		return nil, err
	}
	org, err = svc.orgRepo.Find(org.GetId())
	if err != nil {
		return nil, err
	}
	if err := svc.orgSearch.Index([]model.OrganizationModel{org}); err != nil {
		return nil, err
	}
	if err := svc.orgCache.Set(org); err != nil {
		return nil, nil
	}
	res, err := svc.orgMapper.mapOrganization(org, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Find(id string, userId string) (*Organization, error) {
	user, err := svc.userRepo.Find(userId)
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
	res, err := svc.orgMapper.mapOrganization(org, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Search(query string, userId string) ([]*Organization, error) {
	orgs, err := svc.orgSearch.Query(query)
	if err != nil {
		return nil, err
	}
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	org := make([]*Organization, 0)
	for _, o := range orgs {
		if svc.orgGuard.IsAuthorized(user, o, model.PermissionViewer) {
			v, err := svc.orgMapper.mapOrganization(o, userId)
			if err != nil {
				return nil, err
			}
			org = append(org, v)
		}
	}
	return org, nil
}

func (svc *OrganizationService) SearchMembers(id string, query string, userId string) ([]*User, error) {
	user, err := svc.userRepo.Find(userId)
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
	var members []model.UserModel
	for _, m := range orgMembers {
		for _, u := range users {
			if u.GetId() == m.GetId() {
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

func (svc *OrganizationService) FindAll(userId string) ([]*Organization, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	ids, err := svc.orgRepo.GetIds()
	if err != nil {
		return nil, err
	}
	res := make([]*Organization, 0)
	for _, id := range ids {
		org, err := svc.orgCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.orgGuard.IsAuthorized(user, org, model.PermissionViewer) {
			v, err := svc.orgMapper.mapOrganization(org, userId)
			if err != nil {
				return nil, err
			}
			res = append(res, v)
		}
	}
	return res, nil
}

func (svc *OrganizationService) UpdateName(id string, name string, userId string) (*Organization, error) {
	user, err := svc.userRepo.Find(userId)
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
	if err := svc.orgSearch.Update([]model.OrganizationModel{org}); err != nil {
		return nil, err
	}
	err = svc.orgCache.Set(org)
	if err != nil {
		return nil, err
	}
	res, err := svc.orgMapper.mapOrganization(org, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *OrganizationService) Delete(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
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
	if err := svc.orgCache.Delete(org.GetId()); err != nil {
		return err
	}
	if err := svc.orgSearch.Delete([]string{org.GetId()}); err != nil {
		return err
	}
	return nil
}

func (svc *OrganizationService) RemoveMember(id string, memberId string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	member, err := svc.userRepo.Find(memberId)
	if err != nil {
		return err
	}
	org, err := svc.orgCache.Get(id)
	if err != nil {
		return err
	}

	/* Make sure member is not the last remaining owner of the organization */
	ownerCount, err := svc.orgRepo.GetOwnerCount(org.GetId())
	if err != nil {
		return err
	}
	if svc.orgGuard.IsAuthorized(member, org, model.PermissionOwner) && ownerCount == 1 {
		return errorpkg.NewCannotRemoveLastRemainingOwnerOfOrganizationError(org.GetId())
	}

	if userId != member.GetId() {
		if err := svc.orgGuard.Authorize(user, org, model.PermissionEditor); err != nil {
			return err
		}
	}

	/* Remove member from all groups belonging to this organization */
	groupsIds, err := svc.groupRepo.GetIdsForOrganization(org.GetId())
	if err != nil {
		return err
	}
	for _, groupId := range groupsIds {
		if err := svc.groupService.RemoveMemberUnauthorized(groupId, member.GetId()); err != nil {
			log.Error(err)
		}
	}

	if err := svc.orgRepo.RevokeUserPermission(id, member.GetId()); err != nil {
		return err
	}
	if err := svc.orgRepo.RemoveMember(id, member.GetId()); err != nil {
		return err
	}
	if _, err := svc.orgCache.Refresh(org.GetId()); err != nil {
		return err
	}
	return nil
}

func (svc *OrganizationService) GetMembers(id string, userId string) ([]*User, error) {
	user, err := svc.userRepo.Find(userId)
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

func (svc *OrganizationService) GetGroups(id string, userId string) ([]*Group, error) {
	user, err := svc.userRepo.Find(userId)
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
	res, err := svc.groupMapper.mapGroups(groups, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type organizationMapper struct {
	groupCache *cache.GroupCache
}

func newOrganizationMapper() *organizationMapper {
	return &organizationMapper{
		groupCache: cache.NewGroupCache(),
	}
}

func (mp *organizationMapper) mapOrganization(m model.OrganizationModel, userId string) (*Organization, error) {
	res := &Organization{
		Id:         m.GetId(),
		Name:       m.GetName(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	res.Permission = ""
	for _, p := range m.GetUserPermissions() {
		if p.GetUserId() == userId && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
			res.Permission = p.GetValue()
		}
	}
	for _, p := range m.GetGroupPermissions() {
		g, err := mp.groupCache.Get(p.GetGroupId())
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
