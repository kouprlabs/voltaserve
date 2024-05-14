package service

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
	workspaceRepo  repo.WorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	workspaceGuard *guard.WorkspaceGuard
	fileRepo       repo.FileRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	storageMapper  *storageMapper
	userRepo       repo.UserRepo
}

type NewStorageServiceOptions struct {
	WorkspaceRepo repo.WorkspaceRepo
	FileRepo      repo.FileRepo
	UserRepo      repo.UserRepo
}

func NewStorageService(opts NewStorageServiceOptions) *StorageService {
	svc := &StorageService{
		workspaceCache: cache.NewWorkspaceCache(),
		workspaceGuard: guard.NewWorkspaceGuard(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		storageMapper:  newStorageMapper(),
	}
	if opts.WorkspaceRepo != nil {
		svc.workspaceRepo = opts.WorkspaceRepo
	} else {
		svc.workspaceRepo = repo.NewWorkspaceRepo()
	}
	if opts.FileRepo != nil {
		svc.fileRepo = opts.FileRepo
	} else {
		svc.fileRepo = repo.NewFileRepo()
	}
	if opts.UserRepo != nil {
		svc.userRepo = opts.UserRepo
	} else {
		svc.userRepo = repo.NewUserRepo()
	}
	return svc
}

func (svc *StorageService) GetAccountUsage(userID string) (*StorageUsage, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	ids, err := svc.workspaceRepo.GetIDs()
	if err != nil {
		return nil, err
	}
	workspaces := []model.Workspace{}
	for _, id := range ids {
		var workspace model.Workspace
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
		root, err := svc.fileCache.Get(w.GetRootID())
		if err != nil {
			return nil, err
		}
		size, err := svc.fileRepo.GetSize(root.GetID())
		if err != nil {
			return nil, err
		}
		b = b + size
		maxBytes = maxBytes + w.GetStorageCapacity()
	}
	return svc.storageMapper.mapStorageUsage(b, maxBytes), nil
}

func (svc *StorageService) GetWorkspaceUsage(workspaceID string, userID string) (*StorageUsage, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(workspaceID)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(user, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	root, err := svc.fileCache.Get(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, root, model.PermissionViewer); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.GetSize(root.GetID())
	if err != nil {
		return nil, err
	}
	return svc.storageMapper.mapStorageUsage(size, workspace.GetStorageCapacity()), nil
}

func (svc *StorageService) GetFileUsage(fileID string, userID string) (*StorageUsage, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	size, err := svc.fileRepo.GetSize(file.GetID())
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
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
