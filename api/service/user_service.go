package service

import (
	"sort"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/guard"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type User struct {
	ID         string  `json:"id"`
	FullName   string  `json:"fullName"`
	Picture    *string `json:"picture,omitempty"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime"`
}

type UserList struct {
	Data          []*User `json:"data"`
	TotalPages    uint    `json:"totalPages"`
	TotalElements uint    `json:"totalElements"`
	Page          uint    `json:"page"`
	Size          uint    `json:"size"`
}

type UserListOptions struct {
	Query               string
	OrganizationID      string
	GroupID             string
	NonGroupMembersOnly bool
	SortBy              string
	SortOrder           string
	Page                uint
	Size                uint
}

type UserService struct {
	userRepo   repo.UserRepo
	userMapper *userMapper
	userSearch *search.UserSearch
	orgRepo    repo.OrganizationRepo
	orgCache   *cache.OrganizationCache
	orgGuard   *guard.OrganizationGuard
	groupRepo  repo.GroupRepo
	groupGuard *guard.GroupGuard
	groupCache *cache.GroupCache
	config     config.Config
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:   repo.NewUserRepo(),
		userMapper: newUserMapper(),
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

func (svc *UserService) List(opts UserListOptions, userID string) (*UserList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	if opts.OrganizationID == "" && opts.GroupID == "" {
		return &UserList{
			Data:          []*User{},
			TotalPages:    1,
			TotalElements: 0,
			Page:          1,
			Size:          0,
		}, nil
	}
	var org model.Organization
	if opts.OrganizationID != "" {
		org, err = svc.orgCache.Get(opts.OrganizationID)
		if err != nil {
			return nil, err
		}
		if err := svc.orgGuard.Authorize(user, org, model.PermissionViewer); err != nil {
			return nil, err
		}
	}
	var group model.Group
	if opts.GroupID != "" {
		group, err = svc.groupCache.Get(opts.GroupID)
		if err != nil {
			return nil, err
		}
		if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
			return nil, err
		}
	}
	res := []model.User{}
	if opts.Query == "" {
		if opts.OrganizationID != "" && opts.GroupID != "" && opts.NonGroupMembersOnly {
			orgMembers, err := svc.orgRepo.GetMembers(group.GetOrganizationID())
			if err != nil {
				return nil, err
			}
			groupMembers, err := svc.groupRepo.GetMembers(opts.GroupID)
			if err != nil {
				return nil, err
			}
			for _, om := range orgMembers {
				found := false
				for _, tm := range groupMembers {
					if om.GetID() == tm.GetID() {
						found = true
						break
					}
				}
				if !found {
					res = append(res, om)
				}
			}
		} else if opts.OrganizationID != "" {
			res, err = svc.orgRepo.GetMembers(opts.OrganizationID)
			if err != nil {
				return nil, err
			}
		} else if opts.GroupID != "" {
			res, err = svc.groupRepo.GetMembers(opts.GroupID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		users, err := svc.userSearch.Query(opts.Query)
		if err != nil {
			return nil, err
		}
		var members []model.User
		if opts.OrganizationID != "" {
			members, err = svc.orgRepo.GetMembers(opts.OrganizationID)
			if err != nil {
				return nil, err
			}
		} else if opts.GroupID != "" {
			members, err = svc.groupRepo.GetMembers(opts.GroupID)
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
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(res, opts.SortBy, opts.SortOrder, userID)
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
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *UserService) doSorting(data []model.User, sortBy string, sortOrder string, userID string) []model.User {
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

func (svc *UserService) doPagination(data []model.User, page, size uint) ([]model.User, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return nil, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

type userMapper struct {
}

func newUserMapper() *userMapper {
	return &userMapper{}
}

func (mp *userMapper) mapOne(user model.User) *User {
	return &User{
		ID:         user.GetID(),
		FullName:   user.GetFullName(),
		Picture:    user.GetPicture(),
		Email:      user.GetEmail(),
		Username:   user.GetUsername(),
		CreateTime: user.GetCreateTime(),
		UpdateTime: user.GetUpdateTime(),
	}
}

func (mp *userMapper) mapMany(users []model.User) ([]*User, error) {
	res := []*User{}
	for _, u := range users {
		res = append(res, mp.mapOne(u))
	}
	return res, nil
}
