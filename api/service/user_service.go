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
	"slices"
	"sort"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type UserService struct {
	userMapper     *userMapper
	userRepo       *repo.UserRepo
	userSearch     *search.UserSearch
	orgRepo        *repo.OrganizationRepo
	orgCache       *cache.OrganizationCache
	orgGuard       *guard.OrganizationGuard
	groupRepo      *repo.GroupRepo
	groupGuard     *guard.GroupGuard
	groupCache     *cache.GroupCache
	invitationRepo *repo.InvitationRepo
	config         *config.Config
}

func NewUserService() *UserService {
	return &UserService{
		userMapper:     newUserMapper(),
		userRepo:       repo.NewUserRepo(),
		userSearch:     search.NewUserSearch(),
		orgRepo:        repo.NewOrganizationRepo(),
		orgCache:       cache.NewOrganizationCache(),
		orgGuard:       guard.NewOrganizationGuard(),
		groupRepo:      repo.NewGroupRepo(),
		groupGuard:     guard.NewGroupGuard(),
		groupCache:     cache.NewGroupCache(),
		invitationRepo: repo.NewInvitationRepo(),
		config:         config.GetConfig(),
	}
}

type UserListOptions struct {
	Query               string
	OrganizationID      string
	GroupID             string
	ExcludeGroupMembers bool
	ExcludeMe           bool
	SortBy              string
	SortOrder           string
	Page                uint64
	Size                uint64
}

func (svc *UserService) List(opts UserListOptions, userID string) (*dto.UserList, error) {
	users, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.UserSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.UserSortOrderAsc
	}
	sorted := svc.sort(users, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.userMapper.mapMany(paged)
	if err != nil {
		return nil, err
	}
	return &dto.UserList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *UserService) Probe(opts UserListOptions, userID string) (*dto.UserProbe, error) {
	users, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(users))
	return &dto.UserProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

type ExtractPictureJustification struct {
	OrganizationID *string
	GroupID        *string
	InvitationID   *string
}

func (svc *UserService) ExtractPicture(id string, justification ExtractPictureJustification, userID string, isAdmin bool) ([]byte, *string, *string, error) {
	user, err := svc.findUserForPicture(id, justification, userID, isAdmin)
	if err != nil {
		return nil, nil, nil, err
	}
	if user.GetPicture() == nil {
		return nil, nil, nil, errorpkg.NewPictureNotFoundError(nil)
	}
	mime := helper.Base64ToMIME(*user.GetPicture())
	ext := helper.Base64ToExtension(*user.GetPicture())
	b, err := helper.Base64ToBytes(*user.GetPicture())
	if err != nil {
		return nil, nil, nil, errorpkg.NewPictureNotFoundError(nil)
	}
	return b, &ext, &mime, nil
}

func (svc *UserService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.UserSortByEmail ||
		value == dto.UserSortByFullName ||
		value == dto.UserSortByDateCreated ||
		value == dto.UserSortByDateModified
}

func (svc *UserService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.UserSortOrderAsc || value == dto.UserSortOrderDesc
}

func (svc *UserService) findAll(opts UserListOptions, userID string) ([]model.User, error) {
	if opts.OrganizationID == "" && opts.GroupID == "" {
		return make([]model.User, 0), nil
	}
	if opts.OrganizationID != "" {
		org, err := svc.orgCache.Get(opts.OrganizationID)
		if err != nil {
			return nil, err
		}
		if err := svc.orgGuard.Authorize(userID, org, model.PermissionViewer); err != nil {
			return nil, err
		}
	}
	if opts.GroupID != "" {
		group, err := svc.groupCache.Get(opts.GroupID)
		if err != nil {
			return nil, err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return nil, err
		}
	}
	var res []model.User
	var err error
	if opts.Query == "" {
		res, err = svc.load(opts)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = svc.search(opts)
		if err != nil {
			return nil, err
		}
	}
	if opts.ExcludeMe {
		withoutMe := make([]model.User, 0)
		for _, u := range res {
			if u.GetID() != userID {
				withoutMe = append(withoutMe, u)
			}
		}
		return withoutMe, nil
	} else {
		return res, nil
	}
}

func (svc *UserService) load(opts UserListOptions) ([]model.User, error) {
	res := make([]model.User, 0)
	var err error
	if opts.OrganizationID != "" && opts.GroupID != "" && opts.ExcludeGroupMembers {
		orgMembers, err := svc.orgRepo.FindMembers(opts.OrganizationID)
		if err != nil {
			return nil, err
		}
		groupMembers, err := svc.groupRepo.FindMembers(opts.GroupID)
		if err != nil {
			return nil, err
		}
		for _, om := range orgMembers {
			isGroupMember := false
			for _, gm := range groupMembers {
				if om.GetID() == gm.GetID() {
					isGroupMember = true
					break
				}
			}
			if !isGroupMember {
				res = append(res, om)
			}
		}
	} else if opts.OrganizationID != "" {
		res, err = svc.orgRepo.FindMembers(opts.OrganizationID)
		if err != nil {
			return nil, err
		}
	} else if opts.GroupID != "" {
		res, err = svc.groupRepo.FindMembers(opts.GroupID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *UserService) search(opts UserListOptions) ([]model.User, error) {
	res := make([]model.User, 0)
	var err error
	count, err := svc.userRepo.Count()
	if err != nil {
		return nil, err
	}
	users, err := svc.userSearch.Query(opts.Query, infra.SearchQueryOptions{Limit: count})
	if err != nil {
		return nil, err
	}
	var members []model.User
	if opts.OrganizationID != "" {
		members, err = svc.orgRepo.FindMembers(opts.OrganizationID)
		if err != nil {
			return nil, err
		}
	} else if opts.GroupID != "" {
		members, err = svc.groupRepo.FindMembers(opts.GroupID)
		if err != nil {
			return nil, err
		}
	}
	for _, m := range members {
		for _, u := range users {
			if u.GetID() == m.GetID() {
				res = append(res, m)
			}
		}
	}
	return res, nil
}

func (svc *UserService) sort(data []model.User, sortBy string, sortOrder string) []model.User {
	if sortBy == dto.UserSortByEmail {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.UserSortOrderDesc {
				return data[i].GetEmail() > data[j].GetEmail()
			} else {
				return data[i].GetEmail() < data[j].GetEmail()
			}
		})
		return data
	} else if sortBy == dto.UserSortByFullName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.UserSortOrderDesc {
				return data[i].GetFullName() > data[j].GetFullName()
			} else {
				return data[i].GetFullName() < data[j].GetFullName()
			}
		})
		return data
	}
	return data
}

func (svc *UserService) paginate(data []model.User, page, size uint64) ([]model.User, uint64, uint64) {
	totalElements := uint64(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return []model.User{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *UserService) findUserForPicture(id string, justification ExtractPictureJustification, userID string, isAdmin bool) (model.User, error) {
	user, err := svc.userRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if id == userID || isAdmin {
		return user, nil
	}
	if justification.OrganizationID == nil && justification.GroupID == nil && justification.InvitationID == nil {
		return nil, errorpkg.NewPictureNotFoundError(nil)
	}
	if justification.OrganizationID != nil {
		org, err := svc.orgCache.Get(*justification.OrganizationID)
		if err != nil {
			return nil, err
		}
		if err := svc.orgGuard.Authorize(userID, org, model.PermissionViewer); err != nil {
			return nil, err
		}
		if !slices.Contains(org.GetMembers(), id) {
			return nil, errorpkg.NewPictureNotFoundError(nil)
		}
	} else if justification.GroupID != nil {
		group, err := svc.groupCache.Get(*justification.GroupID)
		if err != nil {
			return nil, err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return nil, err
		}
		if !slices.Contains(group.GetMembers(), id) {
			return nil, errorpkg.NewPictureNotFoundError(nil)
		}
	} else if justification.InvitationID != nil {
		invitation, err := svc.invitationRepo.Find(*justification.InvitationID)
		if err != nil {
			return nil, err
		}
		if invitation.GetOwnerID() != id {
			return nil, errorpkg.NewPictureNotFoundError(nil)
		}
	}
	return user, nil
}

type userMapper struct{}

func newUserMapper() *userMapper {
	return &userMapper{}
}

func (mp *userMapper) mapOne(user model.User) *dto.User {
	res := &dto.User{
		ID:         user.GetID(),
		FullName:   user.GetFullName(),
		Email:      user.GetEmail(),
		Username:   user.GetUsername(),
		CreateTime: user.GetCreateTime(),
		UpdateTime: user.GetUpdateTime(),
	}
	if user.GetPicture() != nil {
		res.Picture = &dto.Picture{
			Extension: helper.Base64ToExtension(*user.GetPicture()),
		}
	}
	return res
}

func (mp *userMapper) mapMany(users []model.User) ([]*dto.User, error) {
	res := make([]*dto.User, 0)
	for _, user := range users {
		res = append(res, mp.mapOne(user))
	}
	return res, nil
}
