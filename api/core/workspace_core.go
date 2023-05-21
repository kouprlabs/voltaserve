package core

import (
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	"github.com/google/uuid"
)

type Workspace struct {
	Id              string       `json:"id"`
	Image           *string      `json:"image,omitempty"`
	Name            string       `json:"name"`
	RootId          string       `json:"rootId,omitempty"`
	StorageCapacity int64        `json:"storageCapacity"`
	Permission      string       `json:"permission"`
	Organization    Organization `json:"organization"`
	CreateTime      string       `json:"createTime"`
	UpdateTime      *string      `json:"updateTime,omitempty"`
}

type WorkspaceSearchOptions struct {
	Text string `json:"text" validate:"required"`
}

type CreateWorkspaceOptions struct {
	Name            string  `json:"name" validate:"required,max=255"`
	Image           *string `json:"image"`
	OrganizationId  string  `json:"organizationId" validate:"required"`
	StorageCapacity int64   `json:"storageCapacity" validate:"required,min=1"`
}

type UpdateWorkspaceNameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateWorkspaceStorageCapacityOptions struct {
	StorageCapacity int64 `json:"storageCapacity" validate:"required,min=1"`
}

type WorkspaceService struct {
	workspaceRepo   *repo.WorkspaceRepo
	workspaceCache  *cache.WorkspaceCache
	workspaceGuard  *guard.WorkspaceGuard
	workspaceSearch *search.WorkspaceSearch
	workspaceMapper *workspaceMapper
	fileRepo        *repo.FileRepo
	fileCache       *cache.FileCache
	fileGuard       *guard.FileGuard
	fileMapper      *FileMapper
	userRepo        *repo.UserRepo
	imageProc       *infra.ImageProcessor
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
		imageProc:       infra.NewImageProcessor(),
		s3:              infra.NewS3Manager(),
		config:          config.GetConfig(),
	}
}

func (svc *WorkspaceService) Create(req CreateWorkspaceOptions, userId string) (*Workspace, error) {
	id := helpers.NewId()
	bucket := strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := svc.s3.CreateBucket(bucket); err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceRepo.Insert(repo.WorkspaceInsertOptions{
		Id:              id,
		Name:            req.Name,
		StorageCapacity: req.StorageCapacity,
		OrganizationId:  req.OrganizationId,
		Image:           req.Image,
		Bucket:          bucket,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceRepo.GrantUserPermission(workspace.GetId(), userId, model.PermissionOwner); err != nil {
		return nil, err
	}
	workspace, err = svc.workspaceRepo.FindByID(workspace.GetId())
	if err != nil {
		return nil, err
	}
	root, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        "root",
		WorkspaceId: workspace.GetId(),
		Type:        model.FileTypeFolder,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(root.GetId(), userId, model.PermissionOwner); err != nil {
		return nil, err
	}
	if err = svc.workspaceRepo.UpdateRootId(workspace.GetId(), root.GetId()); err != nil {
		return nil, err
	}
	if workspace, err = svc.workspaceRepo.FindByID(workspace.GetId()); err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Index([]model.WorkspaceModel{workspace}); err != nil {
		return nil, err
	}
	if root, err = svc.fileRepo.Find(root.GetId()); err != nil {
		return nil, err
	}
	if err := svc.fileCache.Set(root); err != nil {
		return nil, err
	}
	if err = svc.workspaceCache.Set(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapWorkspace(workspace, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Find(id string, userId string) (*Workspace, error) {
	user, err := svc.userRepo.Find(userId)
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
	res, err := svc.workspaceMapper.mapWorkspace(workspace, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) FindAll(userId string) ([]*Workspace, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	ids, err := svc.workspaceRepo.GetIds()
	if err != nil {
		return nil, err
	}
	res := []*Workspace{}
	for _, id := range ids {
		var workspace model.WorkspaceModel
		workspace, err = svc.workspaceCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.workspaceGuard.IsAuthorized(user, workspace, model.PermissionViewer) {
			dto, err := svc.workspaceMapper.mapWorkspace(workspace, userId)
			if err != nil {
				return nil, err
			}
			res = append(res, dto)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) Search(query string, userId string) ([]*Workspace, error) {
	workspaces, err := svc.workspaceSearch.Query(query)
	if err != nil {
		return nil, err
	}
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	res := []*Workspace{}
	for _, w := range workspaces {
		if svc.workspaceGuard.IsAuthorized(user, w, model.PermissionViewer) {
			dto, err := svc.workspaceMapper.mapWorkspace(w, userId)
			if err != nil {
				return nil, err
			}
			res = append(res, dto)
		}
	}
	return res, nil
}

func (svc *WorkspaceService) UpdateName(id string, name string, userId string) (*Workspace, error) {
	user, err := svc.userRepo.Find(userId)
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
	if err = svc.workspaceSearch.Update([]model.WorkspaceModel{workspace}); err != nil {
		return nil, err
	}
	if err = svc.workspaceCache.Set(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapWorkspace(workspace, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) UpdateStorageCapacity(id string, storageCapacity int64, userId string) (*Workspace, error) {
	user, err := svc.userRepo.Find(userId)
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
	size, err := svc.fileRepo.GetSize(workspace.GetRootId())
	if err != nil {
		return nil, err
	}
	if storageCapacity < size {
		return nil, errorpkg.NewInsufficientStorageCapacityError()
	}
	if workspace, err = svc.workspaceRepo.UpdateStorageCapacity(id, storageCapacity); err != nil {
		return nil, err
	}
	if err = svc.workspaceSearch.Update([]model.WorkspaceModel{workspace}); err != nil {
		return nil, err
	}
	if err = svc.workspaceCache.Set(workspace); err != nil {
		return nil, err
	}
	res, err := svc.workspaceMapper.mapWorkspace(workspace, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *WorkspaceService) Delete(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
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
	if workspace, err = svc.workspaceRepo.FindByID(id); err != nil {
		return err
	}
	if err = svc.workspaceRepo.Delete(id); err != nil {
		return err
	}
	if err = svc.workspaceSearch.Delete([]string{workspace.GetId()}); err != nil {
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
	workspace, err := svc.workspaceRepo.FindByID(id)
	if err != nil {
		return false, err
	}
	root, err := svc.fileRepo.Find(workspace.GetRootId())
	if err != nil {
		return false, err
	}
	usage, err := svc.fileRepo.GetSize(root.GetId())
	if err != nil {
		return false, err
	}
	expectedUsage := usage + byteSize
	if expectedUsage > workspace.GetStorageCapacity() {
		return false, err
	}
	return true, nil
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

func (mp *workspaceMapper) mapWorkspace(m model.WorkspaceModel, userId string) (*Workspace, error) {
	org, err := mp.orgCache.Get(m.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	v, err := mp.orgMapper.mapOrganization(org, userId)
	if err != nil {
		return nil, err
	}
	res := &Workspace{
		Id:              m.GetId(),
		Name:            m.GetName(),
		RootId:          m.GetRootId(),
		StorageCapacity: m.GetStorageCapacity(),
		Organization:    *v,
		CreateTime:      m.GetCreateTime(),
		UpdateTime:      m.GetUpdateTime(),
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
