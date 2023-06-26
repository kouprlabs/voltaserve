package service

import (
	"sort"
	"strings"
	"time"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	"github.com/google/uuid"
)

type Workspace struct {
	ID              string       `json:"id"`
	Image           *string      `json:"image,omitempty"`
	Name            string       `json:"name"`
	RootID          string       `json:"rootId,omitempty"`
	StorageCapacity int64        `json:"storageCapacity"`
	Permission      string       `json:"permission"`
	Organization    Organization `json:"organization"`
	CreateTime      string       `json:"createTime"`
	UpdateTime      *string      `json:"updateTime,omitempty"`
}

type WorkspaceList struct {
	Data          []*Workspace `json:"data"`
	TotalPages    uint         `json:"totalPages"`
	TotalElements uint         `json:"totalElements"`
	Page          uint         `json:"page"`
	Size          uint         `json:"size"`
}

type WorkspaceSearchOptions struct {
	Text string `json:"text" validate:"required"`
}

type WorkspaceCreateOptions struct {
	Name            string  `json:"name" validate:"required,max=255"`
	Image           *string `json:"image"`
	OrganizationID  string  `json:"organizationId" validate:"required"`
	StorageCapacity int64   `json:"storageCapacity" validate:"required,min=1"`
}

type WorkspaceUpdateNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type WorkspaceUpdateStorageCapacityOptions struct {
	StorageCapacity int64 `json:"storageCapacity" validate:"required,min=1"`
}

type WorkspaceService struct {
	workspaceRepo   repo.WorkspaceRepo
	workspaceCache  *cache.WorkspaceCache
	workspaceGuard  *guard.WorkspaceGuard
	workspaceSearch *search.WorkspaceSearch
	workspaceMapper *workspaceMapper
	fileRepo        repo.FileRepo
	fileCache       *cache.FileCache
	fileGuard       *guard.FileGuard
	fileMapper      *FileMapper
	userRepo        repo.UserRepo
	s3              *infra.S3Manager
	config          config.Config
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo:   repo.NewWorkspaceRepo(),
		workspaceCache:  cache.NewWorkspaceCache(),
		workspaceSearch: search.NewWorkspaceSearch(),
		workspaceGuard:  guard.NewWorkspaceGuard(),
		workspaceMapper: newWorkspaceMapper(),
		fileRepo:        repo.NewFileRepo(),
		fileCache:       cache.NewFileCache(),
		fileGuard:       guard.NewFileGuard(),
		fileMapper:      NewFileMapper(),
		userRepo:        repo.NewUserRepo(),
		s3:              infra.NewS3Manager(),
		config:          config.GetConfig(),
	}
}

func (svc *WorkspaceService) Create(opts WorkspaceCreateOptions, userID string) (*Workspace, error) {
	id := helper.NewID()
	bucket := strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := svc.s3.CreateBucket(bucket); err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceRepo.Insert(repo.WorkspaceInsertOptions{
		ID:              id,
		Name:            opts.Name,
		StorageCapacity: opts.StorageCapacity,
		OrganizationID:  opts.OrganizationID,
		Image:           opts.Image,
		Bucket:          bucket,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceRepo.GrantUserPermission(workspace.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	workspace, err = svc.workspaceRepo.Find(workspace.GetID())
	if err != nil {
		return nil, err
	}
	root, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        "root",
		WorkspaceID: workspace.GetID(),
		Type:        model.FileTypeFolder,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(root.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	if err = svc.workspaceRepo.UpdateRootID(workspace.GetID(), root.GetID()); err != nil {
		return nil, err
	}
	if workspace, err = svc.workspaceRepo.Find(workspace.GetID()); err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Index([]model.Workspace{workspace}); err != nil {
		return nil, err
	}
	if root, err = svc.fileRepo.Find(root.GetID()); err != nil {
		return nil, err
	}
	if err := svc.fileCache.Set(root); err != nil {
		return nil, err
	}
	if err = svc.workspaceCache.Set(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Find(id string, userID string) (*Workspace, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(user, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) List(page uint, size uint, sortBy string, sortOrder string, userID string) (*WorkspaceList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	ids, err := svc.workspaceRepo.GetIDs()
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorizationByIDs(ids, user)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(authorized, sortBy, sortOrder, userID)
	paged, totalElements, totalPages := svc.doPaging(sorted, page, size)
	mapped, err := svc.workspaceMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &WorkspaceList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *WorkspaceService) Search(query string, page uint, size uint, userID string) (*WorkspaceList, error) {
	workspaces, err := svc.workspaceSearch.Query(query)
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
	mapped, err := svc.workspaceMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &WorkspaceList{
		Data:          mapped,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *WorkspaceService) UpdateName(id string, name string, userID string) (*Workspace, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(user, workspace, model.PermissionEditor); err != nil {
		return nil, err
	}
	if workspace, err = svc.workspaceRepo.UpdateName(id, name); err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Update([]model.Workspace{workspace}); err != nil {
		return nil, err
	}
	if err = svc.workspaceCache.Set(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) UpdateStorageCapacity(id string, storageCapacity int64, userID string) (*Workspace, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(user, workspace, model.PermissionEditor); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.GetSize(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	if storageCapacity < size {
		return nil, errorpkg.NewInsufficientStorageCapacityError()
	}
	if workspace, err = svc.workspaceRepo.UpdateStorageCapacity(id, storageCapacity); err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Update([]model.Workspace{workspace}); err != nil {
		return nil, err
	}
	if err = svc.workspaceCache.Set(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapOne(workspace, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Delete(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.workspaceGuard.Authorize(user, workspace, model.PermissionOwner); err != nil {
		return err
	}
	if workspace, err = svc.workspaceRepo.Find(id); err != nil {
		return err
	}
	if err = svc.workspaceRepo.Delete(id); err != nil {
		return err
	}
	if err = svc.workspaceSearch.Delete([]string{workspace.GetID()}); err != nil {
		return err
	}
	if err = svc.workspaceCache.Delete(id); err != nil {
		return err
	}
	if err = svc.s3.RemoveBucket(workspace.GetBucket()); err != nil {
		return err
	}
	return nil
}

func (svc *WorkspaceService) HasEnoughSpaceForByteSize(id string, byteSize int64) (bool, error) {
	workspace, err := svc.workspaceRepo.Find(id)
	if err != nil {
		return false, err
	}
	root, err := svc.fileRepo.Find(workspace.GetRootID())
	if err != nil {
		return false, err
	}
	usage, err := svc.fileRepo.GetSize(root.GetID())
	if err != nil {
		return false, err
	}
	expectedUsage := usage + byteSize
	if expectedUsage > workspace.GetStorageCapacity() {
		return false, err
	}
	return true, nil
}

func (svc *WorkspaceService) findAll(userID string) ([]*Workspace, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	ids, err := svc.workspaceRepo.GetIDs()
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorizationByIDs(ids, user)
	if err != nil {
		return nil, err
	}
	mapped, err := svc.workspaceMapper.mapMany(authorized, userID)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (svc *WorkspaceService) doAuthorization(data []model.Workspace, user model.User) ([]model.Workspace, error) {
	var res []model.Workspace
	for _, w := range data {
		if svc.workspaceGuard.IsAuthorized(user, w, model.PermissionViewer) {
			res = append(res, w)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) doAuthorizationByIDs(ids []string, user model.User) ([]model.Workspace, error) {
	var res []model.Workspace
	for _, id := range ids {
		var w model.Workspace
		w, err := svc.workspaceCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.workspaceGuard.IsAuthorized(user, w, model.PermissionViewer) {
			res = append(res, w)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) doSorting(data []model.Workspace, sortBy string, sortOrder string, userID string) []model.Workspace {
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

func (svc *WorkspaceService) doPaging(data []model.Workspace, page uint, size uint) (pageData []model.Workspace, totalElements uint, totalPages uint) {
	page = page - 1
	low := size * page
	high := low + size
	if low >= uint(len(data)) {
		pageData = []model.Workspace{}
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

type workspaceMapper struct {
	orgCache   *cache.OrganizationCache
	orgMapper  *organizationMapper
	groupCache *cache.GroupCache
}

func newWorkspaceMapper() *workspaceMapper {
	return &workspaceMapper{
		orgCache:   cache.NewOrganizationCache(),
		orgMapper:  newOrganizationMapper(),
		groupCache: cache.NewGroupCache(),
	}
}

func (mp *workspaceMapper) mapOne(m model.Workspace, userID string) (*Workspace, error) {
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	v, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	res := &Workspace{
		ID:              m.GetID(),
		Name:            m.GetName(),
		RootID:          m.GetRootID(),
		StorageCapacity: m.GetStorageCapacity(),
		Organization:    *v,
		CreateTime:      m.GetCreateTime(),
		UpdateTime:      m.GetUpdateTime(),
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

func (mp *workspaceMapper) mapMany(workspaces []model.Workspace, userID string) ([]*Workspace, error) {
	res := make([]*Workspace, 0)
	for _, f := range workspaces {
		v, err := mp.mapOne(f, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
