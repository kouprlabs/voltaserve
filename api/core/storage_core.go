package core

import (
	"voltaserve/cache"
	"voltaserve/guard"
	"voltaserve/model"
	"voltaserve/repo"
)

type StorageUsage struct {
	Bytes      int64 `json:"bytes"`
	MaxBytes   int64 `json:"maxBytes"`
	Percentage int   `json:"percentage"`
}

type StorageService struct {
	workspaceRepo  repo.CoreWorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	workspaceGuard *guard.WorkspaceGuard
	fileRepo       repo.CoreFileRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	storageMapper  *storageMapper
	userRepo       repo.CoreUserRepo
}

func NewStorageService() *StorageService {
	return &StorageService{
		workspaceRepo:  repo.NewPostgresWorkspaceRepo(),
		workspaceCache: cache.NewWorkspaceCache(),
		workspaceGuard: guard.NewWorkspaceGuard(),
		fileRepo:       repo.NewPostgresFileRepo(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		storageMapper:  newStorageMapper(),
		userRepo:       repo.NewPostgresUserRepo(),
	}
}

func (svc *StorageService) GetAccountUsage(userId string) (*StorageUsage, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	ids, err := svc.workspaceRepo.GetIds()
	if err != nil {
		return nil, err
	}
	workspaces := []model.WorkspaceModel{}
	for _, id := range ids {
		var workspace model.WorkspaceModel
		workspace, err = svc.workspaceCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.workspaceGuard.IsAuthorized(user, workspace, model.PermissionOwner) {
			workspaces = append(workspaces, workspace)
		}
	}
	var maxBytes int64 = 0
	var b int64 = 0
	for _, w := range workspaces {
		root, err := svc.fileCache.Get(w.GetRootId())
		if err != nil {
			return nil, err
		}
		size, err := svc.fileRepo.GetSize(root.GetId())
		if err != nil {
			return nil, err
		}
		b = b + size
		maxBytes = maxBytes + w.GetStorageCapacity()
	}
	return svc.storageMapper.mapStorageUsage(b, maxBytes), nil
}

func (svc *StorageService) GetWorkspaceUsage(workspaceId string, userId string) (*StorageUsage, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(workspaceId)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(user, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	root, err := svc.fileCache.Get(workspace.GetRootId())
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, root, model.PermissionViewer); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.GetSize(root.GetId())
	if err != nil {
		return nil, err
	}
	return svc.storageMapper.mapStorageUsage(size, workspace.GetStorageCapacity()), nil
}

func (svc *StorageService) GetFileUsage(fileId string, userId string) (*StorageUsage, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileId)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.GetSize(file.GetId())
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceId())
	if err != nil {
		return nil, err
	}
	return svc.storageMapper.mapStorageUsage(size, workspace.GetStorageCapacity()), nil
}

type storageMapper struct {
}

func newStorageMapper() *storageMapper {
	return &storageMapper{}
}

func (mp *storageMapper) mapStorageUsage(byteCount int64, maxBytes int64) *StorageUsage {
	res := StorageUsage{
		Bytes:    byteCount,
		MaxBytes: maxBytes,
	}
	if maxBytes != 0 {
		res.Percentage = int(byteCount * 100 / maxBytes)
	}
	return &res
}
