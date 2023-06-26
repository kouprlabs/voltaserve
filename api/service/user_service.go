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
	Query     string
	OrgID     string
	GroupID   string
	SortBy    string
	SortOrder string
	Page      uint
	Size      uint
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
	if opts.OrgID == "" && opts.GroupID == "" {
		return &UserList{
			Data:          []*User{},
			TotalPages:    1,
			TotalElements: 0,
			Page:          1,
			Size:          0,
		}, nil
	}
	var org model.Organization
	if opts.OrgID != "" {
		org, err = svc.orgCache.Get(opts.OrgID)
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
		if opts.OrgID != "" {
			res, err = svc.orgRepo.GetMembers(opts.OrgID)
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
		if opts.OrgID != "" {
			members, err = svc.orgRepo.GetMembers(opts.OrgID)
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
	sorted := svc.doSorting(res, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.doPaging(sorted, opts.Page, opts.Size)
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

func (svc *UserService) doPaging(data []model.User, page uint, size uint) (pageData []model.User, totalElements uint, totalPages uint) {
	page = page - 1
	low := size * page
	high := low + size
	if low >= uint(len(data)) {
		pageData = []model.User{}
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
