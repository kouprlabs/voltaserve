// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package service

import (
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"slices"
	"sort"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type UserService struct {
	userMapper *userMapper
	userRepo   repo.UserRepo
	userSearch *search.UserSearch
	orgRepo    repo.OrganizationRepo
	orgCache   *cache.OrganizationCache
	orgGuard   *guard.OrganizationGuard
	groupRepo  repo.GroupRepo
	groupGuard *guard.GroupGuard
	groupCache *cache.GroupCache
	config     *config.Config
}

func NewUserService() *UserService {
	return &UserService{
		userMapper: newUserMapper(),
		userRepo:   repo.NewUserRepo(),
		userSearch: search.NewUserSearch(),
		orgRepo:    repo.NewOrganizationRepo(),
		orgCache:   cache.NewOrganizationCache(),
		orgGuard:   guard.NewOrganizationGuard(),
		groupRepo:  repo.NewGroupRepo(),
		groupGuard: guard.NewGroupGuard(),
		groupCache: cache.NewGroupCache(),
		config:     config.GetConfig(),
	}
}

type User struct {
	ID         string   `json:"id"`
	FullName   string   `json:"fullName"`
	Picture    *Picture `json:"picture,omitempty"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
	CreateTime string   `json:"createTime"`
	UpdateTime *string  `json:"updateTime"`
}

type Picture struct {
	Extension string `json:"extension"`
}

type UserListOptions struct {
	Query               string
	OrganizationID      string
	GroupID             string
	ExcludeGroupMembers bool
	SortBy              string
	SortOrder           string
	Page                int64
	Size                int64
}

type UserList struct {
	Data          []*User `json:"data"`
	TotalPages    int64   `json:"totalPages"`
	TotalElements int64   `json:"totalElements"`
	Page          int64   `json:"page"`
	Size          int64   `json:"size"`
}

func (svc *UserService) List(opts UserListOptions, userID string) (*UserList, error) {
	users, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(users, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.userMapper.mapMany(paged)
	if err != nil {
		return nil, err
	}
	return &UserList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          int64(len(mapped)),
	}, nil
}

type UserProbe struct {
	TotalPages    int64 `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}

func (svc *UserService) Probe(opts UserListOptions, userID string) (*UserProbe, error) {
	users, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := int64(len(users))
	return &UserProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
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
	res := make([]model.User, 0)
	var err error
	if opts.Query == "" {
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
	} else {
		count, err := svc.userRepo.Count()
		if err != nil {
			return nil, err
		}
		users, err := svc.userSearch.Query(opts.Query, infra.QueryOptions{Limit: count})
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
	}
	return res, nil
}

func (svc *UserService) doSorting(data []model.User, sortBy string, sortOrder string) []model.User {
	if sortBy == SortByEmail {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetEmail() > data[j].GetEmail()
			} else {
				return data[i].GetEmail() < data[j].GetEmail()
			}
		})
		return data
	} else if sortBy == SortByFullName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetFullName() > data[j].GetFullName()
			} else {
				return data[i].GetFullName() < data[j].GetFullName()
			}
		})
		return data
	}
	return data
}

func (svc *UserService) doPagination(data []model.User, page, size int64) ([]model.User, int64, int64) {
	totalElements := int64(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return []model.User{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

func (svc *UserService) ExtractPicture(id string, organizationID *string, groupID *string, userID string, isAdmin bool) ([]byte, *string, *string, error) {
	user, err := svc.find(id, organizationID, groupID, userID, isAdmin)
	if err != nil {
		return nil, nil, nil, err
	}
	if user.GetPicture() == nil {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
	mime := helper.Base64ToMIME(*user.GetPicture())
	ext := helper.Base64ToExtension(*user.GetPicture())
	b, err := helper.Base64ToBytes(*user.GetPicture())
	if err != nil {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
	return b, &ext, &mime, nil
}

func (svc *UserService) find(id string, orgID *string, groupID *string, userID string, isAdmin bool) (model.User, error) {
	user, err := svc.userRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if id == userID || isAdmin {
		return user, nil
	}
	if orgID == nil && groupID == nil {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
	if orgID != nil {
		org, err := svc.orgCache.Get(*orgID)
		if err != nil {
			return nil, err
		}
		if err := svc.orgGuard.Authorize(userID, org, model.PermissionViewer); err != nil {
			return nil, err
		}
		if !slices.Contains(org.GetMembers(), id) {
			return nil, errorpkg.NewS3ObjectNotFoundError(nil)
		}
	} else {
		group, err := svc.groupCache.Get(*groupID)
		if err != nil {
			return nil, err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return nil, err
		}
		if !slices.Contains(group.GetMembers(), id) {
			return nil, errorpkg.NewS3ObjectNotFoundError(nil)
		}
	}
	return user, nil
}

type userMapper struct{}

func newUserMapper() *userMapper {
	return &userMapper{}
}

func (mp *userMapper) mapOne(user model.User) *User {
	res := &User{
		ID:         user.GetID(),
		FullName:   user.GetFullName(),
		Email:      user.GetEmail(),
		Username:   user.GetUsername(),
		CreateTime: user.GetCreateTime(),
		UpdateTime: user.GetUpdateTime(),
	}
	if user.GetPicture() != nil {
		res.Picture = &Picture{
			Extension: helper.Base64ToExtension(*user.GetPicture()),
		}
	}
	return res
}

func (mp *userMapper) mapMany(users []model.User) ([]*User, error) {
	res := make([]*User, 0)
	for _, user := range users {
		res = append(res, mp.mapOne(user))
	}
	return res, nil
}
