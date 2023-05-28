package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/conversion"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	"github.com/reactivex/rxgo/v2"
)

type File struct {
	ID          string      `json:"id"`
	WorkspaceId string      `json:"workspaceId"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	ParentId    *string     `json:"parentId,omitempty"`
	Version     *int64      `json:"version,omitempty"`
	Original    *Download   `json:"original,omitempty"`
	Preview     *Download   `json:"preview,omitempty"`
	Thumbnail   *Thumbnail  `json:"thumbnail,omitempty"`
	Snapshots   []*Snapshot `json:"snapshots,omitempty"`
	Permission  string      `json:"permission"`
	IsShared    bool        `json:"isShared"`
	CreateTime  string      `json:"createTime"`
	UpdateTime  *string     `json:"updateTime,omitempty"`
}

type FileList struct {
	Data          []*File `json:"data"`
	TotalPages    uint    `json:"totalPages"`
	TotalElements uint    `json:"totalElements"`
	Page          uint    `json:"page"`
	Size          uint    `json:"size"`
}

type FileSearchOptions struct {
	Text             string  `json:"text" validate:"required"`
	WorkspaceId      string  `json:"workspaceId" validate:"required"`
	ParentId         *string `json:"parentId,omitempty"`
	Type             *string `json:"type,omitempty" validate:"omitempty,oneof=file folder"`
	CreateTimeAfter  *int64  `json:"createTimeAfter,omitempty"`
	CreateTimeBefore *int64  `json:"createTimeBefore,omitempty"`
	UpdateTimeAfter  *int64  `json:"updateTimeAfter,omitempty"`
	UpdateTimeBefore *int64  `json:"updateTimeBefore,omitempty"`
}

type FileSearchResult struct {
	Query         FileSearchOptions `json:"request"`
	Data          []*File           `json:"data"`
	TotalPages    uint              `json:"totalPages"`
	TotalElements uint              `json:"totalElements"`
	Page          uint              `json:"page"`
	Size          uint              `json:"size"`
}

type FileCreateOptions struct {
	WorkspaceId string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	Type        string  `json:"type" validate:"required,oneof=file folder"`
	ParentId    *string `json:"parentId" validate:"required"`
}

type FileCreateFolderOptions struct {
	WorkspaceId string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	ParentId    *string `json:"parentId"`
}

type FileCopyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileBatchDeleteOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileBatchGetOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileGrantUserPermissionOptions struct {
	UserID     string   `json:"userId" validate:"required"`
	IDs        []string `json:"ids" validate:"required"`
	Permission string   `json:"permission" validate:"required,oneof=viewer editor owner"`
}

type FileRevokeUserPermissionOptions struct {
	IDs    []string `json:"ids" validate:"required"`
	UserID string   `json:"userId" validate:"required"`
}

type FileGrantGroupPermissionOptions struct {
	GroupID    string   `json:"groupId" validate:"required"`
	IDs        []string `json:"ids" validate:"required"`
	Permission string   `json:"permission" validate:"required,oneof=viewer editor owner"`
}

type FileRevokeGroupPermissionOptions struct {
	IDs     []string `json:"ids" validate:"required"`
	GroupID string   `json:"groupId" validate:"required"`
}

type FileMoveOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileRenameOptions struct {
	Name string `json:"name" validate:"required,max=255"`
}

type Snapshot struct {
	ID        string     `json:"id"`
	Version   int64      `json:"version"`
	Original  *Download  `json:"original,omitempty"`
	Preview   *Download  `json:"preview,omitempty"`
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
}

type ImageProps struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Thumbnail struct {
	Base64 string `json:"base64"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Download struct {
	Extension string      `json:"extension"`
	Size      int64       `json:"size"`
	Image     *ImageProps `json:"image,omitempty"`
}

type UserPermission struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	Permission string `json:"permission"`
}

type GroupPermission struct {
	ID         string `json:"id"`
	Group      *Group `json:"group"`
	Permission string `json:"permission"`
}

const SortByName = "name"
const SortByKind = "kind"
const SortBySize = "size"
const SortByDateCreated = "date_created"
const SortByDateModified = "date_modified"

const SortOrderAsc = "asc"
const SortOrderDesc = "desc"

func IsValidSortBy(value string) bool {
	return value == "" || value == SortByName || value == SortByKind || value == SortBySize || value == SortByDateCreated || value == SortByDateModified
}

func IsValidSortOrder(value string) bool {
	return value == "" || value == SortOrderAsc || value == SortOrderDesc
}

type FileService struct {
	fileRepo       repo.FileRepo
	fileSearch     *search.FileSearch
	fileGuard      *guard.FileGuard
	fileMapper     *FileMapper
	fileCache      *cache.FileCache
	workspaceCache *cache.WorkspaceCache
	workspaceRepo  repo.WorkspaceRepo
	workspaceGuard *guard.WorkspaceGuard
	workspaceSvc   *WorkspaceService
	snapshotRepo   repo.SnapshotRepo
	userRepo       repo.UserRepo
	userMapper     *userMapper
	groupCache     *cache.GroupCache
	groupGuard     *guard.GroupGuard
	groupMapper    *groupMapper
	permissionRepo repo.PermissionRepo
	fileIdentifier *infra.FileIdentifier
	s3             *infra.S3Manager
	pdfPipeline    conversion.Pipeline
	imagePipeline  conversion.Pipeline
	officePipeline conversion.Pipeline
	videoPipeline  conversion.Pipeline
	config         config.Config
}

func NewFileService() *FileService {
	return &FileService{
		fileRepo:       repo.NewFileRepo(),
		fileCache:      cache.NewFileCache(),
		fileSearch:     search.NewFileSearch(),
		fileGuard:      guard.NewFileGuard(),
		fileMapper:     NewFileMapper(),
		workspaceGuard: guard.NewWorkspaceGuard(),
		workspaceCache: cache.NewWorkspaceCache(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceSvc:   NewWorkspaceService(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		userRepo:       repo.NewUserRepo(),
		userMapper:     newUserMapper(),
		groupCache:     cache.NewGroupCache(),
		groupGuard:     guard.NewGroupGuard(),
		groupMapper:    newGroupMapper(),
		permissionRepo: repo.NewPermissionRepo(),
		fileIdentifier: infra.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		pdfPipeline:    conversion.NewPDFPipeline(),
		imagePipeline:  conversion.NewImagePipeline(),
		officePipeline: conversion.NewOfficePipeline(),
		videoPipeline:  conversion.NewVideoPipeline(),
		config:         config.GetConfig(),
	}
}

func (svc *FileService) Create(req FileCreateOptions, userId string) (*File, error) {
	if len(*req.ParentId) > 0 {
		if err := svc.validateParent(*req.ParentId, userId); err != nil {
			return nil, err
		}
	}
	file, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        req.Name,
		WorkspaceId: req.WorkspaceId,
		ParentId:    req.ParentId,
		Type:        req.Type,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(file.GetID(), userId, model.PermissionOwner); err != nil {
		return nil, err
	}
	file, err = svc.fileRepo.Find(file.GetID())
	if err != nil {
		return nil, err
	}
	if err = svc.fileSearch.Index([]model.File{file}); err != nil {
		return nil, err
	}
	if err = svc.fileCache.Set(file); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.MapFile(file, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) validateParent(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFolder {
		return errorpkg.NewFileIsNotAFolderError(file)
	}
	return nil
}

func (svc *FileService) Store(fileId string, filePath string, userId string) (*File, error) {
	file, err := svc.fileRepo.Find(fileId)
	if err != nil {
		return nil, err
	}
	if err = svc.fileCache.Set(file); err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	latestVersion, err := svc.snapshotRepo.GetLatestVersionForFile(fileId)
	if err != nil {
		return nil, err
	}
	snapshotId := helper.NewId()
	snapshot := repo.NewSnapshot()
	snapshot.SetID(snapshotId)
	snapshot.SetVersion(latestVersion)
	if err = svc.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if err = svc.snapshotRepo.MapWithFile(snapshotId, fileId); err != nil {
		return nil, err
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	original := model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(fileId + "/" + snapshotId + "/original" + strings.ToLower(filepath.Ext(filePath))),
		Size:   stat.Size(),
	}
	if err = svc.s3.PutFile(original.Key, filePath, infra.DetectMimeFromFile(filePath), workspace.GetBucket()); err != nil {
		return nil, err
	}
	snapshot.SetOriginal(&original)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if stat.Size() >= int64(svc.config.Limits.FileProcessingMaxSizeMB*1000000) {
		v, err := svc.fileMapper.MapFile(file, userId)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if svc.fileIdentifier.IsPDF(filepath.Ext(filePath)) {
		snapshot.SetPreview(&model.S3Object{
			Bucket: original.Bucket,
			Key:    original.Key,
			Size:   original.Size,
		})
		if err = svc.snapshotRepo.Save(snapshot); err != nil {
			return nil, err
		}
		if file, err = svc.fileRepo.Find(fileId); err != nil {
			return nil, err
		}
		if err := svc.fileCache.Set(file); err != nil {
			return nil, err
		}
		if err = svc.pdfPipeline.Run(conversion.PipelineOptions{
			FileID:     fileId,
			SnapshotID: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if svc.fileIdentifier.IsOffice(filepath.Ext(filePath)) || svc.fileIdentifier.IsPlainText(filepath.Ext(filePath)) {
		if err = svc.officePipeline.Run(conversion.PipelineOptions{
			FileID:     fileId,
			SnapshotID: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if svc.fileIdentifier.IsImage(filepath.Ext(filePath)) {
		if err = svc.imagePipeline.Run(conversion.PipelineOptions{
			FileID:     fileId,
			SnapshotID: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if svc.fileIdentifier.IsVideo(filepath.Ext(filePath)) {
		if err = svc.videoPipeline.Run(conversion.PipelineOptions{
			FileID:     fileId,
			SnapshotID: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	v, err := svc.fileMapper.MapFile(file, userId)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (svc *FileService) DownloadOriginalFile(id string, userId string) (string, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return "", nil, nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return "", nil, nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return "", nil, nil, err
	}
	snapshots := file.GetSnapshots()
	if len(snapshots) == 0 {
		return "", nil, nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	latestSnapshot := snapshots[len(snapshots)-1]
	if latestSnapshot.HasOriginal() {
		original := latestSnapshot.GetOriginal()
		path := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(original.Key))
		if err := svc.s3.GetFile(original.Key, path, original.Bucket); err != nil {
			return "", nil, nil, err
		}
		return path, file, latestSnapshot, nil
	} else {
		return "", nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadOriginalBuffer(id string, userId string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, nil, nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	snapshots := file.GetSnapshots()
	if len(snapshots) == 0 {
		return nil, nil, nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	latestSnapshot := snapshots[len(snapshots)-1]
	if latestSnapshot.HasOriginal() {
		original := latestSnapshot.GetOriginal()
		buf, err := svc.s3.GetObject(original.Key, original.Bucket)
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, latestSnapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadPreviewFile(id string, userId string) (string, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return "", nil, nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return "", nil, nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return "", nil, nil, err
	}
	snapshots := file.GetSnapshots()
	if len(snapshots) == 0 {
		return "", nil, nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	latestSnapshot := snapshots[len(snapshots)-1]
	if latestSnapshot.HasPreview() {
		preview := latestSnapshot.GetPreview()
		path := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(preview.Key))
		if err := svc.s3.GetFile(preview.Key, path, preview.Bucket); err != nil {
			return "", nil, nil, err
		}
		return path, file, latestSnapshot, nil
	} else {
		return "", nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadPreviewBuffer(id string, userId string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, nil, nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	snapshots := file.GetSnapshots()
	if len(snapshots) == 0 {
		return nil, nil, nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	latestSnapshot := snapshots[len(snapshots)-1]
	if latestSnapshot.HasPreview() {
		preview := latestSnapshot.GetPreview()
		buf, err := svc.s3.GetObject(preview.Key, preview.Bucket)
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, latestSnapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) FindByID(ids []string, userId string) ([]*File, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	var res []*File
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		f, err := svc.fileMapper.MapFile(file, userId)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (svc *FileService) FindByPath(path string, userId string) (*File, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		return &File{
			ID:          user.GetID(),
			WorkspaceId: "",
			Name:        "/",
			Type:        model.FileTypeFolder,
			Permission:  "owner",
			CreateTime:  user.GetCreateTime(),
			UpdateTime:  nil,
		}, nil
	}
	components := []string{}
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	workspace, err := svc.workspaceSvc.FindByName(components[0], userId)
	if err != nil {
		return nil, err
	}
	if len(components) == 1 {
		return &File{
			ID:          workspace.RootID,
			WorkspaceId: workspace.ID,
			Name:        workspace.Name,
			Type:        model.FileTypeFolder,
			Permission:  workspace.Permission,
			CreateTime:  workspace.CreateTime,
			UpdateTime:  workspace.UpdateTime,
		}, nil
	}
	currentID := workspace.RootID
	components = components[1:]
	for _, component := range components {
		ids, err := svc.fileRepo.GetChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.getAuthorized(ids, user)
		if err != nil {
			return nil, err
		}
		var filtered []model.File
		for _, f := range authorized {
			if f.GetName() == component {
				filtered = append(filtered, f)
			}
		}
		if len(filtered) > 0 {
			item := filtered[0]
			currentID = item.GetID()
			if item.GetType() == model.FileTypeFolder {
				continue
			} else if item.GetType() == model.FileTypeFile {
				break
			}
		} else {
			return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("component not found '%s'", component))
		}
	}
	result, err := svc.FindByID([]string{currentID}, userId)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (svc *FileService) ListByPath(path string, userId string) ([]*File, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		workspaces, err := svc.workspaceSvc.FindAll(userId)
		if err != nil {
			return nil, err
		}
		result := []*File{}
		for _, w := range workspaces {
			result = append(result, &File{
				ID:          w.RootID,
				WorkspaceId: w.ID,
				Name:        w.Name,
				Type:        model.FileTypeFolder,
				Permission:  w.Permission,
				CreateTime:  w.CreateTime,
				UpdateTime:  w.UpdateTime,
			})
		}
		return result, nil
	}
	components := []string{}
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	workspace, err := svc.workspaceRepo.FindByName(components[0])
	if err != nil {
		return nil, err
	}
	currentID := workspace.GetRootID()
	currentType := model.FileTypeFolder
	components = components[1:]
	for _, component := range components {
		ids, err := svc.fileRepo.GetChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.getAuthorized(ids, user)
		if err != nil {
			return nil, err
		}
		var filtered []model.File
		for _, f := range authorized {
			if f.GetName() == component {
				filtered = append(filtered, f)
			}
		}
		if len(filtered) > 0 {
			item := filtered[0]
			currentID = item.GetID()
			currentType = item.GetType()
			if item.GetType() == model.FileTypeFolder {
				continue
			} else if item.GetType() == model.FileTypeFile {
				break
			}
		} else {
			return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("component not found '%s'", component))
		}
	}
	if currentType == model.FileTypeFolder {
		ids, err := svc.fileRepo.GetChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.getAuthorized(ids, user)
		if err != nil {
			return nil, err
		}
		result, err := svc.fileMapper.MapFiles(authorized, userId)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else if currentType == model.FileTypeFile {
		result, err := svc.FindByID([]string{currentID}, userId)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("invalid file type %s", currentType))
	}
}

func (svc *FileService) ListByID(id string, page uint, size uint, sortBy string, sortOrder string, fileType string, userId string) (*FileList, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if page < 1 {
		return nil, errorpkg.NewInvalidPageParameterError()
	}
	if size < 1 {
		return nil, errorpkg.NewInvalidSizeParameterError()
	}
	ids, err := svc.fileRepo.GetChildrenIDs(id)
	if err != nil {
		return nil, err
	}
	authorizedFiles, err := svc.getAuthorized(ids, user)
	if err != nil {
		return nil, err
	}
	var filteredFiles []model.File
	for _, f := range authorizedFiles {
		if fileType == "" || f.GetType() == fileType {
			filteredFiles = append(filteredFiles, f)
		}
	}
	sortedFiles := svc.doSorting(filteredFiles, sortBy, sortOrder, userId)
	data, totalElements, totalPages := svc.doPaging(sortedFiles, page, size)
	mappedFiles, err := svc.fileMapper.MapFiles(data, userId)
	if err != nil {
		return nil, err
	}
	return &FileList{
		Data:          mappedFiles,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          page,
		Size:          size,
	}, nil
}

func (svc *FileService) doSorting(data []model.File, sortBy string, sortOrder string, userId string) []model.File {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == SortBySize {
		sort.Slice(data, func(i, j int) bool {
			fileA, err := svc.fileMapper.MapFile(data[i], userId)
			if err != nil {
				return false
			}
			fileB, err := svc.fileMapper.MapFile(data[j], userId)
			if err != nil {
				return false
			}
			var sizeA int64 = 0
			if fileA.Original != nil {
				sizeA = int64(fileA.Original.Size)
			}
			var sizeB int64 = 0
			if fileB.Original != nil {
				sizeB = int64(fileB.Original.Size)
			}
			if sortOrder == SortOrderDesc {
				return sizeA > sizeB
			} else {
				return sizeA < sizeB
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
	} else if sortBy == SortByKind {
		folders, _ := rxgo.Just(data)().
			Filter(func(v interface{}) bool {
				return v.(model.File).GetType() == model.FileTypeFolder
			}).
			ToSlice(0)
		files, _ := rxgo.Just(data)().
			Filter(func(v interface{}) bool {
				return v.(model.File).GetType() == model.FileTypeFile
			}).
			ToSlice(0)
		images, _ := rxgo.Just(files)().
			Filter(func(v interface{}) bool {
				f, err := svc.fileMapper.MapFile(v.(model.File), userId)
				if err != nil {
					return false
				}
				if f.Original == nil {
					return false
				}
				if svc.fileIdentifier.IsImage(f.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		pdfs, _ := rxgo.Just(files)().
			Filter(func(v interface{}) bool {
				f, err := svc.fileMapper.MapFile(v.(model.File), userId)
				if err != nil {
					return false
				}
				if f.Original == nil {
					return false
				}
				if svc.fileIdentifier.IsPDF(f.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		documents, _ := rxgo.Just(files)().
			Filter(func(v interface{}) bool {
				f, err := svc.fileMapper.MapFile(v.(model.File), userId)
				if err != nil {
					return false
				}
				if f.Original == nil {
					return false
				}
				if svc.fileIdentifier.IsOffice(f.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		videos, _ := rxgo.Just(files)().
			Filter(func(v interface{}) bool {
				f, err := svc.fileMapper.MapFile(v.(model.File), userId)
				if err != nil {
					return false
				}
				if f.Original == nil {
					return false
				}
				if svc.fileIdentifier.IsVideo(f.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		texts, _ := rxgo.Just(files)().
			Filter(func(v interface{}) bool {
				f, err := svc.fileMapper.MapFile(v.(model.File), userId)
				if err != nil {
					return false
				}
				if f.Original == nil {
					return false
				}
				if svc.fileIdentifier.IsPlainText(f.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		others, _ := rxgo.Just(files)().
			Filter(func(v interface{}) bool {
				f, err := svc.fileMapper.MapFile(v.(model.File), userId)
				if err != nil {
					return false
				}
				if f.Original == nil {
					return false
				}
				if !svc.fileIdentifier.IsImage(f.Original.Extension) &&
					!svc.fileIdentifier.IsPDF(f.Original.Extension) &&
					!svc.fileIdentifier.IsOffice(f.Original.Extension) &&
					!svc.fileIdentifier.IsVideo(f.Original.Extension) &&
					!svc.fileIdentifier.IsPlainText(f.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		var res []model.File
		for _, v := range folders {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range images {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range pdfs {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range documents {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range videos {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range texts {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		for _, v := range others {
			var file model.File
			file, err := svc.fileCache.Get(v.(model.File).GetID())
			if err != nil {
				return data
			}
			res = append(res, file)
		}
		return res
	}
	return data
}

func (svc *FileService) doPaging(files []model.File, page uint, size uint) ([]model.File, uint, uint) {
	page = page - 1
	low := size * page
	high := low + size
	var pagedFiles []model.File
	if low >= uint(len(files)) {
		pagedFiles = []model.File{}
	} else if high >= uint(len(files)) {
		high = uint(len(files))
		pagedFiles = files[low:high]
	} else {
		pagedFiles = files[low:high]
	}
	totalElements := uint(len(files))
	var totalPages uint
	if totalElements == 0 {
		totalPages = 1
	} else {
		if size > uint(len(files)) {
			size = uint(len(files))
		}
		totalPages = totalElements / size
		if totalPages == 0 {
			totalPages = 1
		}
		if totalElements%size > 0 {
			totalPages = totalPages + 1
		}
	}
	return pagedFiles, totalElements, totalPages
}

func (svc *FileService) getAuthorized(ids []string, user model.User) ([]model.File, error) {
	var res []model.File
	for _, id := range ids {
		var file model.File
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.fileGuard.IsAuthorized(user, file, model.PermissionViewer) {
			res = append(res, file)
		}
	}
	return res, nil
}

func (svc *FileService) Search(req FileSearchOptions, page uint, size uint, userId string) (*FileSearchResult, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceRepo.FindByID(req.WorkspaceId)
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceGuard.Authorize(user, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	files, err := svc.fileSearch.Query(req.Text)
	if err != nil {
		return nil, err
	}
	data, totalElements, totalPages, err := svc.doFilteringAndPaging(req, files, page, size, userId)
	if err != nil {
		return nil, err
	}
	v, err := svc.fileMapper.MapFiles(data, userId)
	if err != nil {
		return nil, err
	}
	res := &FileSearchResult{
		Data:          v,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          page,
		Size:          size,
		Query:         req,
	}
	return res, nil
}

func (svc *FileService) doFilteringAndPaging(req FileSearchOptions, files []model.File, page uint, size uint, userId string) ([]model.File, uint, uint, error) {
	filtered, _ := rxgo.Just(files)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetWorkspaceID() == req.WorkspaceId
		}).
		Filter(func(v interface{}) bool {
			if req.Type != nil {
				return v.(model.File).GetType() == *req.Type
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			file := v.(model.File)
			if req.ParentId != nil {
				res, err := svc.fileRepo.IsGrandChildOf(file.GetID(), *req.ParentId)
				if err != nil {
					return false
				}
				return res
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if req.CreateTimeBefore != nil {
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return t.UnixMilli() >= *req.CreateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if req.CreateTimeBefore != nil {
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return t.UnixMilli() <= *req.CreateTimeBefore
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if req.UpdateTimeAfter != nil {
				file := v.(model.File)
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return file.GetUpdateTime() != nil && t.UnixMilli() >= *req.UpdateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if req.UpdateTimeBefore != nil {
				file := v.(model.File)
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return file.GetUpdateTime() != nil && t.UnixMilli() <= *req.UpdateTimeBefore
			} else {
				return true
			}
		}).
		Skip((page - 1) * size).
		Take(size).
		ToSlice(0)
	var res []model.File
	for _, v := range filtered {
		var file model.File
		file, err := svc.fileCache.Get(v.(model.File).GetID())
		if err != nil {
			return nil, 0, 0, err
		}
		res = append(res, file)
	}
	totalPages := uint(len(res)) / size
	if totalPages == 0 {
		totalPages = 1
	}
	totalElements := uint(len(res))
	return res, totalElements, totalPages, nil
}

func (svc *FileService) GetPath(id string, userId string) ([]*File, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	path, err := svc.fileRepo.FindPath(id)
	if err != nil {
		return nil, err
	}
	res := []*File{}
	for _, f := range path {
		v, err := svc.fileMapper.MapFile(f, userId)
		if err != nil {
			return nil, err
		}
		res = append([]*File{v}, res...)
	}
	return res, nil
}

func (svc *FileService) Copy(targetId string, sourceIds []string, userId string) ([]*File, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	target, err := svc.fileCache.Get(targetId)
	if err != nil {
		return nil, err
	}

	/* Do checks */
	for _, sourceId := range sourceIds {
		var source model.File
		if source, err = svc.fileCache.Get(sourceId); err != nil {
			return nil, err
		}
		if err = svc.fileGuard.Authorize(user, target, model.PermissionEditor); err != nil {
			return nil, err
		}
		if err = svc.fileGuard.Authorize(user, source, model.PermissionEditor); err != nil {
			return nil, err
		}
		if source.GetID() == target.GetID() {
			return nil, errorpkg.NewFileCannotBeCopiedIntoIselfError(source)
		}
		if target.GetType() != model.FileTypeFolder {
			return nil, errorpkg.NewFileIsNotAFolderError(target)
		}
		if yes, _ := svc.fileRepo.IsGrandChildOf(target.GetID(), source.GetID()); yes {
			return nil, errorpkg.NewFileCannotBeCopiedIntoOwnSubtreeError(source)
		}
	}

	/* Do copying */
	allClones := []model.File{}
	for _, sourceId := range sourceIds {
		/* Get original tree */
		var sourceTree []model.File
		if sourceTree, err = svc.fileRepo.FindTree(sourceId); err != nil {
			return nil, err
		}

		/* Clone source tree */
		var rootCloneIndex int
		var cloneIds = make(map[string]string)
		var originalIds = make(map[string]string)
		var clones []model.File
		var permissions []*repo.UserPermission
		for i, o := range sourceTree {
			c := repo.NewFile()
			c.SetID(helper.NewId())
			c.SetParentID(o.GetParentID())
			c.SetWorkspaceID(o.GetWorkspaceID())
			c.SetType(o.GetType())
			c.SetName(o.GetName())
			c.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
			if o.GetID() == sourceId {
				rootCloneIndex = i
			}
			cloneIds[o.GetID()] = c.GetID()
			originalIds[c.GetID()] = o.GetID()
			clones = append(clones, c)
			permissions = append(permissions, &repo.UserPermission{
				ID:         helper.NewId(),
				UserID:     userId,
				ResourceID: c.GetID(),
				Permission: model.PermissionOwner,
				CreateTime: time.Now().UTC().Format(time.RFC3339),
			})
		}

		/* Set parent Ids of clones */
		for i, c := range clones {
			id := cloneIds[*c.GetParentID()]
			clones[i].SetParentID(&id)
		}

		/* Parent ID of root clone is target ID */
		if clones != nil {
			clones[rootCloneIndex].SetParentID(&targetId)
		}

		/* Persist clones */
		if err = svc.fileRepo.BulkInsert(clones, 100); err != nil {
			return nil, err
		}

		/* Persist permissions */
		if err = svc.fileRepo.BulkInsertPermissions(permissions, 100); err != nil {
			return nil, err
		}

		/* Assign snapshots to clones */
		for _, c := range clones {
			if err := svc.fileRepo.AssignSnapshots(c.GetID(), originalIds[c.GetID()]); err != nil {
				return nil, err
			}
		}

		/* Index clones for search */
		if err := svc.fileSearch.Index(clones); err != nil {
			return nil, err
		}

		/* Create cache for clones */
		for _, c := range clones {
			if _, err := svc.fileCache.Refresh(c.GetID()); err != nil {
				return nil, err
			}
		}

		allClones = append(allClones, clones...)
	}

	/* Refresh updateTime on target */
	timeNow := time.Now().UTC().Format(time.RFC3339)
	target.SetUpdateTime(&timeNow)
	if err := svc.fileRepo.Save(target); err != nil {
		return nil, err
	}

	res, err := svc.fileMapper.MapFiles(allClones, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *FileService) Move(targetId string, sourceIds []string, userId string) ([]string, error) {
	res := []string{}
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return []string{}, err
	}
	target, err := svc.fileCache.Get(targetId)
	if err != nil {
		return []string{}, err
	}

	/* Do checks */
	for _, id := range sourceIds {
		source, err := svc.fileCache.Get(id)
		if err != nil {
			return []string{}, err
		}
		if err := svc.fileGuard.Authorize(user, target, model.PermissionEditor); err != nil {
			return []string{}, err
		}
		if err := svc.fileGuard.Authorize(user, source, model.PermissionEditor); err != nil {
			return []string{}, err
		}
		if source.GetParentID() != nil && *source.GetParentID() == target.GetID() {
			return []string{}, errorpkg.NewFileAlreadyChildOfDestinationError(source, target)
		}
		if target.GetID() == source.GetID() {
			return []string{}, errorpkg.NewFileCannotBeMovedIntoItselfError(source)
		}
		if target.GetType() != model.FileTypeFolder {
			return []string{}, errorpkg.NewFileIsNotAFolderError(target)
		}
		targetIsGrandChildOfSource, _ := svc.fileRepo.IsGrandChildOf(target.GetID(), source.GetID())
		if targetIsGrandChildOfSource {
			return []string{}, errorpkg.NewTargetIsGrandChildOfSourceError(source)
		}
	}

	/* Do moving */
	for _, id := range sourceIds {
		source, _ := svc.fileCache.Get(id)

		/* Add old parent */
		res = append(res, *source.GetParentID())

		/* Move source into target */
		if err := svc.fileRepo.MoveSourceIntoTarget(target.GetID(), source.GetID()); err != nil {
			return []string{}, err
		}

		/* Get updated source */
		source, err = svc.fileRepo.Find(source.GetID())
		if err != nil {
			return []string{}, err
		}

		// Add new parent
		res = append(res, *source.GetParentID())

		/* Refresh updateTime on source and target */
		timeNow := time.Now().UTC().Format(time.RFC3339)
		source.SetUpdateTime(&timeNow)
		if err := svc.fileRepo.Save(source); err != nil {
			return []string{}, err
		}
		target.SetUpdateTime(&timeNow)
		if err := svc.fileRepo.Save(target); err != nil {
			return []string{}, err
		}
		if err := svc.fileSearch.Update([]model.File{source}); err != nil {
			return []string{}, err
		}
		sourceTree, err := svc.fileRepo.FindTree(source.GetID())
		if err != nil {
			return []string{}, err
		}
		for _, f := range sourceTree {
			if err := svc.fileCache.Set(f); err != nil {
				return []string{}, err
			}
		}
	}
	return res, nil
}

func (svc *FileService) Rename(id string, name string, userId string) (*File, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	file.SetName(name)
	if err = svc.fileRepo.Save(file); err != nil {
		return nil, err
	}
	if err = svc.fileSearch.Update([]model.File{file}); err != nil {
		return nil, err
	}
	err = svc.fileCache.Set(file)
	if err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.MapFile(file, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) Delete(ids []string, userId string) ([]string, error) {
	var res []string
	for _, id := range ids {
		var user model.User
		user, err := svc.userRepo.Find(userId)
		if err != nil {
			return nil, err
		}
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		if file.GetParentID() == nil {
			workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
			if err != nil {
				return []string{}, err
			}
			return nil, errorpkg.NewCannotDeleteWorkspaceRootError(file, workspace)
		}
		if err = svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
			return nil, err
		}

		// Add parent
		res = append(res, *file.GetParentID())

		var tree []model.File
		tree, err = svc.fileRepo.FindTree(file.GetID())
		if err != nil {
			return nil, err
		}
		var treeIds []string
		for _, f := range tree {
			treeIds = append(treeIds, f.GetID())
		}
		if err := svc.fileSearch.Delete(treeIds); err != nil {
			// Here we don't return an error or panic but we just print the error
			fmt.Println(err)
		}
		for _, f := range tree {
			if err = svc.fileCache.Delete(f.GetID()); err != nil {
				return nil, err
			}
		}
		for _, f := range tree {
			if err = svc.fileRepo.Delete(f.GetID()); err != nil {
				return nil, err
			}
			if err = svc.snapshotRepo.DeleteMappingsForFile(f.GetID()); err != nil {
				return nil, err
			}
		}
		var danglingSnapshots []model.Snapshot
		danglingSnapshots, err = svc.snapshotRepo.FindAllDangling()
		if err != nil {
			return nil, err
		}
		for _, s := range danglingSnapshots {
			if s.HasOriginal() {
				if err = svc.s3.RemoveObject(s.GetOriginal().Key, s.GetOriginal().Bucket); err != nil {
					return nil, err
				}
			}
			if s.HasPreview() {
				if err = svc.s3.RemoveObject(s.GetPreview().Key, s.GetPreview().Bucket); err != nil {
					return nil, err
				}
			}
		}
		if err = svc.snapshotRepo.DeleteAllDangling(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *FileService) GetSize(id string, userId string) (int64, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return -1, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return -1, err
	}
	if err := svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return -1, err
	}
	res, err := svc.fileRepo.GetSize(id)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (svc *FileService) GetItemCount(id string, userId string) (int64, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return 0, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return 0, err
	}
	if err := svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return 0, err
	}
	res, err := svc.fileRepo.GetItemCount(id)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (svc *FileService) GrantUserPermission(ids []string, assigneeId string, permission string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err = svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
			return err
		}
		if _, err := svc.userRepo.Find(assigneeId); err != nil {
			return err
		}
		if err = svc.fileRepo.GrantUserPermission(id, assigneeId, permission); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		workspace, err := svc.workspaceRepo.FindByID(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		if err = svc.workspaceCache.Set(workspace); err != nil {
			return err
		}
		path, err := svc.fileRepo.FindPath(id)
		if err != nil {
			return err
		}
		for _, f := range path {
			if err := svc.fileCache.Set(f); err != nil {
				return err
			}
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		for _, f := range tree {
			if err := svc.fileCache.Set(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FileService) RevokeUserPermission(ids []string, assigneeId string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err := svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
			return err
		}
		if _, err := svc.userRepo.Find(assigneeId); err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeUserPermission(id, assigneeId); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *FileService) GrantGroupPermission(ids []string, groupId string, permission string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err = svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
			return err
		}
		group, err := svc.groupCache.Get(groupId)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
			return err
		}
		if err = svc.fileRepo.GrantGroupPermission(id, groupId, permission); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		workspace, err := svc.workspaceRepo.FindByID(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		if err = svc.workspaceCache.Set(workspace); err != nil {
			return err
		}
		path, err := svc.fileRepo.FindPath(id)
		if err != nil {
			return err
		}
		for _, f := range path {
			if err := svc.fileCache.Set(f); err != nil {
				return err
			}
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		for _, f := range tree {
			if err := svc.fileCache.Set(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FileService) RevokeGroupPermission(ids []string, groupId string, userId string) error {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err := svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
			return err
		}
		group, err := svc.groupCache.Get(groupId)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeGroupPermission(id, groupId); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *FileService) GetUserPermissions(id string, userId string) ([]*UserPermission, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.GetUserPermissions(id)
	if err != nil {
		return nil, err
	}
	res := make([]*UserPermission, 0)
	for _, p := range permissions {
		if p.UserID == userId {
			continue
		}
		u, err := svc.userRepo.Find(p.UserID)
		if err != nil {
			return nil, err
		}
		res = append(res, &UserPermission{
			ID:         p.ID,
			User:       svc.userMapper.mapUser(u),
			Permission: p.Permission,
		})
	}
	return res, nil
}

func (svc *FileService) GetGroupPermissions(id string, userId string) ([]*GroupPermission, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(user, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.GetGroupPermissions(id)
	if err != nil {
		return nil, err
	}
	res := make([]*GroupPermission, 0)
	for _, p := range permissions {
		m, err := svc.groupCache.Get(p.GroupID)
		if err != nil {
			return nil, err
		}
		g, err := svc.groupMapper.mapGroup(m, userId)
		if err != nil {
			return nil, err
		}
		res = append(res, &GroupPermission{
			ID:         p.ID,
			Group:      g,
			Permission: p.Permission,
		})
	}
	return res, nil
}

type FileMapper struct {
	groupCache *cache.GroupCache
	config     config.Config
}

func NewFileMapper() *FileMapper {
	return &FileMapper{
		groupCache: cache.NewGroupCache(),
		config:     config.GetConfig(),
	}
}

func (mp *FileMapper) MapFile(m model.File, userId string) (*File, error) {
	snapshots := m.GetSnapshots()
	res := &File{
		ID:          m.GetID(),
		WorkspaceId: m.GetWorkspaceID(),
		Name:        m.GetName(),
		Type:        m.GetType(),
		ParentId:    m.GetParentID(),
		Snapshots:   mp.MapSnapshots(snapshots, m.GetID()),
		CreateTime:  m.GetCreateTime(),
		UpdateTime:  m.GetUpdateTime(),
	}
	if len(snapshots) > 0 {
		latest := mp.MapSnapshot(snapshots[len(snapshots)-1])
		res.Version = &latest.Version
		res.Original = latest.Original
		res.Preview = latest.Preview
		res.Thumbnail = latest.Thumbnail
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
	shareCount := 0
	for _, p := range m.GetUserPermissions() {
		if p.GetUserID() != userId {
			shareCount++
		}
	}
	shareCount += len(m.GetGroupPermissions())
	if shareCount > 0 {
		res.IsShared = true
	} else {
		res.IsShared = false
	}
	return res, nil
}

func (mp *FileMapper) MapFiles(files []model.File, userId string) ([]*File, error) {
	res := make([]*File, 0)
	for _, f := range files {
		v, err := mp.MapFile(f, userId)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func (mp *FileMapper) MapSnapshot(m model.Snapshot) *Snapshot {
	s := &Snapshot{
		ID:      m.GetID(),
		Version: m.GetVersion(),
	}
	if m.HasOriginal() {
		s.Original = mp.MapOriginal(m.GetOriginal())
	}
	if m.HasPreview() {
		s.Preview = mp.MapPreview(m.GetPreview())
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.MapThumbnail(m.GetThumbnail())
	}
	return s
}

func (mp *FileMapper) MapOriginal(m *model.S3Object) *Download {
	download := &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
	if m.Image != nil {
		download.Image = &ImageProps{
			Width:  m.Image.Width,
			Height: m.Image.Height,
		}
	}
	return download
}

func (mp *FileMapper) MapPreview(m *model.S3Object) *Download {
	download := &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
	if m.Image != nil {
		download.Image = &ImageProps{
			Width:  m.Image.Width,
			Height: m.Image.Height,
		}
	}
	return download
}

func (mp *FileMapper) MapOcr(m *model.S3Object) *Download {
	return &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
}

func (mp *FileMapper) MapThumbnail(m *model.Thumbnail) *Thumbnail {
	return &Thumbnail{
		Base64: m.Base64,
		Width:  m.Width,
		Height: m.Height,
	}
}

func (mp *FileMapper) MapText(m *model.S3Object) *Download {
	return &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
}

func (mp *FileMapper) MapSnapshots(snapshots []model.Snapshot, fileId string) []*Snapshot {
	res := make([]*Snapshot, 0)
	for _, s := range snapshots {
		res = append(res, mp.MapSnapshot(s))
	}
	return res
}
