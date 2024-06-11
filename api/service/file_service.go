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
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	"github.com/minio/minio-go/v7"
	"github.com/reactivex/rxgo/v2"
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
	taskSvc        *TaskService
	fileIdent      *infra.FileIdentifier
	s3             *infra.S3Manager
	pipelineClient *client.PipelineClient
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
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotSvc:    NewSnapshotService(),
		userRepo:       repo.NewUserRepo(),
		userMapper:     newUserMapper(),
		groupCache:     cache.NewGroupCache(),
		groupGuard:     guard.NewGroupGuard(),
		groupMapper:    newGroupMapper(),
		permissionRepo: repo.NewPermissionRepo(),
		taskSvc:        NewTaskService(),
		fileIdent:      infra.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		pipelineClient: client.NewPipelineClient(),
		config:         config.GetConfig(),
	}
}

type FileCreateOptions struct {
	WorkspaceID string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	Type        string  `json:"type" validate:"required,oneof=file folder"`
	ParentID    *string `json:"parentId" validate:"required"`
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

func (svc *FileService) Store(id string, path string, userID string) (*File, error) {
	file, err := svc.fileRepo.Find(id)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	latestVersion, err := svc.snapshotRepo.GetLatestVersionForFile(id)
	if err != nil {
		return nil, err
	}
	snapshotID := helper.NewID()
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
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	exceedsProcessingLimit := stat.Size() > helper.MegabyteToByte(svc.config.Limits.FileProcessingMaxSizeMB)
	original := model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    snapshotID + "/original" + strings.ToLower(filepath.Ext(path)),
		Size:   helper.ToPtr(stat.Size()),
	}
	if err = svc.s3.PutFile(original.Key, path, infra.DetectMimeFromPath(path), workspace.GetBucket(), minio.PutObjectOptions{}); err != nil {
		return nil, err
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
			Payload:         map[string]string{"fileId": file.GetID()},
		})
		if err != nil {
			return nil, err
		}
		snapshot.SetTaskID(helper.ToPtr(task.GetID()))
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			return nil, err
		}
		if err := svc.pipelineClient.Run(&client.PipelineRunOptions{
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

func (svc *FileService) DownloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasWatermark() {
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
			return nil, err
		}
	} else {
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
	}
	if snapshot.HasOriginal() {
		objectInfo, err := svc.s3.StatObject(snapshot.GetOriginal().Key, snapshot.GetOriginal().Bucket, minio.StatObjectOptions{})
		if err != nil {
			return nil, err
		}
		opts := minio.GetObjectOptions{}
		var ri *infra.RangeInterval
		if rangeHeader != "" {
			ri = infra.NewRangeInterval(rangeHeader, objectInfo.Size)
			ri.ApplyToMinIOGetObjectOptions(&opts)
		}
		if _, err := svc.s3.GetObjectWithBuffer(snapshot.GetOriginal().Key, snapshot.GetOriginal().Bucket, buf, opts); err != nil {
			return nil, err
		}
		return &DownloadResult{
			Buffer:   buf,
			File:     file,
			Snapshot: snapshot,
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasWatermark() {
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
			return nil, err
		}
	} else {
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
	}
	if snapshot.HasPreview() {
		objectInfo, err := svc.s3.StatObject(snapshot.GetOriginal().Key, snapshot.GetOriginal().Bucket, minio.StatObjectOptions{})
		if err != nil {
			return nil, err
		}
		opts := minio.GetObjectOptions{}
		var ri *infra.RangeInterval
		if rangeHeader != "" {
			ri = infra.NewRangeInterval(rangeHeader, objectInfo.Size)
			ri.ApplyToMinIOGetObjectOptions(&opts)
		}
		if _, err := svc.s3.GetObjectWithBuffer(snapshot.GetPreview().Key, snapshot.GetPreview().Bucket, buf, opts); err != nil {
			return nil, err
		}
		return &DownloadResult{
			RangeInterval: ri,
			Buffer:        buf,
			File:          file,
			Snapshot:      snapshot,
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) DownloadThumbnailBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, nil, err
	}
	if snapshot.HasThumbnail() {
		buf, _, err := svc.s3.GetObject(snapshot.GetThumbnail().Key, snapshot.GetThumbnail().Bucket, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, snapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
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
	components := []string{}
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
		ids, err := svc.fileRepo.GetChildrenIDs(currentID)
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
		workspaces, err := svc.workspaceSvc.findAll(userID)
		if err != nil {
			return nil, err
		}
		result := []*File{}
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
	components := []string{}
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
		ids, err := svc.fileRepo.GetChildrenIDs(currentID)
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
		ids, err := svc.fileRepo.GetChildrenIDs(currentID)
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
	Text             string  `json:"text" validate:"required"`
	Type             *string `json:"type,omitempty" validate:"omitempty,oneof=file folder"`
	CreateTimeAfter  *int64  `json:"createTimeAfter,omitempty"`
	CreateTimeBefore *int64  `json:"createTimeBefore,omitempty"`
	UpdateTimeAfter  *int64  `json:"updateTimeAfter,omitempty"`
	UpdateTimeBefore *int64  `json:"updateTimeBefore,omitempty"`
}

type FileList struct {
	Data          []*File    `json:"data"`
	TotalPages    uint       `json:"totalPages"`
	TotalElements uint       `json:"totalElements"`
	Page          uint       `json:"page"`
	Size          uint       `json:"size"`
	Query         *FileQuery `json:"query,omitempty"`
}

type FileListOptions struct {
	Page      uint
	Size      uint
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
	if opts.Page < 1 {
		return nil, errorpkg.NewInvalidPageParameterError()
	}
	if opts.Size < 1 {
		return nil, errorpkg.NewInvalidSizeParameterError()
	}
	ids, err := svc.fileRepo.GetChildrenIDs(id)
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
		data, err = svc.fileSearch.Query(opts.Query.Text)
	} else {
		ids, err := svc.fileRepo.GetChildrenIDs(id)
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
	if err != nil {
		return nil, err
	}
	filtered, err := svc.doQueryFiltering(data, *opts.Query, parent)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorization(filtered, userID)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(authorized, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	v, err := svc.fileMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	res := &FileList{
		Data:          v,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          opts.Page,
		Size:          opts.Size,
		Query:         opts.Query,
	}
	return res, nil
}

func (svc *FileService) GetPath(id string, userID string) ([]*File, error) {
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
	res := []*File{}
	for _, file := range path {
		f, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append([]*File{f}, res...)
	}
	return res, nil
}

func (svc *FileService) Copy(targetID string, sourceIDs []string, userID string) (copiedFiles []*File, err error) {
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return nil, err
	}

	/* Do checks */
	for _, sourceID := range sourceIDs {
		var source model.File
		if source, err = svc.fileCache.Get(sourceID); err != nil {
			return nil, err
		}
		if err = svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
			return nil, err
		}
		if err = svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
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
	for _, sourceID := range sourceIDs {
		/* Get original tree */
		var sourceTree []model.File
		if sourceTree, err = svc.fileRepo.FindTree(sourceID); err != nil {
			return nil, err
		}

		/* Clone source tree */
		var rootCloneIndex int
		var cloneIDs = make(map[string]string)
		var originalIDs = make(map[string]string)
		var clones []model.File
		var permissions []model.UserPermission
		for i, o := range sourceTree {
			c := repo.NewFile()
			c.SetID(helper.NewID())
			c.SetParentID(o.GetParentID())
			c.SetWorkspaceID(o.GetWorkspaceID())
			c.SetSnapshotID(o.GetSnapshotID())
			c.SetType(o.GetType())
			c.SetName(o.GetName())
			c.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
			if o.GetID() == sourceID {
				rootCloneIndex = i
			}
			cloneIDs[o.GetID()] = c.GetID()
			originalIDs[c.GetID()] = o.GetID()
			clones = append(clones, c)

			p := repo.NewUserPermission()
			p.SetID(helper.NewID())
			p.SetUserID(userID)
			p.SetResourceID(c.GetID())
			p.SetPermission(model.PermissionOwner)
			p.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
			permissions = append(permissions, p)
		}

		/* Set parent IDs of clones */
		for i, c := range clones {
			id := cloneIDs[*c.GetParentID()]
			clones[i].SetParentID(&id)
		}

		rootClone := clones[rootCloneIndex]

		/* Parent ID of root clone is target ID */
		if clones != nil {
			rootClone.SetParentID(&targetID)
		}

		/* If there is a file with similar name, append a prefix */
		existing, err := svc.getChildWithName(targetID, rootClone.GetName())
		if err != nil {
			return nil, err
		}
		if existing != nil {
			rootClone.SetName(fmt.Sprintf("Copy of %s", rootClone.GetName()))
		}

		/* Persist clones */
		if err = svc.fileRepo.BulkInsert(clones, 100); err != nil {
			return nil, err
		}

		/* Persist permissions */
		if err = svc.fileRepo.BulkInsertPermissions(permissions, 100); err != nil {
			return nil, err
		}

		/* Attach latest snapshot to clones */
		for _, c := range clones {
			if err := svc.snapshotRepo.Attach(originalIDs[c.GetID()], c.GetID()); err != nil {
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

	copiedFiles, err = svc.fileMapper.mapMany(allClones, userID)
	if err != nil {
		return nil, err
	}

	return copiedFiles, nil
}

func (svc *FileService) Move(targetID string, sourceIDs []string, userID string) (parentIDs []string, err error) {
	parentIDs = []string{}
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return []string{}, err
	}

	/* Do checks */
	for _, id := range sourceIDs {
		source, err := svc.fileCache.Get(id)
		if err != nil {
			return []string{}, err
		}
		if source.GetParentID() != nil {
			existing, err := svc.getChildWithName(targetID, source.GetName())
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return nil, errorpkg.NewFileWithSimilarNameExistsError()
			}
		}
		if err := svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
			return []string{}, err
		}
		if err := svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
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
	for _, id := range sourceIDs {
		source, _ := svc.fileCache.Get(id)

		/* Add old parent */
		parentIDs = append(parentIDs, *source.GetParentID())

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
		parentIDs = append(parentIDs, *source.GetParentID())

		/* Refresh updateTime on source and target */
		timeNow := time.Now().UTC().Format(time.RFC3339)
		source.SetUpdateTime(&timeNow)
		if err := svc.fileRepo.Save(source); err != nil {
			return []string{}, err
		}
		if err := svc.sync(source); err != nil {
			return []string{}, err
		}
		target.SetUpdateTime(&timeNow)
		if err := svc.fileRepo.Save(target); err != nil {
			return []string{}, err
		}
		if err := svc.sync(target); err != nil {
			return []string{}, err
		}
		sourceTree, err := svc.fileRepo.FindTree(source.GetID())
		if err != nil {
			return []string{}, err
		}
		for _, f := range sourceTree {
			if err := svc.sync(f); err != nil {
				return []string{}, err
			}
		}
	}
	return parentIDs, nil
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

func (svc *FileService) Delete(ids []string, userID string) ([]string, error) {
	var res []string
	for _, id := range ids {
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
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
			return nil, err
		}

		// Add parent
		res = append(res, *file.GetParentID())

		var tree []model.File
		tree, err = svc.fileRepo.FindTree(file.GetID())
		if err != nil {
			return nil, err
		}
		var treeIDs []string
		for _, f := range tree {
			treeIDs = append(treeIDs, f.GetID())
		}
		/* Delete from search */
		if err := svc.fileSearch.Delete(treeIDs); err != nil {
			/* Here we intentionally don't return an error or panic, we just print the error,
			that's because we still want to delete the file in the repo afterwards even
			if we fail to delete it from the search. */
			fmt.Println(err)
		}
		/* Delete from cache */
		for _, f := range tree {
			if err = svc.fileCache.Delete(f.GetID()); err != nil {
				// Same thing as above for the search
				fmt.Println(err)
			}
		}
		/* Delete from repo */
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
				if err = svc.s3.RemoveObject(s.GetOriginal().Key, s.GetOriginal().Bucket, minio.RemoveObjectOptions{}); err != nil {
					return nil, err
				}
			}
			if s.HasPreview() {
				if err = svc.s3.RemoveObject(s.GetPreview().Key, s.GetPreview().Bucket, minio.RemoveObjectOptions{}); err != nil {
					return nil, err
				}
			}
			if s.HasText() {
				if err = svc.s3.RemoveObject(s.GetText().Key, s.GetText().Bucket, minio.RemoveObjectOptions{}); err != nil {
					return nil, err
				}
			}
			if err := svc.snapshotCache.Delete(s.GetID()); err != nil {
				return nil, err
			}
		}
		if err = svc.snapshotRepo.DeleteAllDangling(); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *FileService) GetSize(id string, userID string) (*int64, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.fileRepo.GetSize(id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *FileService) GetCount(id string, userID string) (*int64, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	res, err := svc.fileRepo.GetItemCount(id)
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

func (svc *FileService) GetUserPermissions(id string, userID string) ([]*UserPermission, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.GetUserPermissions(id)
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

func (svc *FileService) GetGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	permissions, err := svc.permissionRepo.GetGroupPermissions(id)
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
			return nil, err
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
				sizeA = int64(*fileA.Snapshot.Original.Size)
			}
			var sizeB int64 = 0
			if fileB.Snapshot != nil && fileB.Snapshot.Original != nil {
				sizeB = int64(*fileB.Snapshot.Original.Size)
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
				if f.Snapshot.Original == nil {
					return false
				}
				if svc.fileIdent.IsImage(f.Snapshot.Original.Extension) {
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
				if f.Snapshot.Original == nil {
					return false
				}
				if svc.fileIdent.IsPDF(f.Snapshot.Original.Extension) {
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
				if f.Snapshot.Original == nil {
					return false
				}
				if svc.fileIdent.IsOffice(f.Snapshot.Original.Extension) {
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
				if f.Snapshot.Original == nil {
					return false
				}
				if svc.fileIdent.IsVideo(f.Snapshot.Original.Extension) {
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
				if f.Snapshot.Original == nil {
					return false
				}
				if svc.fileIdent.IsPlainText(f.Snapshot.Original.Extension) {
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
				if f.Snapshot.Original == nil {
					return false
				}
				if !svc.fileIdent.IsImage(f.Snapshot.Original.Extension) &&
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

func (svc *FileService) doPagination(data []model.File, page, size uint) (pageData []model.File, totalElements uint, totalPages uint) {
	totalElements = uint(len(data))
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
	config         config.Config
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
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}
