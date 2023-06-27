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
	WorkspaceID string      `json:"workspaceId"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	ParentID    *string     `json:"parentId,omitempty"`
	Version     *int64      `json:"version,omitempty"`
	Original    *Download   `json:"original,omitempty"`
	Preview     *Download   `json:"preview,omitempty"`
	OCR         *Download   `json:"ocr,omitempty"`
	Thumbnail   *Thumbnail  `json:"thumbnail,omitempty"`
	Language    *string     `json:"language,omitempty"`
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
	WorkspaceID      string  `json:"workspaceId" validate:"required"`
	ParentID         *string `json:"parentId,omitempty"`
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
	WorkspaceID string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	Type        string  `json:"type" validate:"required,oneof=file folder"`
	ParentID    *string `json:"parentId" validate:"required"`
}

type FileCreateFolderOptions struct {
	WorkspaceID string  `json:"workspaceId" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	ParentID    *string `json:"parentId"`
}

type FileListByIDOptions struct {
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
	FileType  string
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
	OCR       *Download  `json:"ocr,omitempty"`
	Thumbnail *Thumbnail `json:"thumbnail,omitempty"`
	Language  *string    `json:"language,omitempty"`
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

type SnapshotUpdateOptions struct {
	Options   infra.RunPipelineOptions `json:"options"`
	Original  *model.S3Object          `json:"original,omitempty"`
	Preview   *model.S3Object          `json:"preview,omitempty"`
	Text      *model.S3Object          `json:"text,omitempty"`
	OCR       *model.S3Object          `json:"ocr,omitempty"`
	Thumbnail *model.Thumbnail         `json:"thumbnail,omitempty"`
	Language  *string                  `json:"language,omitempty"`
}

type FileService struct {
	fileRepo         repo.FileRepo
	fileSearch       *search.FileSearch
	fileGuard        *guard.FileGuard
	fileMapper       *FileMapper
	fileCache        *cache.FileCache
	workspaceCache   *cache.WorkspaceCache
	workspaceRepo    repo.WorkspaceRepo
	workspaceGuard   *guard.WorkspaceGuard
	workspaceSvc     *WorkspaceService
	snapshotRepo     repo.SnapshotRepo
	userRepo         repo.UserRepo
	userMapper       *userMapper
	groupCache       *cache.GroupCache
	groupGuard       *guard.GroupGuard
	groupMapper      *groupMapper
	permissionRepo   repo.PermissionRepo
	fileIdentifier   *infra.FileIdentifier
	s3               *infra.S3Manager
	conversionClient *infra.ConversionClient
	config           config.Config
}

func NewFileService() *FileService {
	return &FileService{
		fileRepo:         repo.NewFileRepo(),
		fileCache:        cache.NewFileCache(),
		fileSearch:       search.NewFileSearch(),
		fileGuard:        guard.NewFileGuard(),
		fileMapper:       NewFileMapper(),
		workspaceGuard:   guard.NewWorkspaceGuard(),
		workspaceCache:   cache.NewWorkspaceCache(),
		workspaceRepo:    repo.NewWorkspaceRepo(),
		workspaceSvc:     NewWorkspaceService(),
		snapshotRepo:     repo.NewSnapshotRepo(),
		userRepo:         repo.NewUserRepo(),
		userMapper:       newUserMapper(),
		groupCache:       cache.NewGroupCache(),
		groupGuard:       guard.NewGroupGuard(),
		groupMapper:      newGroupMapper(),
		permissionRepo:   repo.NewPermissionRepo(),
		fileIdentifier:   infra.NewFileIdentifier(),
		s3:               infra.NewS3Manager(),
		conversionClient: infra.NewConversionClient(),
		config:           config.GetConfig(),
	}
}

func (svc *FileService) Create(opts FileCreateOptions, userID string) (*File, error) {
	if len(*opts.ParentID) > 0 {
		if err := svc.validateParent(*opts.ParentID, userID); err != nil {
			return nil, err
		}
		nameExists, err := svc.hasChildWithName(*opts.ParentID, opts.Name)
		if err != nil {
			return nil, err
		}
		if nameExists {
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
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) validateParent(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
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

func (svc *FileService) Store(fileID string, filePath string, userID string) (*File, error) {
	file, err := svc.fileRepo.Find(fileID)
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
	latestVersion, err := svc.snapshotRepo.GetLatestVersionForFile(fileID)
	if err != nil {
		return nil, err
	}
	snapshotID := helper.NewID()
	snapshot := repo.NewSnapshot()
	snapshot.SetID(snapshotID)
	snapshot.SetVersion(latestVersion)
	if err = svc.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if err = svc.snapshotRepo.MapWithFile(snapshotID, fileID); err != nil {
		return nil, err
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	original := model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    fileID + "/" + snapshotID + "/original" + strings.ToLower(filepath.Ext(filePath)),
		Size:   stat.Size(),
	}
	if err = svc.s3.PutFile(original.Key, filePath, infra.DetectMimeFromFile(filePath), workspace.GetBucket()); err != nil {
		return nil, err
	}
	snapshot.SetOriginal(&original)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	file, err = svc.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	if err := svc.conversionClient.RunPipeline(&infra.RunPipelineOptions{
		FileID:     file.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     original.Bucket,
		Key:        original.Key,
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) UpdateSnapshot(opts SnapshotUpdateOptions, apiKey string) error {
	if apiKey != svc.config.Security.APIKey {
		return errorpkg.NewInvalidAPIKeyError()
	}
	if err := svc.snapshotRepo.Update(opts.Options.SnapshotID, repo.SnapshotUpdateOptions{
		Thumbnail: opts.Thumbnail,
		Original:  opts.Original,
		Preview:   opts.Preview,
		Text:      opts.Text,
		OCR:       opts.OCR,
		Language:  opts.Language,
	}); err != nil {
		return err
	}
	file, err := svc.fileCache.Refresh(opts.Options.FileID)
	if err != nil {
		return err
	}
	if err = svc.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	return nil
}

func (svc *FileService) DownloadOriginalBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userID)
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

func (svc *FileService) DownloadPreviewBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userID)
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

func (svc *FileService) DownloadOCRBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userID)
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
	if latestSnapshot.HasOCR() {
		ocr := latestSnapshot.GetOCR()
		buf, err := svc.s3.GetObject(ocr.Key, ocr.Bucket)
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, latestSnapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *FileService) FindByID(ids []string, userID string) ([]*File, error) {
	user, err := svc.userRepo.Find(userID)
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
		authorized, err := svc.doAuthorizationByIDs(ids, user)
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
	result, err := svc.FindByID([]string{currentID}, userID)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (svc *FileService) ListByPath(path string, userID string) ([]*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
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
		authorized, err := svc.doAuthorizationByIDs(ids, user)
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
		authorized, err := svc.doAuthorizationByIDs(ids, user)
		if err != nil {
			return nil, err
		}
		result, err := svc.fileMapper.mapMany(authorized, userID)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else if currentType == model.FileTypeFile {
		result, err := svc.FindByID([]string{currentID}, userID)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("invalid file type %s", currentType))
	}
}

func (svc *FileService) ListByID(id string, opts FileListByIDOptions, userID string) (*FileList, error) {
	user, err := svc.userRepo.Find(userID)
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
	authorized, err := svc.doAuthorizationByIDs(ids, user)
	if err != nil {
		return nil, err
	}
	var filteredFiles []model.File
	for _, f := range authorized {
		if opts.FileType == "" || f.GetType() == opts.FileType {
			filteredFiles = append(filteredFiles, f)
		}
	}
	sorted := svc.doSorting(filteredFiles, opts.SortBy, opts.SortOrder, userID)
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

func (svc *FileService) Search(opts FileSearchOptions, page uint, size uint, userID string) (*FileSearchResult, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceRepo.Find(opts.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceGuard.Authorize(user, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	files, err := svc.fileSearch.Query(opts.Text)
	if err != nil {
		return nil, err
	}
	filtered, err := svc.doFiltering(opts, files, userID)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.doAuthorization(filtered, user)
	if err != nil {
		return nil, err
	}
	paged, totalElements, totalPages := svc.doPagination(authorized, page, size)
	v, err := svc.fileMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	res := &FileSearchResult{
		Data:          v,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          page,
		Size:          size,
		Query:         opts,
	}
	return res, nil
}

func (svc *FileService) GetPath(id string, userID string) ([]*File, error) {
	user, err := svc.userRepo.Find(userID)
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
		v, err := svc.fileMapper.mapOne(f, userID)
		if err != nil {
			return nil, err
		}
		res = append([]*File{v}, res...)
	}
	return res, nil
}

func (svc *FileService) Copy(targetID string, sourceIDs []string, userID string) (copiedFiles []*File, err error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
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
		var permissions []*repo.UserPermission
		for i, o := range sourceTree {
			c := repo.NewFile()
			c.SetID(helper.NewID())
			c.SetParentID(o.GetParentID())
			c.SetWorkspaceID(o.GetWorkspaceID())
			c.SetType(o.GetType())
			c.SetName(o.GetName())
			c.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
			if o.GetID() == sourceID {
				rootCloneIndex = i
			}
			cloneIDs[o.GetID()] = c.GetID()
			originalIDs[c.GetID()] = o.GetID()
			clones = append(clones, c)
			permissions = append(permissions, &repo.UserPermission{
				ID:         helper.NewID(),
				UserID:     userID,
				ResourceID: c.GetID(),
				Permission: model.PermissionOwner,
				CreateTime: time.Now().UTC().Format(time.RFC3339),
			})
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
		nameExists, err := svc.hasChildWithName(targetID, rootClone.GetName())
		if err != nil {
			return nil, err
		}
		if nameExists {
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

		/* Assign snapshots to clones */
		for _, c := range clones {
			if err := svc.fileRepo.AssignSnapshots(c.GetID(), originalIDs[c.GetID()]); err != nil {
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
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return []string{}, err
	}
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
			nameExists, err := svc.hasChildWithName(targetID, source.GetName())
			if err != nil {
				return nil, err
			}
			if nameExists {
				return nil, errorpkg.NewFileWithSimilarNameExistsError()
			}
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
	return parentIDs, nil
}

func (svc *FileService) Rename(id string, name string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetParentID() != nil {
		nameExists, err := svc.hasChildWithName(*file.GetParentID(), name)
		if err != nil {
			return nil, err
		}
		if nameExists {
			return nil, errorpkg.NewFileWithSimilarNameExistsError()
		}
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
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileService) Delete(ids []string, userID string) ([]string, error) {
	var res []string
	for _, id := range ids {
		var user model.User
		user, err := svc.userRepo.Find(userID)
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
		var treeIDs []string
		for _, f := range tree {
			treeIDs = append(treeIDs, f.GetID())
		}
		if err := svc.fileSearch.Delete(treeIDs); err != nil {
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
			if s.HasText() {
				if err = svc.s3.RemoveObject(s.GetText().Key, s.GetText().Bucket); err != nil {
					return nil, err
				}
			}
			if s.HasOCR() {
				if err = svc.s3.RemoveObject(s.GetOCR().Key, s.GetOCR().Bucket); err != nil {
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

func (svc *FileService) GetSize(id string, userID string) (int64, error) {
	user, err := svc.userRepo.Find(userID)
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

func (svc *FileService) GetItemCount(id string, userID string) (int64, error) {
	user, err := svc.userRepo.Find(userID)
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

func (svc *FileService) GrantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	user, err := svc.userRepo.Find(userID)
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

func (svc *FileService) RevokeUserPermission(ids []string, assigneeID string, userID string) error {
	user, err := svc.userRepo.Find(userID)
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
		if _, err := svc.userRepo.Find(assigneeID); err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeUserPermission(id, assigneeID); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *FileService) GrantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	user, err := svc.userRepo.Find(userID)
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
		group, err := svc.groupCache.Get(groupID)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
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

func (svc *FileService) RevokeGroupPermission(ids []string, groupID string, userID string) error {
	user, err := svc.userRepo.Find(userID)
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
		group, err := svc.groupCache.Get(groupID)
		if err != nil {
			return err
		}
		if err := svc.groupGuard.Authorize(user, group, model.PermissionViewer); err != nil {
			return err
		}
		if err := svc.fileRepo.RevokeGroupPermission(id, groupID); err != nil {
			return err
		}
		if _, err := svc.fileCache.Refresh(file.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *FileService) GetUserPermissions(id string, userID string) ([]*UserPermission, error) {
	user, err := svc.userRepo.Find(userID)
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
		if p.UserID == userID {
			continue
		}
		u, err := svc.userRepo.Find(p.UserID)
		if err != nil {
			return nil, err
		}
		res = append(res, &UserPermission{
			ID:         p.ID,
			User:       svc.userMapper.mapOne(u),
			Permission: p.Permission,
		})
	}
	return res, nil
}

func (svc *FileService) GetGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
	user, err := svc.userRepo.Find(userID)
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
		g, err := svc.groupMapper.mapOne(m, userID)
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

func (svc *FileService) doAuthorization(data []model.File, user model.User) ([]model.File, error) {
	var res []model.File
	for _, f := range data {
		if svc.fileGuard.IsAuthorized(user, f, model.PermissionViewer) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *FileService) doAuthorizationByIDs(ids []string, user model.User) ([]model.File, error) {
	var res []model.File
	for _, id := range ids {
		var f model.File
		f, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		if svc.fileGuard.IsAuthorized(user, f, model.PermissionViewer) {
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
				f, err := svc.fileMapper.mapOne(v.(model.File), userID)
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
				f, err := svc.fileMapper.mapOne(v.(model.File), userID)
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
				f, err := svc.fileMapper.mapOne(v.(model.File), userID)
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
				f, err := svc.fileMapper.mapOne(v.(model.File), userID)
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
				f, err := svc.fileMapper.mapOne(v.(model.File), userID)
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
				f, err := svc.fileMapper.mapOne(v.(model.File), userID)
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

func (svc *FileService) doPagination(data []model.File, page, size uint) ([]model.File, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		page = totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

func (svc *FileService) doFiltering(opts FileSearchOptions, data []model.File, userID string) ([]model.File, error) {
	filtered, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetWorkspaceID() == opts.WorkspaceID
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
			if opts.ParentID != nil {
				res, err := svc.fileRepo.IsGrandChildOf(file.GetID(), *opts.ParentID)
				if err != nil {
					return false
				}
				return res
			} else {
				return true
			}
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

func (svc *FileService) hasChildWithName(id string, name string) (bool, error) {
	children, err := svc.fileRepo.FindChildren(id)
	if err != nil {
		return false, err
	}
	for _, child := range children {
		if child.GetName() == name {
			return true, nil
		}
	}
	return false, nil
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

func (mp *FileMapper) mapOne(m model.File, userID string) (*File, error) {
	snapshots := m.GetSnapshots()
	res := &File{
		ID:          m.GetID(),
		WorkspaceID: m.GetWorkspaceID(),
		Name:        m.GetName(),
		Type:        m.GetType(),
		ParentID:    m.GetParentID(),
		Snapshots:   mp.mapSnapshots(snapshots, m.GetID()),
		CreateTime:  m.GetCreateTime(),
		UpdateTime:  m.GetUpdateTime(),
	}
	if len(snapshots) > 0 {
		latest := mp.mapSnapshot(snapshots[len(snapshots)-1])
		res.Version = &latest.Version
		res.Original = latest.Original
		res.Preview = latest.Preview
		res.OCR = latest.OCR
		res.Thumbnail = latest.Thumbnail
		res.Language = latest.Language
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
	shareCount += len(m.GetGroupPermissions())
	if shareCount > 0 {
		res.IsShared = true
	} else {
		res.IsShared = false
	}
	return res, nil
}

func (mp *FileMapper) mapMany(files []model.File, userID string) ([]*File, error) {
	res := make([]*File, 0)
	for _, f := range files {
		v, err := mp.mapOne(f, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func (mp *FileMapper) mapSnapshot(m model.Snapshot) *Snapshot {
	s := &Snapshot{
		ID:      m.GetID(),
		Version: m.GetVersion(),
	}
	if m.HasOriginal() {
		s.Original = mp.mapOriginal(m.GetOriginal())
	}
	if m.HasPreview() {
		s.Preview = mp.mapPreview(m.GetPreview())
	}
	if m.HasOCR() {
		s.OCR = mp.mapOCR(m.GetOCR())
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.mapThumbnail(m.GetThumbnail())
	}
	if m.HasLanguage() {
		s.Language = m.GetLanguage()
	}
	return s
}

func (mp *FileMapper) mapOriginal(m *model.S3Object) *Download {
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

func (mp *FileMapper) mapPreview(m *model.S3Object) *Download {
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

func (mp *FileMapper) mapOCR(m *model.S3Object) *Download {
	return &Download{
		Extension: filepath.Ext(m.Key),
		Size:      m.Size,
	}
}

func (mp *FileMapper) mapThumbnail(m *model.Thumbnail) *Thumbnail {
	return &Thumbnail{
		Base64: m.Base64,
		Width:  m.Width,
		Height: m.Height,
	}
}

func (mp *FileMapper) mapSnapshots(snapshots []model.Snapshot, fileID string) []*Snapshot {
	res := make([]*Snapshot, 0)
	for _, s := range snapshots {
		res = append(res, mp.mapSnapshot(s))
	}
	return res
}
