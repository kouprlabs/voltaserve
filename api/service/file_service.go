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
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/reactivex/rxgo/v2"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

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
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	userRepo       repo.UserRepo
	userMapper     *userMapper
	groupCache     *cache.GroupCache
	groupGuard     *guard.GroupGuard
	groupMapper    *groupMapper
	permissionRepo repo.PermissionRepo
	taskCache      *cache.TaskCache
	taskSvc        *TaskService
	fileIdent      *infra.FileIdentifier
	s3             *infra.S3Manager
	pipelineClient *conversion_client.PipelineClient
	config         *config.Config
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
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotSvc:    NewSnapshotService(),
		userRepo:       repo.NewUserRepo(),
		userMapper:     newUserMapper(),
		groupCache:     cache.NewGroupCache(),
		groupGuard:     guard.NewGroupGuard(),
		groupMapper:    newGroupMapper(),
		permissionRepo: repo.NewPermissionRepo(),
		taskCache:      cache.NewTaskCache(),
		taskSvc:        NewTaskService(),
		fileIdent:      infra.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		pipelineClient: conversion_client.NewPipelineClient(),
		config:         config.GetConfig(),
	}
}

type FileCreateOptions struct {
	WorkspaceID string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name"        validate:"required,max=255"`
	Type        string  `json:"type"        validate:"required,oneof=file folder"`
	ParentID    *string `json:"parentId"    validate:"required"`
}

type File struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ParentID    *string   `json:"parentId,omitempty"`
	Permission  string    `json:"permission"`
	IsShared    *bool     `json:"isShared,omitempty"`
	Snapshot    *Snapshot `json:"snapshot,omitempty"`
	CreateTime  string    `json:"createTime"`
	UpdateTime  *string   `json:"updateTime,omitempty"`
}

func (svc *FileService) Create(opts FileCreateOptions, userID string) (*File, error) {
	var components []string
	for _, component := range strings.Split(opts.Name, "/") {
		if component != "" {
			components = append(components, component)
		}
	}
	parentID := opts.ParentID
	if len(components) > 1 {
		for _, component := range components[:len(components)-1] {
			existing, err := svc.getChildWithName(*parentID, component)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				parentID = new(string)
				*parentID = existing.GetID()
			} else {
				res, err := svc.create(FileCreateOptions{
					Name:        component,
					Type:        model.FileTypeFolder,
					ParentID:    parentID,
					WorkspaceID: opts.WorkspaceID,
				}, userID)
				if err != nil {
					return nil, err
				}
				parentID = &res.ID
			}
		}
	}
	name := components[len(components)-1]
	return svc.create(FileCreateOptions{
		WorkspaceID: opts.WorkspaceID,
		Name:        name,
		Type:        opts.Type,
		ParentID:    parentID,
	}, userID)
}

func (svc *FileService) create(opts FileCreateOptions, userID string) (*File, error) {
	if len(*opts.ParentID) > 0 {
		if err := svc.validateParent(*opts.ParentID, userID); err != nil {
			return nil, err
		}
		existing, err := svc.getChildWithName(*opts.ParentID, opts.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errorpkg.NewFileWithSimilarNameExistsError()
		}
	}
	file, err := svc.fileRepo.Insert(repo.FileInsertOptions{
		Name:        opts.Name,
		WorkspaceID: opts.WorkspaceID,
		ParentID:    opts.ParentID,
		Type:        opts.Type,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.fileRepo.GrantUserPermission(file.GetID(), userID, model.PermissionOwner); err != nil {
		return nil, err
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	if err = svc.fileSearch.Index([]model.File{file}); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) validateParent(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFolder {
		return errorpkg.NewFileIsNotAFolderError(file)
	}
	return nil
}

type StoreOptions struct {
	S3Reference *model.S3Reference
	Path        *string
}

func (svc *FileService) Store(id string, opts StoreOptions, userID string) (*File, error) {
	file, err := svc.fileRepo.Find(id)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	latestVersion, err := svc.snapshotRepo.FindLatestVersionForFile(id)
	if err != nil {
		return nil, err
	}
	var snapshotID string
	var size int64
	var path string
	var original model.S3Object
	var bucket string
	var contentType string
	if opts.S3Reference == nil {
		snapshotID = helper.NewID()
		path = *opts.Path
		stat, err := os.Stat(*opts.Path)
		if err != nil {
			return nil, err
		}
		size = stat.Size()
		original = model.S3Object{
			Bucket: workspace.GetBucket(),
			Key:    snapshotID + "/original" + strings.ToLower(filepath.Ext(path)),
			Size:   helper.ToPtr(size),
		}
		bucket = workspace.GetBucket()
		contentType = infra.DetectMIMEFromPath(path)
	} else {
		snapshotID = opts.S3Reference.SnapshotID
		path = opts.S3Reference.Key
		size = opts.S3Reference.Size
		original = model.S3Object{
			Bucket: opts.S3Reference.Bucket,
			Key:    opts.S3Reference.Key,
			Size:   helper.ToPtr(size),
		}
		bucket = opts.S3Reference.Bucket
		contentType = opts.S3Reference.ContentType
	}
	snapshot := repo.NewSnapshot()
	snapshot.SetID(snapshotID)
	snapshot.SetVersion(latestVersion + 1)
	if err = svc.snapshotRepo.Insert(snapshot); err != nil {
		return nil, err
	}
	snapshot, err = svc.snapshotCache.Get(snapshotID)
	if err != nil {
		return nil, err
	}
	if err = svc.snapshotRepo.MapWithFile(snapshotID, id); err != nil {
		return nil, err
	}
	exceedsProcessingLimit := size > helper.MegabyteToByte(svc.fileIdent.GetProcessingLimitMB(path))
	if opts.S3Reference == nil {
		if err = svc.s3.PutFile(original.Key, path, contentType, bucket, minio.PutObjectOptions{}); err != nil {
			return nil, err
		}
	}
	snapshot.SetOriginal(&original)
	if exceedsProcessingLimit {
		snapshot.SetStatus(model.SnapshotStatusReady)
	} else {
		snapshot.SetStatus(model.SnapshotStatusWaiting)
	}
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return nil, err
	}
	file.SetSnapshotID(&snapshotID)
	if err := svc.fileRepo.Save(file); err != nil {
		return nil, err
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	if !exceedsProcessingLimit {
		task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
			ID:              helper.NewID(),
			Name:            "Waiting.",
			UserID:          userID,
			IsIndeterminate: true,
			Status:          model.TaskStatusWaiting,
			Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
		})
		if err != nil {
			return nil, err
		}
		snapshot.SetTaskID(helper.ToPtr(task.GetID()))
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			return nil, err
		}
		if err := svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
			TaskID:     task.GetID(),
			SnapshotID: snapshot.GetID(),
			Bucket:     original.Bucket,
			Key:        original.Key,
		}); err != nil {
			return nil, err
		}
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) DownloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (model.File, model.Snapshot, *infra.RangeInterval, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, nil, err
	}
	if snapshot.HasOriginal() {
		objectInfo, err := svc.s3.StatObject(snapshot.GetOriginal().Key, snapshot.GetOriginal().Bucket, minio.StatObjectOptions{})
		if err != nil {
			return nil, nil, nil, err
		}
		opts := minio.GetObjectOptions{}
		var ri *infra.RangeInterval
		if rangeHeader != "" {
			ri = infra.NewRangeInterval(rangeHeader, objectInfo.Size)
			ri.ApplyToMinIOGetObjectOptions(&opts)
		}
		if _, err := svc.s3.GetObjectWithBuffer(snapshot.GetOriginal().Key, snapshot.GetOriginal().Bucket, buf, opts); err != nil {
			return nil, nil, nil, err
		}
		return file, snapshot, ri, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (model.File, model.Snapshot, *infra.RangeInterval, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, nil, err
	}
	if snapshot.HasPreview() {
		objectInfo, err := svc.s3.StatObject(snapshot.GetOriginal().Key, snapshot.GetOriginal().Bucket, minio.StatObjectOptions{})
		if err != nil {
			return nil, nil, nil, err
		}
		opts := minio.GetObjectOptions{}
		var ri *infra.RangeInterval
		if rangeHeader != "" {
			ri = infra.NewRangeInterval(rangeHeader, objectInfo.Size)
			ri.ApplyToMinIOGetObjectOptions(&opts)
		}
		if _, err := svc.s3.GetObjectWithBuffer(snapshot.GetPreview().Key, snapshot.GetPreview().Bucket, buf, opts); err != nil {
			return nil, nil, nil, err
		}
		return file, snapshot, ri, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadThumbnailBuffer(id string, buf *bytes.Buffer, userID string) (model.Snapshot, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasThumbnail() {
		if _, err := svc.s3.GetObjectWithBuffer(snapshot.GetThumbnail().Key, snapshot.GetThumbnail().Bucket, buf, minio.GetObjectOptions{}); err != nil {
			return nil, err
		}
		return snapshot, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) Find(ids []string, userID string) ([]*File, error) {
	var res []*File
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			continue
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		f, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func (svc *FileService) FindByPath(path string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		return &File{
			ID:          user.GetID(),
			WorkspaceID: "",
			Name:        "/",
			Type:        model.FileTypeFolder,
			Permission:  "owner",
			CreateTime:  user.GetCreateTime(),
			UpdateTime:  nil,
		}, nil
	}
	components := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	workspace, err := svc.workspaceSvc.Find(helper.WorkspaceIDFromSlug(components[0]), userID)
	if err != nil {
		return nil, err
	}
	if len(components) == 1 {
		return &File{
			ID:          workspace.RootID,
			WorkspaceID: workspace.ID,
			Name:        helper.SlugFromWorkspace(workspace.ID, workspace.Name),
			Type:        model.FileTypeFolder,
			Permission:  workspace.Permission,
			CreateTime:  workspace.CreateTime,
			UpdateTime:  workspace.UpdateTime,
		}, nil
	}
	currentID := workspace.RootID
	components = components[1:]
	for _, component := range components {
		ids, err := svc.fileRepo.FindChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.doAuthorizationByIDs(ids, userID)
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
	result, err := svc.Find([]string{currentID}, userID)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("item not found '%s'", currentID))
	}
	return result[0], nil
}

func (svc *FileService) ListByPath(path string, userID string) ([]*File, error) {
	if path == "/" {
		workspaces, err := svc.workspaceSvc.findAllWithoutOptions(userID)
		if err != nil {
			return nil, err
		}
		result := make([]*File, 0)
		for _, w := range workspaces {
			result = append(result, &File{
				ID:          w.RootID,
				WorkspaceID: w.ID,
				Name:        helper.SlugFromWorkspace(w.ID, w.Name),
				Type:        model.FileTypeFolder,
				Permission:  w.Permission,
				CreateTime:  w.CreateTime,
				UpdateTime:  w.UpdateTime,
			})
		}
		return result, nil
	}
	components := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	workspace, err := svc.workspaceRepo.Find(helper.WorkspaceIDFromSlug(components[0]))
	if err != nil {
		return nil, err
	}
	currentID := workspace.GetRootID()
	currentType := model.FileTypeFolder
	components = components[1:]
	for _, component := range components {
		ids, err := svc.fileRepo.FindChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.doAuthorizationByIDs(ids, userID)
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
		ids, err := svc.fileRepo.FindChildrenIDs(currentID)
		if err != nil {
			return nil, err
		}
		authorized, err := svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
		result, err := svc.fileMapper.mapMany(authorized, userID)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else if currentType == model.FileTypeFile {
		result, err := svc.Find([]string{currentID}, userID)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("invalid file type %s", currentType))
	}
}

type FileQuery struct {
	Text             string  `json:"text"                       validate:"required"`
	Type             *string `json:"type,omitempty"             validate:"omitempty,oneof=file folder"`
	CreateTimeAfter  *int64  `json:"createTimeAfter,omitempty"`
	CreateTimeBefore *int64  `json:"createTimeBefore,omitempty"`
	UpdateTimeAfter  *int64  `json:"updateTimeAfter,omitempty"`
	UpdateTimeBefore *int64  `json:"updateTimeBefore,omitempty"`
}

type FileList struct {
	Data          []*File    `json:"data"`
	TotalPages    uint64     `json:"totalPages"`
	TotalElements uint64     `json:"totalElements"`
	Page          uint64     `json:"page"`
	Size          uint64     `json:"size"`
	Query         *FileQuery `json:"query,omitempty"`
}

type FileListOptions struct {
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
	Query     *FileQuery
}

func (svc *FileService) List(id string, opts FileListOptions, userID string) (*FileList, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	ids, err := svc.fileRepo.FindChildrenIDs(id)
	if err != nil {
		return nil, err
	}
	var data []model.File
	for _, id := range ids {
		var f model.File
		f, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		data = append(data, f)
	}
	var filtered []model.File
	for _, f := range data {
		if opts.Query == nil || *opts.Query.Type == "" || f.GetType() == *opts.Query.Type {
			filtered = append(filtered, f)
		}
	}
	authorized, err := svc.doAuthorization(filtered, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(authorized, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.fileMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &FileList{
		Data:          mapped,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          opts.Page,
		Size:          opts.Size,
	}, nil
}

type FileProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

func (svc *FileService) Probe(id string, opts FileListOptions, userID string) (*FileProbe, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	totalElements, err := svc.fileRepo.CountChildren(id)
	if err != nil {
		return nil, err
	}
	return &FileProbe{
		TotalElements: uint64(totalElements),
		TotalPages:    (uint64(totalElements) + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *FileService) Search(id string, opts FileListOptions, userID string) (*FileList, error) {
	parent, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceRepo.Find(parent.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	var data []model.File
	if opts.Query.Text != "" {
		count, err := svc.fileRepo.Count()
		if err != nil {
			return nil, err
		}
		data, err = svc.fileSearch.Query(opts.Query.Text, infra.QueryOptions{Limit: count})
		if err != nil {
			return nil, err
		}
	} else {
		ids, err := svc.fileRepo.FindChildrenIDs(id)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			var f model.File
			f, err := svc.fileCache.Get(id)
			if err != nil {
				return nil, err
			}
			data = append(data, f)
		}
	}
	filteredData, err := svc.doQueryFiltering(data, *opts.Query, parent)
	if err != nil {
		return nil, err
	}
	authorizedData, err := svc.doAuthorization(filteredData, userID)
	if err != nil {
		return nil, err
	}
	sortedData := svc.doSorting(authorizedData, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.doPagination(sortedData, opts.Page, opts.Size)
	mappedData, err := svc.fileMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	res := &FileList{
		Data:          mappedData,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          opts.Page,
		Size:          opts.Size,
		Query:         opts.Query,
	}
	return res, nil
}

func (svc *FileService) FindPath(id string, userID string) ([]*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	path, err := svc.fileRepo.FindPath(id)
	if err != nil {
		return nil, err
	}
	res := make([]*File, 0)
	for _, file := range path {
		f, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append([]*File{f}, res...)
	}
	return res, nil
}

func (svc *FileService) CopyOne(sourceID string, targetID string, userID string) (*File, error) {
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return nil, err
	}
	source, err := svc.fileCache.Get(sourceID)
	if err != nil {
		return nil, err
	}

	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Copying.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: source.GetName()},
	})
	if err != nil {
		return nil, err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())

	/* Do checks */
	if err := svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
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

	/* Read original tree */
	var sourceIds []string
	sourceIds, err = svc.fileRepo.FindTreeIDs(source.GetID())
	if err != nil {
		return nil, err
	}
	var sourceTree []model.File
	for _, id := range sourceIds {
		sourceFile, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		sourceTree = append(sourceTree, sourceFile)
	}

	/* Clone source tree */
	var rootCloneIndex int
	cloneIDs := make(map[string]string)
	originalIDs := make(map[string]string)
	var clones []model.File
	var permissions []model.UserPermission
	for i, sourceFile := range sourceTree {
		f := repo.NewFile()
		f.SetID(helper.NewID())
		f.SetParentID(sourceFile.GetParentID())
		f.SetWorkspaceID(sourceFile.GetWorkspaceID())
		f.SetSnapshotID(sourceFile.GetSnapshotID())
		f.SetType(sourceFile.GetType())
		f.SetName(sourceFile.GetName())
		f.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
		if sourceFile.GetID() == source.GetID() {
			rootCloneIndex = i
		}
		cloneIDs[sourceFile.GetID()] = f.GetID()
		originalIDs[f.GetID()] = sourceFile.GetID()
		clones = append(clones, f)

		p := repo.NewUserPermission()
		p.SetID(helper.NewID())
		p.SetUserID(userID)
		p.SetResourceID(f.GetID())
		p.SetPermission(model.PermissionOwner)
		p.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
		permissions = append(permissions, p)
	}

	/* Set parent IDs of clones */
	for i, f := range clones {
		id := cloneIDs[*f.GetParentID()]
		clones[i].SetParentID(&id)
	}

	rootClone := clones[rootCloneIndex]

	/* Parent ID of root clone is target ID */
	if clones != nil {
		rootClone.SetParentID(helper.ToPtr(target.GetID()))
	}

	/* If there is a file with similar name, append a prefix */
	existing, err := svc.getChildWithName(target.GetID(), rootClone.GetName())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		rootClone.SetName(helper.UniqueFilename(rootClone.GetName()))
	}

	const BulkInsertChunkSize = 1000

	/* Persist clones */
	if err = svc.fileRepo.BulkInsert(clones, BulkInsertChunkSize); err != nil {
		return nil, err
	}

	/* Persist permissions */
	if err = svc.fileRepo.BulkInsertPermissions(permissions, BulkInsertChunkSize); err != nil {
		return nil, err
	}

	/* Attach latest snapshot to clones */
	var mappings []*repo.SnapshotFileEntity
	for i, f := range clones {
		original := sourceTree[i]
		if original.GetSnapshotID() != nil {
			mappings = append(mappings, &repo.SnapshotFileEntity{
				SnapshotID: *original.GetSnapshotID(),
				FileID:     f.GetID(),
			})
		}
	}
	if err := svc.snapshotRepo.BulkMapWithFile(mappings, BulkInsertChunkSize); err != nil {
		return nil, err
	}

	/* Create cache for clones */
	for _, clone := range clones {
		if _, err := svc.fileCache.RefreshWithExisting(clone, userID); err != nil {
			log.GetLogger().Error(err)
		}
	}

	/* Index clones for search */
	go func() {
		if err := svc.fileSearch.Index(clones); err != nil {
			log.GetLogger().Error(err)
		}
	}()

	/* Refresh updateTime on target */
	timeNow := helper.NewTimestamp()
	target.SetUpdateTime(&timeNow)
	if err := svc.fileRepo.Save(target); err != nil {
		return nil, err
	}

	res, err := svc.fileMapper.mapOne(rootClone, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type FileCopyManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileCopyManyResult struct {
	New       []string `json:"new"`
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *FileService) CopyMany(opts FileCopyManyOptions, userID string) (*FileCopyManyResult, error) {
	res := &FileCopyManyResult{
		New:       make([]string, 0),
		Succeeded: make([]string, 0),
		Failed:    make([]string, 0),
	}
	for _, id := range opts.SourceIDs {
		file, err := svc.CopyOne(id, opts.TargetID, userID)
		if err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.New = append(res.New, file.ID)
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

func (svc *FileService) MoveOne(sourceID string, targetID string, userID string) (*File, error) {
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return nil, err
	}
	source, err := svc.fileCache.Get(sourceID)
	if err != nil {
		return nil, err
	}

	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Moving.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: source.GetName()},
	})
	if err != nil {
		return nil, err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())

	/* Do checks */
	if source.GetParentID() != nil {
		existing, err := svc.getChildWithName(target.GetID(), source.GetName())
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errorpkg.NewFileWithSimilarNameExistsError()
		}
	}
	if err := svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
		return nil, err
	}
	if source.GetParentID() != nil && *source.GetParentID() == target.GetID() {
		return nil, errorpkg.NewFileAlreadyChildOfDestinationError(source, target)
	}
	if target.GetID() == source.GetID() {
		return nil, errorpkg.NewFileCannotBeMovedIntoItselfError(source)
	}
	if target.GetType() != model.FileTypeFolder {
		return nil, errorpkg.NewFileIsNotAFolderError(target)
	}
	targetIsGrandChildOfSource, _ := svc.fileRepo.IsGrandChildOf(target.GetID(), source.GetID())
	if targetIsGrandChildOfSource {
		return nil, errorpkg.NewTargetIsGrandChildOfSourceError(source)
	}

	/* Move source into target */
	if err := svc.fileRepo.MoveSourceIntoTarget(target.GetID(), source.GetID()); err != nil {
		return nil, err
	}

	/* Read updated source */
	source, err = svc.fileRepo.Find(source.GetID())
	if err != nil {
		return nil, err
	}

	/* Refresh updateTime on source and target */
	timeNow := time.Now().UTC().Format(time.RFC3339)
	source.SetUpdateTime(&timeNow)
	if err := svc.fileRepo.Save(source); err != nil {
		return nil, err
	}
	if err := svc.sync(source); err != nil {
		return nil, err
	}
	target.SetUpdateTime(&timeNow)
	if err := svc.fileRepo.Save(target); err != nil {
		return nil, err
	}
	if err := svc.sync(target); err != nil {
		return nil, err
	}

	res, err := svc.fileMapper.mapOne(source, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type FileMoveManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileMoveManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *FileService) MoveMany(opts FileMoveManyOptions, userID string) (*FileMoveManyResult, error) {
	res := &FileMoveManyResult{
		Failed:    make([]string, 0),
		Succeeded: make([]string, 0),
	}
	for _, id := range opts.SourceIDs {
		if _, err := svc.MoveOne(id, opts.TargetID, userID); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

func (svc *FileService) PatchName(id string, name string, userID string) (*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetParentID() != nil {
		existing, err := svc.getChildWithName(*file.GetParentID(), name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errorpkg.NewFileWithSimilarNameExistsError()
		}
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	file.SetName(name)
	if err = svc.fileRepo.Save(file); err != nil {
		return nil, err
	}
	if err := svc.sync(file); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ReprocessResponse struct {
	Accepted []string `json:"accepted"`
	Rejected []string `json:"rejected"`
}

func (r *ReprocessResponse) AppendAccepted(id string) {
	if !slices.Contains(r.Accepted, id) {
		r.Accepted = append(r.Accepted, id)
	}
}

func (r *ReprocessResponse) AppendRejected(id string) {
	if !slices.Contains(r.Rejected, id) {
		r.Rejected = append(r.Rejected, id)
	}
}

func (svc *FileService) Reprocess(id string, userID string) (res *ReprocessResponse, err error) {
	res = &ReprocessResponse{
		// We intend to send an empty array to the caller, better than nil
		Accepted: []string{},
		Rejected: []string{},
	}

	var ancestor model.File
	ancestor, err = svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}

	var tree []model.File
	if ancestor.GetType() == model.FileTypeFolder {
		if err = svc.fileGuard.Authorize(userID, ancestor, model.PermissionViewer); err != nil {
			return nil, err
		}
		tree, err = svc.fileRepo.FindTree(ancestor.GetID())
		if err != nil {
			return nil, err
		}
	} else if ancestor.GetType() == model.FileTypeFile {
		var file model.File
		file, err = svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		tree = append(tree, file)
	}

	for _, file := range tree {
		if file.GetType() != model.FileTypeFile {
			continue
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
			log.GetLogger().Error(err)
			continue
		}
		if !svc.canReprocessFile(file) {
			res.AppendRejected(file.GetID())
			continue
		}

		var snapshot model.Snapshot
		snapshot, err = svc.snapshotCache.Get(*file.GetSnapshotID())
		if err != nil {
			log.GetLogger().Error(err)
			continue
		}
		if !svc.canReprocessSnapshot(snapshot) {
			res.AppendRejected(file.GetID())
			continue
		}

		// Create a task
		var task model.Task
		task, err = svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
			ID:              helper.NewID(),
			Name:            "Waiting.",
			UserID:          userID,
			IsIndeterminate: true,
			Status:          model.TaskStatusWaiting,
			Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
		})
		if err != nil {
			log.GetLogger().Error(err)
			continue
		}
		snapshot.SetTaskID(helper.ToPtr(task.GetID()))
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			log.GetLogger().Error(err)
			continue
		}

		// Forward to conversion microservice
		if err = svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
			TaskID:     task.GetID(),
			SnapshotID: snapshot.GetID(),
			Bucket:     snapshot.GetOriginal().Bucket,
			Key:        snapshot.GetOriginal().Key,
		}); err != nil {
			log.GetLogger().Error(err)
			continue
		} else {
			res.AppendAccepted(file.GetID())
		}
	}
	return res, nil
}

func (svc *FileService) canReprocessFile(file model.File) bool {
	// We don't reprocess if there is no active snapshot
	return file.GetSnapshotID() != nil
}

func (svc *FileService) canReprocessSnapshot(snapshot model.Snapshot) bool {
	// We don't reprocess if there is a pending task
	if snapshot.GetTaskID() != nil {
		task, err := svc.taskCache.Get(*snapshot.GetTaskID())
		if err != nil {
			log.GetLogger().Error(err)
			return false
		}
		if task.GetStatus() == model.TaskStatusWaiting || task.GetStatus() == model.TaskStatusRunning {
			return false
		}
	}
	// We don't reprocess without an "original" on the active snapshot
	if !snapshot.HasOriginal() {
		return false
	}
	return true
}

func (svc *FileService) DeleteOne(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}

	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Deleting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())

	if file.GetParentID() == nil {
		workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		return errorpkg.NewCannotDeleteWorkspaceRootError(file, workspace)
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return err
	}

	treeIDs, err := svc.fileRepo.FindTreeIDs(file.GetID())
	if err != nil {
		return err
	}

	/* Delete file from repo */
	if err := svc.fileRepo.Delete(id); err != nil {
		return err
	}

	go func(ids []string) {
		/* Delete tree from repo */
		const ChunkSize = 1000
		for i := 0; i < len(ids); i += ChunkSize {
			end := i + ChunkSize
			if end > len(ids) {
				end = len(ids)
			}
			chunk := ids[i:end]
			if err := svc.fileRepo.DeleteChunk(chunk); err != nil {
				log.GetLogger().Error(err)
			}
		}
		/* Delete snapshot mappings from tree */
		if err := svc.snapshotRepo.DeleteMappingsForTree(id); err != nil {
			log.GetLogger().Error(err)
		}
		/* Fetch dangling snapshots */
		var danglingSnapshots []model.Snapshot
		danglingSnapshots, err = svc.snapshotRepo.FindAllDangling()
		if err != nil {
			log.GetLogger().Error(err)
		}
		/* Delete dangling snapshots from S3 */
		for _, s := range danglingSnapshots {
			if s.HasOriginal() {
				if err = svc.s3.RemoveObject(s.GetOriginal().Key, s.GetOriginal().Bucket, minio.RemoveObjectOptions{}); err != nil {
					log.GetLogger().Error(err)
				}
			}
			if s.HasPreview() {
				if err = svc.s3.RemoveObject(s.GetPreview().Key, s.GetPreview().Bucket, minio.RemoveObjectOptions{}); err != nil {
					log.GetLogger().Error(err)
				}
			}
			if s.HasText() {
				if err = svc.s3.RemoveObject(s.GetText().Key, s.GetText().Bucket, minio.RemoveObjectOptions{}); err != nil {
					log.GetLogger().Error(err)
				}
			}
			if s.HasThumbnail() {
				if err = svc.s3.RemoveObject(s.GetThumbnail().Key, s.GetThumbnail().Bucket, minio.RemoveObjectOptions{}); err != nil {
					log.GetLogger().Error(err)
				}
			}
			if err := svc.snapshotCache.Delete(s.GetID()); err != nil {
				log.GetLogger().Error(err)
			}
		}
		/* Delete dangling snapshots from cache */
		for _, s := range danglingSnapshots {
			if err = svc.snapshotCache.Delete(s.GetID()); err != nil {
				log.GetLogger().Error(err)
			}
		}
		/* Delete dangling snapshots from repo */
		if err = svc.snapshotRepo.DeleteAllDangling(); err != nil {
			log.GetLogger().Error(err)
		}
	}(treeIDs)

	/* Delete from cache */
	go func(ids []string) {
		for _, treeID := range treeIDs {
			if err := svc.fileCache.Delete(treeID); err != nil {
				// Here we intentionally don't return an error or panic, we just print the error
				log.GetLogger().Error(err)
			}
		}
	}(treeIDs)

	/* Delete from search */
	go func(ids []string) {
		if err := svc.fileSearch.Delete(treeIDs); err != nil {
			// Here we intentionally don't return an error or panic, we just print the error
			log.GetLogger().Error(err)
		}
	}(treeIDs)

	return nil
}

type FileDeleteManyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileDeleteManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *FileService) DeleteMany(opts FileDeleteManyOptions, userID string) (*FileDeleteManyResult, error) {
	res := &FileDeleteManyResult{
		Failed:    make([]string, 0),
		Succeeded: make([]string, 0),
	}
	for _, id := range opts.IDs {
		if err := svc.DeleteOne(id, userID); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

func (svc *FileService) ComputeSize(id string, userID string) (*int64, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.fileRepo.ComputeSize(id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *FileService) Count(id string, userID string) (*int64, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.fileRepo.CountItems(id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *FileService) GrantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		if _, err := svc.userRepo.Find(assigneeID); err != nil {
			return err
		}
		if err = svc.fileRepo.GrantUserPermission(id, assigneeID, permission); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		workspace, err := svc.workspaceRepo.Find(file.GetWorkspaceID())
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
			if err := svc.sync(f); err != nil {
				return err
			}
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		for _, f := range tree {
			if err := svc.sync(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FileService) RevokeUserPermission(ids []string, assigneeID string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		if _, err := svc.userRepo.Find(assigneeID); err != nil {
			return err
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeUserPermission(tree, assigneeID); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		for _, f := range tree {
			if _, err := svc.fileCache.Refresh(f.GetID()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FileService) GrantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		group, err := svc.groupCache.Get(groupID)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return err
		}
		if err = svc.fileRepo.GrantGroupPermission(id, groupID, permission); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		workspace, err := svc.workspaceRepo.Find(file.GetWorkspaceID())
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
			if err := svc.sync(f); err != nil {
				return err
			}
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		for _, f := range tree {
			if err := svc.sync(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *FileService) RevokeGroupPermission(ids []string, groupID string, userID string) error {
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return err
		}
		if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return err
		}
		group, err := svc.groupCache.Get(groupID)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
			return err
		}
		tree, err := svc.fileRepo.FindTree(id)
		if err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeGroupPermission(tree, groupID); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
		for _, f := range tree {
			if _, err := svc.fileCache.Refresh(f.GetID()); err != nil {
				return err
			}
		}
	}
	return nil
}

type UserPermission struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	Permission string `json:"permission"`
}

func (svc *FileService) FindUserPermissions(id string, userID string) ([]*UserPermission, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.FindUserPermissions(id)
	if err != nil {
		return nil, err
	}
	res := make([]*UserPermission, 0)
	for _, p := range permissions {
		if p.GetUserID() == userID {
			continue
		}
		u, err := svc.userRepo.Find(p.GetUserID())
		if err != nil {
			return nil, err
		}
		res = append(res, &UserPermission{
			ID:         p.GetID(),
			User:       svc.userMapper.mapOne(u),
			Permission: p.GetPermission(),
		})
	}
	return res, nil
}

type GroupPermission struct {
	ID         string `json:"id"`
	Group      *Group `json:"group"`
	Permission string `json:"permission"`
}

func (svc *FileService) FindGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.FindGroupPermissions(id)
	if err != nil {
		return nil, err
	}
	res := make([]*GroupPermission, 0)
	for _, p := range permissions {
		m, err := svc.groupCache.Get(p.GetGroupID())
		if err != nil {
			return nil, err
		}
		g, err := svc.groupMapper.mapOne(m, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, &GroupPermission{
			ID:         p.GetID(),
			Group:      g,
			Permission: p.GetPermission(),
		})
	}
	return res, nil
}

func (svc *FileService) doAuthorization(data []model.File, userID string) ([]model.File, error) {
	var res []model.File
	for _, f := range data {
		if svc.fileGuard.IsAuthorized(userID, f, model.PermissionViewer) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *FileService) doAuthorizationByIDs(ids []string, userID string) ([]model.File, error) {
	var res []model.File
	for _, id := range ids {
		var f model.File
		f, err := svc.fileCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if svc.fileGuard.IsAuthorized(userID, f, model.PermissionViewer) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *FileService) doSorting(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
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
			fileA, err := svc.fileMapper.mapOne(data[i], userID)
			if err != nil {
				return false
			}
			fileB, err := svc.fileMapper.mapOne(data[j], userID)
			if err != nil {
				return false
			}
			var sizeA int64 = 0
			if fileA.Snapshot != nil && fileA.Snapshot.Original != nil {
				sizeA = *fileA.Snapshot.Original.Size
			}
			var sizeB int64 = 0
			if fileB.Snapshot != nil && fileB.Snapshot.Original != nil {
				sizeB = *fileB.Snapshot.Original.Size
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
			Filter(func(file interface{}) bool {
				f, err := svc.fileMapper.mapOne(file.(model.File), userID)
				if err != nil {
					return false
				}
				if f.Snapshot != nil && f.Snapshot.Original == nil {
					return false
				}
				if f.Snapshot != nil && svc.fileIdent.IsImage(f.Snapshot.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		pdfs, _ := rxgo.Just(files)().
			Filter(func(file interface{}) bool {
				f, err := svc.fileMapper.mapOne(file.(model.File), userID)
				if err != nil {
					return false
				}
				if f.Snapshot != nil && f.Snapshot.Original == nil {
					return false
				}
				if f.Snapshot != nil && svc.fileIdent.IsPDF(f.Snapshot.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		documents, _ := rxgo.Just(files)().
			Filter(func(file interface{}) bool {
				f, err := svc.fileMapper.mapOne(file.(model.File), userID)
				if err != nil {
					return false
				}
				if f.Snapshot != nil && f.Snapshot.Original == nil {
					return false
				}
				if f.Snapshot != nil && svc.fileIdent.IsOffice(f.Snapshot.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		videos, _ := rxgo.Just(files)().
			Filter(func(file interface{}) bool {
				f, err := svc.fileMapper.mapOne(file.(model.File), userID)
				if err != nil {
					return false
				}
				if f.Snapshot != nil && f.Snapshot.Original == nil {
					return false
				}
				if f.Snapshot != nil && svc.fileIdent.IsVideo(f.Snapshot.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		texts, _ := rxgo.Just(files)().
			Filter(func(file interface{}) bool {
				f, err := svc.fileMapper.mapOne(file.(model.File), userID)
				if err != nil {
					return false
				}
				if f.Snapshot != nil && f.Snapshot.Original == nil {
					return false
				}
				if f.Snapshot != nil && svc.fileIdent.IsPlainText(f.Snapshot.Original.Extension) {
					return true
				}
				return false
			}).
			ToSlice(0)
		others, _ := rxgo.Just(files)().
			Filter(func(file interface{}) bool {
				f, err := svc.fileMapper.mapOne(file.(model.File), userID)
				if err != nil {
					return false
				}
				if f.Snapshot != nil && f.Snapshot.Original == nil {
					return false
				}
				if f.Snapshot != nil &&
					!svc.fileIdent.IsImage(f.Snapshot.Original.Extension) &&
					!svc.fileIdent.IsPDF(f.Snapshot.Original.Extension) &&
					!svc.fileIdent.IsOffice(f.Snapshot.Original.Extension) &&
					!svc.fileIdent.IsVideo(f.Snapshot.Original.Extension) &&
					!svc.fileIdent.IsPlainText(f.Snapshot.Original.Extension) {
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

func (svc *FileService) doPagination(data []model.File, page, size uint64) (pageData []model.File, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.File{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *FileService) doQueryFiltering(data []model.File, opts FileQuery, parent model.File) ([]model.File, error) {
	filtered, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetWorkspaceID() == parent.GetWorkspaceID()
		}).
		Filter(func(v interface{}) bool {
			if opts.Type != nil {
				return v.(model.File).GetType() == *opts.Type
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			file := v.(model.File)
			res, err := svc.fileRepo.IsGrandChildOf(file.GetID(), parent.GetID())
			if err != nil {
				return false
			}
			return res
		}).
		Filter(func(v interface{}) bool {
			if opts.CreateTimeBefore != nil {
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return t.UnixMilli() >= *opts.CreateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.CreateTimeBefore != nil {
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return t.UnixMilli() <= *opts.CreateTimeBefore
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeAfter != nil {
				file := v.(model.File)
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return file.GetUpdateTime() != nil && t.UnixMilli() >= *opts.UpdateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeBefore != nil {
				file := v.(model.File)
				t, _ := time.Parse(time.RFC3339, v.(model.File).GetCreateTime())
				return file.GetUpdateTime() != nil && t.UnixMilli() <= *opts.UpdateTimeBefore
			} else {
				return true
			}
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range filtered {
		var file model.File
		file, err := svc.fileCache.Get(v.(model.File).GetID())
		if err != nil {
			return nil, err
		}
		res = append(res, file)
	}
	return res, nil
}

func (svc *FileService) getChildWithName(id string, name string) (model.File, error) {
	children, err := svc.fileRepo.FindChildren(id)
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		if child.GetName() == name {
			return child, nil
		}
	}
	return nil, nil
}

func (svc *FileService) sync(file model.File) error {
	if err := svc.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	if err := svc.fileCache.Set(file); err != nil {
		return err
	}
	return nil
}

type FileMapper struct {
	groupCache     *cache.GroupCache
	snapshotMapper *SnapshotMapper
	snapshotCache  *cache.SnapshotCache
	snapshotRepo   repo.SnapshotRepo
	config         *config.Config
}

func NewFileMapper() *FileMapper {
	return &FileMapper{
		groupCache:     cache.NewGroupCache(),
		snapshotMapper: NewSnapshotMapper(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		config:         config.GetConfig(),
	}
}

func (mp *FileMapper) mapOne(m model.File, userID string) (*File, error) {
	res := &File{
		ID:          m.GetID(),
		WorkspaceID: m.GetWorkspaceID(),
		Name:        m.GetName(),
		Type:        m.GetType(),
		ParentID:    m.GetParentID(),
		CreateTime:  m.GetCreateTime(),
		UpdateTime:  m.GetUpdateTime(),
	}
	if m.GetSnapshotID() != nil {
		snapshot, err := mp.snapshotCache.Get(*m.GetSnapshotID())
		if err != nil {
			return nil, err
		}
		res.Snapshot = mp.snapshotMapper.mapOne(snapshot)
		res.Snapshot.IsActive = true
	}
	res.Permission = model.PermissionNone
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
		for _, u := range g.GetMembers() {
			if u == userID && model.GetPermissionWeight(p.GetValue()) > model.GetPermissionWeight(res.Permission) {
				res.Permission = p.GetValue()
			}
		}
	}
	shareCount := 0
	for _, p := range m.GetUserPermissions() {
		if p.GetUserID() != userID {
			shareCount++
		}
	}
	if res.Permission == model.PermissionOwner {
		shareCount += len(m.GetGroupPermissions())
		res.IsShared = new(bool)
		if shareCount > 0 {
			*res.IsShared = true
		} else {
			*res.IsShared = false
		}
	}
	return res, nil
}

func (mp *FileMapper) mapMany(data []model.File, userID string) ([]*File, error) {
	res := make([]*File, 0)
	for _, file := range data {
		f, err := mp.mapOne(file, userID)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, f)
	}
	return res, nil
}
