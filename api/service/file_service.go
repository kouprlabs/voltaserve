// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

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

	"github.com/gosimple/slug"
	"github.com/minio/minio-go/v7"
	"github.com/reactivex/rxgo/v2"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
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
	fileCreate      *fileCreate
	fileStore       *fileStore
	fileDelete      *fileDelete
	fileMove        *fileMove
	fileCopy        *fileCopy
	fileDownload    *fileDownload
	fileFetch       *fileFetch
	fileList        *fileList
	fileSortService *fileSortService
	fileReprocess   *fileReprocess
	filePermission  *filePermission
	fileCompute     *fileCompute
	filePatch       *filePatch
}

func NewFileService() *FileService {
	return &FileService{
		fileCreate:      newFileCreate(),
		fileStore:       newFileStore(),
		fileDelete:      newFileDelete(),
		fileMove:        newFileMove(),
		fileCopy:        newFileCopy(),
		fileDownload:    newFileDownload(),
		fileFetch:       newFileFetch(),
		fileList:        newFileList(),
		fileSortService: newFileSortService(),
		fileReprocess:   newFileReprocess(),
		filePermission:  newFilePermission(),
		fileCompute:     newFileCompute(),
		filePatch:       newFilePatch(),
	}
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

const (
	FileSortByName         = "name"
	FileSortByKind         = "kind"
	FileSortBySize         = "size"
	FileSortByDateCreated  = "date_created"
	FileSortByDateModified = "date_modified"
)

const (
	FileSortOrderAsc  = "asc"
	FileSortOrderDesc = "desc"
)

func (svc *FileService) Create(opts FileCreateOptions, userID string) (*File, error) {
	return svc.fileCreate.create(opts, userID)
}

func (svc *FileService) Find(ids []string, userID string) ([]*File, error) {
	return svc.fileFetch.find(ids, userID)
}

func (svc *FileService) FindByPath(path string, userID string) (*File, error) {
	return svc.fileFetch.findByPath(path, userID)
}

func (svc *FileService) ListByPath(path string, userID string) ([]*File, error) {
	return svc.fileFetch.listByPath(path, userID)
}

func (svc *FileService) FindPath(id string, userID string) ([]*File, error) {
	return svc.fileFetch.findPath(id, userID)
}

func (svc *FileService) GetPathString(files []*File) string {
	return svc.fileFetch.getPathString(files)
}

func (svc *FileService) GetPathStringWithoutWorkspace(files []*File) string {
	return svc.fileFetch.getPathStringWithoutWorkspace(files)
}

func (svc *FileService) Probe(id string, opts FileListOptions, userID string) (*FileProbe, error) {
	return svc.fileList.probe(id, opts, userID)
}

func (svc *FileService) List(id string, opts FileListOptions, userID string) (*FileList, error) {
	return svc.fileList.list(id, opts, userID)
}

func (svc *FileService) IsValidSortBy(value string) bool {
	return svc.fileSortService.isValidSortBy(value)
}

func (svc *FileService) IsValidSortOrder(value string) bool {
	return svc.fileSortService.isValidSortOrder(value)
}

func (svc *FileService) ComputeSize(id string, userID string) (*int64, error) {
	return svc.fileCompute.computeSize(id, userID)
}

func (svc *FileService) Count(id string, userID string) (*int64, error) {
	return svc.fileCompute.count(id, userID)
}

func (svc *FileService) Copy(sourceID string, targetID string, userID string) (*File, error) {
	return svc.fileCopy.copy(sourceID, targetID, userID)
}

func (svc *FileService) CopyMany(opts FileCopyManyOptions, userID string) (*FileCopyManyResult, error) {
	return svc.fileCopy.copyMany(opts, userID)
}

func (svc *FileService) Delete(id string, userID string) error {
	return svc.fileDelete.delete(id, userID)
}

func (svc *FileService) DeleteMany(opts FileDeleteManyOptions, userID string) (*FileDeleteManyResult, error) {
	return svc.fileDelete.deleteMany(opts, userID)
}

func (svc *FileService) DownloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.fileDownload.downloadOriginalBuffer(id, rangeHeader, buf, userID)
}

func (svc *FileService) DownloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.fileDownload.downloadPreviewBuffer(id, rangeHeader, buf, userID)
}

func (svc *FileService) DownloadThumbnailBuffer(id string, buf *bytes.Buffer, userID string) (model.Snapshot, error) {
	return svc.fileDownload.downloadThumbnailBuffer(id, buf, userID)
}

func (svc *FileService) Move(sourceID string, targetID string, userID string) (*File, error) {
	return svc.fileMove.move(sourceID, targetID, userID)
}

func (svc *FileService) MoveMany(opts FileMoveManyOptions, userID string) (*FileMoveManyResult, error) {
	return svc.fileMove.moveMany(opts, userID)
}

func (svc *FileService) PatchName(id string, name string, userID string) (*File, error) {
	return svc.filePatch.patchName(id, name, userID)
}

func (svc *FileService) GrantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	return svc.filePermission.grantUserPermission(ids, assigneeID, permission, userID)
}

func (svc *FileService) RevokeUserPermission(ids []string, assigneeID string, userID string) error {
	return svc.filePermission.revokeUserPermission(ids, assigneeID, userID)
}

func (svc *FileService) GrantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	return svc.filePermission.grantGroupPermission(ids, groupID, permission, userID)
}

func (svc *FileService) RevokeGroupPermission(ids []string, groupID string, userID string) error {
	return svc.filePermission.revokeGroupPermission(ids, groupID, userID)
}

func (svc *FileService) FindUserPermissions(id string, userID string) ([]*UserPermission, error) {
	return svc.filePermission.findUserPermissions(id, userID)
}

func (svc *FileService) FindGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
	return svc.filePermission.findGroupPermissions(id, userID)
}

func (svc *FileService) Reprocess(id string, userID string) (*FileReprocessResult, error) {
	return svc.fileReprocess.reprocess(id, userID)
}

func (svc *FileService) Store(id string, opts FileStoreOptions, userID string) (*File, error) {
	return svc.fileStore.store(id, opts, userID)
}

type fileCreate struct {
	fileRepo    *repo.FileRepo
	fileSearch  *search.FileSearch
	fileCache   *cache.FileCache
	fileGuard   *guard.FileGuard
	fileMapper  *fileMapper
	fileCoreSvc *fileCoreService
}

func newFileCreate() *fileCreate {
	return &fileCreate{
		fileRepo:    repo.NewFileRepo(),
		fileSearch:  search.NewFileSearch(),
		fileCache:   cache.NewFileCache(),
		fileGuard:   guard.NewFileGuard(),
		fileMapper:  newFileMapper(),
		fileCoreSvc: newFileCoreService(),
	}
}

type FileCreateOptions struct {
	WorkspaceID string `json:"workspaceId" validate:"required"`
	Name        string `json:"name"        validate:"required,max=255"`
	Type        string `json:"type"        validate:"required,oneof=file folder"`
	ParentID    string `json:"parentId"    validate:"required"`
}

func (svc *fileCreate) create(opts FileCreateOptions, userID string) (*File, error) {
	path := helper.PathFromFilename(opts.Name)
	parentID := opts.ParentID
	if len(path) > 1 {
		newParentID, err := svc.createDirectoriesForPath(path, parentID, opts.WorkspaceID, userID)
		if err != nil {
			return nil, err
		}
		parentID = *newParentID
	}
	return svc.performCreate(FileCreateOptions{
		WorkspaceID: opts.WorkspaceID,
		Name:        helper.FilenameFromPath(path),
		Type:        opts.Type,
		ParentID:    parentID,
	}, len(path) == 1, userID)
}

func (svc *fileCreate) createDirectoriesForPath(path []string, parentID string, workspaceID string, userID string) (*string, error) {
	for _, component := range path[:len(path)-1] {
		existing, err := svc.fileCoreSvc.getChildWithName(parentID, component)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			parentID = existing.GetID()
		} else {
			folder, err := svc.performCreate(FileCreateOptions{
				Name:        component,
				Type:        model.FileTypeFolder,
				ParentID:    parentID,
				WorkspaceID: workspaceID,
			}, false, userID)
			if err != nil {
				return nil, err
			}
			parentID = folder.ID
		}
	}
	return &parentID, nil
}

func (svc *fileCreate) performCreate(opts FileCreateOptions, failOnDuplicateName bool, userID string) (*File, error) {
	if len(opts.ParentID) > 0 {
		if err := svc.validateParent(opts.ParentID, userID); err != nil {
			return nil, err
		}
		existing, err := svc.fileCoreSvc.getChildWithName(opts.ParentID, opts.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			if failOnDuplicateName {
				return nil, errorpkg.NewFileWithSimilarNameExistsError()
			} else {
				res, err := svc.fileMapper.mapOne(existing, userID)
				if err != nil {
					return nil, err
				}
				return res, nil
			}
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

func (svc *fileCreate) validateParent(id string, userID string) error {
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

type fileFetch struct {
	fileCache      *cache.FileCache
	fileRepo       *repo.FileRepo
	fileSearch     *search.FileSearch
	fileGuard      *guard.FileGuard
	fileMapper     *fileMapper
	fileCoreSvc    *fileCoreService
	fileIdent      *infra.FileIdentifier
	userRepo       *repo.UserRepo
	workspaceRepo  *repo.WorkspaceRepo
	workspaceSvc   *WorkspaceService
	workspaceGuard *guard.WorkspaceGuard
}

func newFileFetch() *fileFetch {
	return &fileFetch{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileSearch:     search.NewFileSearch(),
		fileGuard:      guard.NewFileGuard(),
		fileMapper:     newFileMapper(),
		fileCoreSvc:    newFileCoreService(),
		fileIdent:      infra.NewFileIdentifier(),
		userRepo:       repo.NewUserRepo(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceSvc:   NewWorkspaceService(),
		workspaceGuard: guard.NewWorkspaceGuard(),
	}
}

func (svc *fileFetch) find(ids []string, userID string) ([]*File, error) {
	var res []*File
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			continue
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		mapped, err := svc.fileMapper.mapOne(file, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, mapped)
	}
	return res, nil
}

func (svc *fileFetch) findByPath(path string, userID string) (*File, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	if err := svc.validatePath(path); err != nil {
		return nil, err
	}
	if path == "/" {
		return svc.getUserAsFile(user), nil
	}
	components, err := svc.getComponentsFromPath(path)
	if err != nil {
		return nil, err
	}
	workspace, err := svc.workspaceSvc.Find(svc.workspaceIDFromSlug(components[0]), userID)
	if err != nil {
		return nil, err
	}
	if len(components) == 1 {
		return svc.getWorkspaceAsFile(workspace), nil
	}
	file, err := svc.getFileFromComponents(components, userID)
	if err != nil {
		return nil, err
	}
	if err := svc.validatePathForFileType(path, file.GetType()); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileFetch) listByPath(path string, userID string) ([]*File, error) {
	if err := svc.validatePath(path); err != nil {
		return nil, err
	}
	if path == "/" {
		return svc.getWorkspacesAsFiles(userID)
	}
	components, err := svc.getComponentsFromPath(path)
	if err != nil {
		return nil, err
	}
	file, err := svc.getFileFromComponents(components, userID)
	if err != nil {
		return nil, err
	}
	if err := svc.validatePathForFileType(path, file.GetType()); err != nil {
		return nil, err
	}
	if file.GetType() == model.FileTypeFolder {
		children, err := svc.getAuthorizedChildren(file.GetID(), userID)
		if err != nil {
			return nil, err
		}
		res, err := svc.fileMapper.mapMany(children, userID)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if file.GetType() == model.FileTypeFile {
		res, err := svc.find([]string{file.GetID()}, userID)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else {
		// This should never happen
		return nil, errorpkg.NewFileTypeIsInvalid(file.GetType())
	}
}

func (svc *fileFetch) findPath(id string, userID string) ([]*File, error) {
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

func (svc *fileFetch) getPathString(files []*File) string {
	return strings.Join(svc.getPathStrings(files), "/")
}

func (svc *fileFetch) getPathStringWithoutWorkspace(files []*File) string {
	return strings.Join(svc.getPathStrings(files)[1:], "/")
}

func (svc *fileFetch) getPathStrings(files []*File) []string {
	var components []string
	for _, f := range files {
		components = append(components, f.Name)
	}
	return components
}

func (svc *fileFetch) getWorkspacesAsFiles(userID string) ([]*File, error) {
	workspaces, err := svc.workspaceSvc.findAll(userID)
	if err != nil {
		return nil, err
	}
	res := make([]*File, 0)
	for _, w := range workspaces {
		res = append(res, svc.getWorkspaceAsFile(w))
	}
	return res, nil
}

func (svc *fileFetch) getWorkspaceAsFile(workspace *Workspace) *File {
	return &File{
		ID:          workspace.RootID,
		WorkspaceID: workspace.ID,
		Name:        svc.slugFromWorkspace(workspace.ID, workspace.Name),
		Type:        model.FileTypeFolder,
		Permission:  workspace.Permission,
		CreateTime:  workspace.CreateTime,
		UpdateTime:  workspace.UpdateTime,
	}
}

func (svc *fileFetch) getUserAsFile(user model.User) *File {
	return &File{
		ID:          user.GetID(),
		WorkspaceID: "",
		Name:        "/",
		Type:        model.FileTypeFolder,
		Permission:  "owner",
		CreateTime:  user.GetCreateTime(),
		UpdateTime:  nil,
	}
}

func (svc *fileFetch) getFileFromComponents(components []string, userID string) (model.File, error) {
	workspace, err := svc.workspaceRepo.Find(svc.workspaceIDFromSlug(components[0]))
	if err != nil {
		return nil, err
	}
	res, err := svc.fileCache.Get(workspace.GetRootID())
	if err != nil {
		return nil, err
	}
	for _, component := range components[1:] {
		file, err := svc.findComponentInFolder(component, res.GetID(), userID)
		if err != nil {
			return nil, err
		}
		res = file
		if file.GetType() == model.FileTypeFolder {
			continue
		} else if file.GetType() == model.FileTypeFile {
			break
		}
	}
	return res, err
}

func (svc *fileFetch) findComponentInFolder(component string, id string, userID string) (model.File, error) {
	children, err := svc.getAuthorizedChildren(id, userID)
	if err != nil {
		return nil, err
	}
	var filtered []model.File
	for _, f := range children {
		if f.GetName() == component {
			filtered = append(filtered, f)
		}
	}
	if len(filtered) > 0 {
		return filtered[0], nil
	} else {
		return nil, errorpkg.NewFileNotFoundError(fmt.Errorf("item not found '%s'", component))
	}
}

func (svc *fileFetch) getAuthorizedChildren(id string, userID string) ([]model.File, error) {
	childrenIDs, err := svc.fileRepo.FindChildrenIDs(id)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.fileCoreSvc.authorizeIDs(userID, childrenIDs, model.PermissionViewer)
	if err != nil {
		return nil, err
	}
	return authorized, nil
}

func (svc *fileFetch) getComponentsFromPath(path string) ([]string, error) {
	components := make([]string, 0)
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			components = append(components, v)
		}
	}
	if len(components) == 0 || components[0] == "" {
		return nil, errorpkg.NewInvalidPathError(fmt.Errorf("invalid path '%s'", path))
	}
	return components, nil
}

func (svc *fileFetch) slugFromWorkspace(id string, name string) string {
	return fmt.Sprintf("%s-%s", slug.Make(name), id)
}

func (svc *fileFetch) workspaceIDFromSlug(slug string) string {
	parts := strings.Split(slug, "-")
	return parts[len(parts)-1]
}

func (svc *fileFetch) validatePath(path string) error {
	if !strings.HasPrefix(path, "/") {
		return errorpkg.NewFilePathMissingLeadingSlash()
	}
	return nil
}

func (svc *fileFetch) validatePathForFileType(path string, fileType string) error {
	if fileType == model.FileTypeFile && strings.HasSuffix(path, "/") {
		return errorpkg.NewFilePathOfTypeFileHasTrailingSlash()
	}
	return nil
}

type fileList struct {
	fileCache      *cache.FileCache
	fileRepo       *repo.FileRepo
	fileSearch     *search.FileSearch
	fileGuard      *guard.FileGuard
	fileCoreSvc    *fileCoreService
	fileFilterSvc  *fileFilterService
	fileSortSvc    *fileSortService
	fileMapper     *fileMapper
	workspaceRepo  *repo.WorkspaceRepo
	workspaceGuard *guard.WorkspaceGuard
}

func newFileList() *fileList {
	return &fileList{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileSearch:     search.NewFileSearch(),
		fileGuard:      guard.NewFileGuard(),
		fileCoreSvc:    newFileCoreService(),
		fileFilterSvc:  newFileFilterService(),
		fileSortSvc:    newFileSortService(),
		fileMapper:     newFileMapper(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceGuard: guard.NewWorkspaceGuard(),
	}
}

type FileQuery struct {
	Text             *string `json:"text"                       validate:"required"`
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

type FileProbe struct {
	TotalPages    uint64 `json:"totalPages"`
	TotalElements uint64 `json:"totalElements"`
}

func (svc *fileList) probe(id string, opts FileListOptions, userID string) (*FileProbe, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFolder {
		return nil, errorpkg.NewFileIsNotAFolderError(file)
	}
	totalElements, err := svc.fileRepo.CountChildren(id)
	if err != nil {
		return nil, err
	}
	return &FileProbe{
		TotalElements: uint64(totalElements),                               // #nosec G115 integer overflow conversion
		TotalPages:    (uint64(totalElements) + opts.Size - 1) / opts.Size, // #nosec G115 integer overflow conversion
	}, nil
}

func (svc *fileList) list(id string, opts FileListOptions, userID string) (*FileList, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFolder {
		return nil, errorpkg.NewFileIsNotAFolderError(file)
	}
	workspace, err := svc.workspaceRepo.Find(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	if err := svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	var data []model.File
	if opts.Query != nil && opts.Query.Text != nil {
		data, err = svc.search(opts.Query, workspace)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = svc.getChildren(id)
		if err != nil {
			return nil, err
		}
	}
	return svc.createList(data, file, opts, userID)
}

func (svc *fileList) search(query *FileQuery, workspace model.Workspace) ([]model.File, error) {
	var res []model.File
	filter := fmt.Sprintf("workspaceId=\"%s\"", workspace.GetID())
	if query.Type != nil {
		filter += fmt.Sprintf(" AND type=\"%s\"", *query.Type)
	}
	hits, err := svc.fileSearch.Query(*query.Text, infra.QueryOptions{Filter: filter})
	if err != nil {
		return nil, err
	}
	for _, hit := range hits {
		var file model.File
		file, err := svc.fileCache.Get(hit.GetID())
		if err != nil {
			var e *errorpkg.ErrorResponse
			// We don't want to break if the search engine contains files that shouldn't be there
			if errors.As(err, &e) && e.Code == errorpkg.NewFileNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		res = append(res, file)
	}
	return res, nil
}

func (svc *fileList) getChildren(id string) ([]model.File, error) {
	var res []model.File
	ids, err := svc.fileRepo.FindChildrenIDs(id)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		var file model.File
		file, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		res = append(res, file)
	}
	return res, nil
}

func (svc *fileList) createList(data []model.File, parent model.File, opts FileListOptions, userID string) (*FileList, error) {
	var filtered []model.File
	var err error
	if opts.Query != nil {
		filtered, err = svc.fileFilterSvc.filterWithQuery(data, *opts.Query, parent)
		if err != nil {
			return nil, err
		}
	} else {
		filtered = data
	}
	authorized, err := svc.fileCoreSvc.authorize(userID, filtered, model.PermissionViewer)
	if err != nil {
		return nil, err
	}
	sorted := svc.fileSortSvc.sort(authorized, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
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

func (svc *fileList) paginate(data []model.File, page, size uint64) (pageData []model.File, totalElements uint64, totalPages uint64) {
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

type fileCompute struct {
	fileCache *cache.FileCache
	fileRepo  *repo.FileRepo
	fileGuard *guard.FileGuard
}

func newFileCompute() *fileCompute {
	return &fileCompute{
		fileCache: cache.NewFileCache(),
		fileRepo:  repo.NewFileRepo(),
		fileGuard: guard.NewFileGuard(),
	}
}

func (svc *fileCompute) computeSize(id string, userID string) (*int64, error) {
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

func (svc *fileCompute) count(id string, userID string) (*int64, error) {
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

type fileCopy struct {
	fileRepo     *repo.FileRepo
	fileSearch   *search.FileSearch
	fileCache    *cache.FileCache
	fileGuard    *guard.FileGuard
	fileMapper   *fileMapper
	fileCoreSvc  *fileCoreService
	taskSvc      *TaskService
	snapshotRepo *repo.SnapshotRepo
}

func newFileCopy() *fileCopy {
	return &fileCopy{
		fileRepo:     repo.NewFileRepo(),
		fileSearch:   search.NewFileSearch(),
		fileCache:    cache.NewFileCache(),
		fileGuard:    guard.NewFileGuard(),
		fileMapper:   newFileMapper(),
		fileCoreSvc:  newFileCoreService(),
		taskSvc:      NewTaskService(),
		snapshotRepo: repo.NewSnapshotRepo(),
	}
}

func (svc *fileCopy) copy(sourceID string, targetID string, userID string) (*File, error) {
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return nil, err
	}
	source, err := svc.fileCache.Get(sourceID)
	if err != nil {
		return nil, err
	}
	task, err := svc.createTask(source, userID)
	if err != nil {
		return nil, err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(source, target, userID); err != nil {
		return nil, err
	}
	return svc.performCopy(source, target, userID)
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

func (svc *fileCopy) copyMany(opts FileCopyManyOptions, userID string) (*FileCopyManyResult, error) {
	res := &FileCopyManyResult{
		New:       make([]string, 0),
		Succeeded: make([]string, 0),
		Failed:    make([]string, 0),
	}
	for _, id := range opts.SourceIDs {
		file, err := svc.copy(id, opts.TargetID, userID)
		if err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.New = append(res.New, file.ID)
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

func (svc *fileCopy) performCopy(source model.File, target model.File, userID string) (*File, error) {
	tree, err := svc.getTree(source)
	if err != nil {
		return nil, err
	}
	cloneResult, err := svc.cloneTree(source, target, tree, userID)
	if err != nil {
		return nil, err
	}
	if err := svc.persist(cloneResult.Clones, cloneResult.Permissions); err != nil {
		return nil, err
	}
	if err := svc.attachSnapshots(cloneResult.Clones, tree); err != nil {
		return nil, err
	}
	svc.cache(cloneResult.Clones, userID)
	go svc.index(cloneResult.Clones)
	if err := svc.refreshUpdateTime(target); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(cloneResult.Root, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileCopy) createTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Copying.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileCopy) check(source model.File, target model.File, userID string) error {
	if err := svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
		return err
	}
	if err := svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
		return err
	}
	if source.GetID() == target.GetID() {
		return errorpkg.NewFileCannotBeCopiedIntoItselfError(source)
	}
	if target.GetType() != model.FileTypeFolder {
		return errorpkg.NewFileIsNotAFolderError(target)
	}
	isGrandChild, err := svc.fileRepo.IsGrandChildOf(target.GetID(), source.GetID())
	if err != nil {
		return err
	}
	if isGrandChild {
		return errorpkg.NewFileCannotBeCopiedIntoOwnSubtreeError(source)
	}
	return nil
}

func (svc *fileCopy) getTree(source model.File) ([]model.File, error) {
	var ids []string
	ids, err := svc.fileRepo.FindTreeIDs(source.GetID())
	if err != nil {
		return nil, err
	}
	var tree []model.File
	for _, id := range ids {
		sourceFile, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		tree = append(tree, sourceFile)
	}
	return tree, nil
}

type cloneTreeResult struct {
	Root        model.File
	Clones      []model.File
	Permissions []model.UserPermission
}

func (svc *fileCopy) cloneTree(source model.File, target model.File, tree []model.File, userID string) (*cloneTreeResult, error) {
	var rootIndex int
	ids := make(map[string]string)
	var clones []model.File
	var permissions []model.UserPermission
	for index, leaf := range tree {
		clone := svc.newClone(leaf)
		if leaf.GetID() == source.GetID() {
			rootIndex = index
		}
		ids[leaf.GetID()] = clone.GetID()
		clones = append(clones, clone)
		permissions = append(permissions, svc.newUserPermission(clone, userID))
	}
	root := clones[rootIndex]
	for index, clone := range clones {
		id := ids[*clone.GetParentID()]
		clones[index].SetParentID(&id)
	}
	root.SetParentID(helper.ToPtr(target.GetID()))
	existing, err := svc.fileCoreSvc.getChildWithName(target.GetID(), root.GetName())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		root.SetName(helper.UniqueFilename(root.GetName()))
	}
	return &cloneTreeResult{
		Root:        root,
		Clones:      clones,
		Permissions: permissions,
	}, nil
}

func (svc *fileCopy) newClone(file model.File) model.File {
	f := repo.NewFileModel()
	f.SetID(helper.NewID())
	f.SetParentID(file.GetParentID())
	f.SetWorkspaceID(file.GetWorkspaceID())
	f.SetSnapshotID(file.GetSnapshotID())
	f.SetType(file.GetType())
	f.SetName(file.GetName())
	f.SetCreateTime(helper.NewTimeString())
	return f
}

func (svc *fileCopy) newUserPermission(file model.File, userID string) model.UserPermission {
	p := repo.NewUserPermissionModel()
	p.SetID(helper.NewID())
	p.SetUserID(userID)
	p.SetResourceID(file.GetID())
	p.SetPermission(model.PermissionOwner)
	p.SetCreateTime(helper.NewTimeString())
	return p
}

func (svc *fileCopy) persist(clones []model.File, permissions []model.UserPermission) error {
	const BulkInsertChunkSize = 1000
	if err := svc.fileRepo.BulkInsert(clones, BulkInsertChunkSize); err != nil {
		return err
	}
	if err := svc.fileRepo.BulkInsertPermissions(permissions, BulkInsertChunkSize); err != nil {
		return err
	}
	return nil
}

func (svc *fileCopy) attachSnapshots(clones []model.File, tree []model.File) error {
	const BulkInsertChunkSize = 1000
	var mappings []*repo.SnapshotFileEntity
	for index, clone := range clones {
		leaf := tree[index]
		if leaf.GetSnapshotID() != nil {
			mappings = append(mappings, &repo.SnapshotFileEntity{
				SnapshotID: *leaf.GetSnapshotID(),
				FileID:     clone.GetID(),
			})
		}
	}
	if err := svc.snapshotRepo.BulkMapWithFile(mappings, BulkInsertChunkSize); err != nil {
		return err
	}
	return nil
}

func (svc *fileCopy) cache(clones []model.File, userID string) {
	for _, clone := range clones {
		if _, err := svc.fileCache.RefreshWithExisting(clone, userID); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *fileCopy) index(clones []model.File) {
	if err := svc.fileSearch.Index(clones); err != nil {
		log.GetLogger().Error(err)
	}
}

func (svc *fileCopy) refreshUpdateTime(target model.File) error {
	now := helper.NewTimeString()
	target.SetUpdateTime(&now)
	if err := svc.fileRepo.Save(target); err != nil {
		return err
	}
	if err := svc.fileCoreSvc.sync(target); err != nil {
		return err
	}
	return nil
}

type fileDelete struct {
	fileRepo       *repo.FileRepo
	fileSearch     *search.FileSearch
	fileGuard      *guard.FileGuard
	fileCache      *cache.FileCache
	workspaceCache *cache.WorkspaceCache
	taskSvc        *TaskService
	snapshotRepo   *repo.SnapshotRepo
	snapshotSvc    *SnapshotService
}

func newFileDelete() *fileDelete {
	return &fileDelete{
		fileRepo:       repo.NewFileRepo(),
		fileCache:      cache.NewFileCache(),
		fileSearch:     search.NewFileSearch(),
		fileGuard:      guard.NewFileGuard(),
		workspaceCache: cache.NewWorkspaceCache(),
		taskSvc:        NewTaskService(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotSvc:    NewSnapshotService(),
	}
}

func (svc *fileDelete) delete(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return err
	}
	task, err := svc.createTask(file, userID)
	if err != nil {
		return err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(file); err != nil {
		return err
	}
	return svc.performDelete(file)
}

type FileDeleteManyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileDeleteManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *fileDelete) deleteMany(opts FileDeleteManyOptions, userID string) (*FileDeleteManyResult, error) {
	res := &FileDeleteManyResult{
		Failed:    make([]string, 0),
		Succeeded: make([]string, 0),
	}
	for _, id := range opts.IDs {
		if err := svc.delete(id, userID); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

func (svc *fileDelete) performDelete(file model.File) error {
	if file.GetType() == model.FileTypeFolder {
		return svc.deleteFolder(file.GetID())
	} else if file.GetType() == model.FileTypeFile {
		return svc.deleteFile(file.GetID())
	}
	return nil
}

func (svc *fileDelete) check(file model.File) error {
	if file.GetParentID() == nil {
		workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		return errorpkg.NewCannotDeleteWorkspaceRootError(file, workspace)
	}
	return nil
}

func (svc *fileDelete) createTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Deleting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileDelete) deleteFolder(id string) error {
	treeIDs, err := svc.fileRepo.FindTreeIDs(id)
	if err != nil {
		return err
	}
	// Start by deleting the folder's root to give a quick user feedback
	if err := svc.fileCache.Delete(id); err != nil {
		return err
	}
	if err := svc.fileRepo.Delete(id); err != nil {
		return err
	}
	go func(treeIDs []string) {
		svc.deleteSnapshots(treeIDs)
		svc.deleteFromCache(treeIDs)
		svc.deleteFromRepo(treeIDs)
		svc.deleteFromSearch(treeIDs)
	}(treeIDs)
	return nil
}

func (svc *fileDelete) deleteFile(id string) error {
	if err := svc.snapshotRepo.DeleteMappingsForTree(id); err != nil {
		log.GetLogger().Error(err)
	}
	if err := svc.snapshotSvc.deleteForFile(id); err != nil {
		log.GetLogger().Error(err)
	}
	if err := svc.fileCache.Delete(id); err != nil {
		log.GetLogger().Error(err)
	}
	if err := svc.fileRepo.Delete(id); err != nil {
		log.GetLogger().Error(err)
	}
	if err := svc.fileSearch.Delete([]string{id}); err != nil {
		log.GetLogger().Error(err)
	}
	return nil
}

func (svc *fileDelete) deleteSnapshots(ids []string) {
	for _, id := range ids {
		if err := svc.snapshotSvc.deleteForFile(id); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *fileDelete) deleteFromRepo(ids []string) {
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
}

func (svc *fileDelete) deleteFromCache(ids []string) {
	for _, id := range ids {
		if err := svc.fileCache.Delete(id); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *fileDelete) deleteFromSearch(ids []string) {
	if err := svc.fileSearch.Delete(ids); err != nil {
		log.GetLogger().Error(err)
	}
}

type fileDownload struct {
	fileCache     *cache.FileCache
	fileGuard     *guard.FileGuard
	snapshotCache *cache.SnapshotCache
	s3            infra.S3Manager
}

func newFileDownload() *fileDownload {
	return &fileDownload{
		fileCache:     cache.NewFileCache(),
		fileGuard:     guard.NewFileGuard(),
		snapshotCache: cache.NewSnapshotCache(),
		s3:            infra.NewS3Manager(),
	}
}

type DownloadResult struct {
	File          model.File
	Snapshot      model.Snapshot
	RangeInterval *infra.RangeInterval
}

func (svc *fileDownload) downloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if err = svc.check(file); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasOriginal() {
		rangeInterval, err := svc.downloadS3Object(snapshot.GetOriginal(), rangeHeader, buf)
		if err != nil {
			return nil, err
		}
		return &DownloadResult{
			File:          file,
			Snapshot:      snapshot,
			RangeInterval: rangeInterval,
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *fileDownload) downloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if err = svc.check(file); err != nil {
		return nil, err
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasPreview() {
		rangeInterval, err := svc.downloadS3Object(snapshot.GetPreview(), rangeHeader, buf)
		if err != nil {
			return nil, err
		}
		return &DownloadResult{
			File:          file,
			Snapshot:      snapshot,
			RangeInterval: rangeInterval,
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *fileDownload) downloadThumbnailBuffer(id string, buf *bytes.Buffer, userID string) (model.Snapshot, error) {
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

func (svc *fileDownload) check(file model.File) error {
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	return nil
}

func (svc *fileDownload) downloadS3Object(s3Object *model.S3Object, rangeHeader string, buf *bytes.Buffer) (*infra.RangeInterval, error) {
	objectInfo, err := svc.s3.StatObject(s3Object.Key, s3Object.Bucket, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	opts := minio.GetObjectOptions{}
	var rangeInterval *infra.RangeInterval
	if rangeHeader != "" {
		rangeInterval = infra.NewRangeInterval(rangeHeader, objectInfo.Size)
		if err := rangeInterval.ApplyToMinIOGetObjectOptions(&opts); err != nil {
			return nil, err
		}
	}
	if _, err := svc.s3.GetObjectWithBuffer(s3Object.Key, s3Object.Bucket, buf, opts); err != nil {
		return nil, err
	}
	return rangeInterval, nil
}

type fileMove struct {
	fileRepo    *repo.FileRepo
	fileSearch  *search.FileSearch
	fileCache   *cache.FileCache
	fileGuard   *guard.FileGuard
	fileMapper  *fileMapper
	fileCoreSvc *fileCoreService
	taskSvc     *TaskService
}

func newFileMove() *fileMove {
	return &fileMove{
		fileRepo:    repo.NewFileRepo(),
		fileSearch:  search.NewFileSearch(),
		fileCache:   cache.NewFileCache(),
		fileGuard:   guard.NewFileGuard(),
		fileMapper:  newFileMapper(),
		fileCoreSvc: newFileCoreService(),
		taskSvc:     NewTaskService(),
	}
}

func (svc *fileMove) move(sourceID string, targetID string, userID string) (*File, error) {
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return nil, err
	}
	source, err := svc.fileCache.Get(sourceID)
	if err != nil {
		return nil, err
	}
	task, err := svc.createTask(source, userID)
	if err != nil {
		return nil, err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(source, target, userID); err != nil {
		return nil, err
	}
	return svc.performMove(source, target, userID)
}

type FileMoveManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileMoveManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *fileMove) moveMany(opts FileMoveManyOptions, userID string) (*FileMoveManyResult, error) {
	res := &FileMoveManyResult{
		Failed:    make([]string, 0),
		Succeeded: make([]string, 0),
	}
	for _, id := range opts.SourceIDs {
		if _, err := svc.move(id, opts.TargetID, userID); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

func (svc *fileMove) performMove(source model.File, target model.File, userID string) (*File, error) {
	if err := svc.fileRepo.MoveSourceIntoTarget(target.GetID(), source.GetID()); err != nil {
		return nil, err
	}
	var err error
	source, err = svc.fileRepo.Find(source.GetID())
	if err != nil {
		return nil, err
	}
	if err := svc.refreshUpdateAndCreateTime(source, target); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(source, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileMove) createTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Moving.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileMove) check(source model.File, target model.File, userID string) error {
	if source.GetParentID() != nil {
		existing, err := svc.fileCoreSvc.getChildWithName(target.GetID(), source.GetName())
		if err != nil {
			return err
		}
		if existing != nil {
			return errorpkg.NewFileWithSimilarNameExistsError()
		}
	}
	if err := svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
		return err
	}
	if err := svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
		return err
	}
	if source.GetParentID() != nil && *source.GetParentID() == target.GetID() {
		return errorpkg.NewFileAlreadyChildOfDestinationError(source, target)
	}
	if target.GetID() == source.GetID() {
		return errorpkg.NewFileCannotBeMovedIntoItselfError(source)
	}
	if target.GetType() != model.FileTypeFolder {
		return errorpkg.NewFileIsNotAFolderError(target)
	}
	isGrandChild, err := svc.fileRepo.IsGrandChildOf(target.GetID(), source.GetID())
	if err != nil {
		return err
	}
	if isGrandChild {
		return errorpkg.NewTargetIsGrandChildOfSourceError(source)
	}
	return nil
}

func (svc *fileMove) refreshUpdateAndCreateTime(source model.File, target model.File) error {
	now := helper.NewTimeString()
	source.SetUpdateTime(&now)
	if err := svc.fileRepo.Save(source); err != nil {
		return err
	}
	if err := svc.fileCoreSvc.sync(source); err != nil {
		return err
	}
	target.SetUpdateTime(&now)
	if err := svc.fileRepo.Save(target); err != nil {
		return err
	}
	if err := svc.fileCoreSvc.sync(target); err != nil {
		return err
	}
	return nil
}

type filePatch struct {
	fileCache   *cache.FileCache
	fileRepo    *repo.FileRepo
	fileGuard   *guard.FileGuard
	fileCoreSvc *fileCoreService
	fileMapper  *fileMapper
}

func newFilePatch() *filePatch {
	return &filePatch{
		fileCache:   cache.NewFileCache(),
		fileRepo:    repo.NewFileRepo(),
		fileGuard:   guard.NewFileGuard(),
		fileCoreSvc: newFileCoreService(),
		fileMapper:  newFileMapper(),
	}
}

func (svc *filePatch) patchName(id string, name string, userID string) (*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if file.GetParentID() != nil {
		existing, err := svc.fileCoreSvc.getChildWithName(*file.GetParentID(), name)
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
	if err := svc.fileCoreSvc.sync(file); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type filePermission struct {
	fileCache      *cache.FileCache
	fileRepo       *repo.FileRepo
	fileGuard      *guard.FileGuard
	fileCoreSvc    *fileCoreService
	userRepo       *repo.UserRepo
	userMapper     *userMapper
	workspaceRepo  *repo.WorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	groupCache     *cache.GroupCache
	groupGuard     *guard.GroupGuard
	groupMapper    *groupMapper
	permissionRepo *repo.PermissionRepo
}

func newFilePermission() *filePermission {
	return &filePermission{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileGuard:      guard.NewFileGuard(),
		fileCoreSvc:    newFileCoreService(),
		userRepo:       repo.NewUserRepo(),
		userMapper:     newUserMapper(),
		workspaceRepo:  repo.NewWorkspaceRepo(),
		workspaceCache: cache.NewWorkspaceCache(),
		groupCache:     cache.NewGroupCache(),
		groupGuard:     guard.NewGroupGuard(),
		groupMapper:    newGroupMapper(),
		permissionRepo: repo.NewPermissionRepo(),
	}
}

func (svc *filePermission) grantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	for _, id := range ids {
		if err := svc.grantOneUserPermission(id, assigneeID, permission, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) grantOneUserPermission(id string, assigneeID string, permission string, userID string) error {
	file, err := svc.authorizeUserPermission(id, assigneeID, userID)
	if err != nil {
		return err
	}
	if err = svc.fileRepo.GrantUserPermission(id, assigneeID, permission); err != nil {
		return err
	}
	if _, err := svc.workspaceCache.Refresh(file.GetWorkspaceID()); err != nil {
		return err
	}
	if err := svc.refreshPathAndTree(id); err != nil {
		return nil
	}
	return nil
}

func (svc *filePermission) revokeUserPermission(ids []string, assigneeID string, userID string) error {
	for _, id := range ids {
		if err := svc.revokeOneUserPermission(id, assigneeID, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) revokeOneUserPermission(id string, assigneeID string, userID string) error {
	if _, err := svc.authorizeUserPermission(id, assigneeID, userID); err != nil {
		return err
	}
	tree, err := svc.fileRepo.FindTree(id)
	if err != nil {
		return err
	}
	if err := svc.fileRepo.RevokeUserPermission(tree, assigneeID); err != nil {
		return err
	}
	for _, leaf := range tree {
		if _, err := svc.fileCache.Refresh(leaf.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) grantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	for _, id := range ids {
		if err := svc.grantOneGroupPermission(id, groupID, permission, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) grantOneGroupPermission(id string, groupID string, permission string, userID string) error {
	file, _, err := svc.authorizeGroupPermission(id, groupID, userID)
	if err != nil {
		return err
	}
	if err := svc.fileRepo.GrantGroupPermission(id, groupID, permission); err != nil {
		return err
	}
	if _, err := svc.workspaceCache.Refresh(file.GetWorkspaceID()); err != nil {
		return err
	}
	if err := svc.refreshPathAndTree(id); err != nil {
		return nil
	}
	return nil
}

func (svc *filePermission) revokeGroupPermission(ids []string, groupID string, userID string) error {
	for _, id := range ids {
		if err := svc.revokeOneGroupPermission(id, groupID, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) revokeOneGroupPermission(id string, groupID string, userID string) error {
	if _, _, err := svc.authorizeGroupPermission(id, groupID, userID); err != nil {
		return err
	}
	tree, err := svc.fileRepo.FindTree(id)
	if err != nil {
		return err
	}
	if err := svc.fileRepo.RevokeGroupPermission(tree, groupID); err != nil {
		return err
	}
	for _, leaf := range tree {
		if _, err := svc.fileCache.Refresh(leaf.GetID()); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) authorizeUserPermission(id string, assigneeID string, userID string) (model.File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, err
	}
	if _, err := svc.userRepo.Find(assigneeID); err != nil {
		return nil, err
	}
	return file, nil
}

func (svc *filePermission) authorizeGroupPermission(id string, groupID string, userID string) (model.File, model.Group, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, err
	}
	if err := svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return nil, nil, err
	}
	group, err := svc.groupCache.Get(groupID)
	if err != nil {
		return nil, nil, err
	}
	if err := svc.groupGuard.Authorize(userID, group, model.PermissionViewer); err != nil {
		return nil, nil, err
	}
	return file, group, nil
}

func (svc *filePermission) refreshPathAndTree(id string) error {
	path, err := svc.fileRepo.FindPath(id)
	if err != nil {
		return err
	}
	for _, f := range path {
		if _, err := svc.fileCache.Refresh(f.GetID()); err != nil {
			return err
		}
	}
	tree, err := svc.fileRepo.FindTree(id)
	if err != nil {
		return err
	}
	for _, leaf := range tree {
		if _, err := svc.fileCache.Refresh(leaf.GetID()); err != nil {
			return err
		}
	}
	return nil
}

type UserPermission struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	Permission string `json:"permission"`
}

func (svc *filePermission) findUserPermissions(id string, userID string) ([]*UserPermission, error) {
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

func (svc *filePermission) findGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
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

type fileReprocess struct {
	fileCache      *cache.FileCache
	fileRepo       *repo.FileRepo
	fileGuard      *guard.FileGuard
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	taskCache      *cache.TaskCache
	taskSvc        *TaskService
	fileIdent      *infra.FileIdentifier
	pipelineClient conversion_client.PipelineClient
}

func newFileReprocess() *fileReprocess {
	return &fileReprocess{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileGuard:      guard.NewFileGuard(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotSvc:    NewSnapshotService(),
		taskCache:      cache.NewTaskCache(),
		taskSvc:        NewTaskService(),
		fileIdent:      infra.NewFileIdentifier(),
		pipelineClient: conversion_client.NewPipelineClient(),
	}
}

type FileReprocessResult struct {
	Accepted []string `json:"accepted"`
	Rejected []string `json:"rejected"`
}

func (r *FileReprocessResult) AppendAccepted(id string) {
	if !slices.Contains(r.Accepted, id) {
		r.Accepted = append(r.Accepted, id)
	}
}

func (r *FileReprocessResult) AppendRejected(id string) {
	if !slices.Contains(r.Rejected, id) {
		r.Rejected = append(r.Rejected, id)
	}
}

func (svc *fileReprocess) reprocess(id string, userID string) (*FileReprocessResult, error) {
	resp := &FileReprocessResult{
		// We intend to send an empty array to the caller, better than nil
		Accepted: []string{},
		Rejected: []string{},
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	tree, err := svc.getTree(file, userID)
	if err != nil {
		return nil, err
	}
	for _, leaf := range tree {
		snapshot, err := svc.snapshotCache.Get(*leaf.GetSnapshotID())
		if err != nil {
			return nil, err
		}
		if *snapshot.GetOriginal().Size <= helper.MegabyteToByte(svc.fileIdent.GetProcessingLimitMB(leaf.GetName())) {
			if svc.performReprocess(leaf, userID) {
				resp.AppendAccepted(leaf.GetID())
			} else {
				resp.AppendRejected(leaf.GetID())
			}
		}
	}
	return resp, nil
}

func (svc *fileReprocess) performReprocess(leaf model.File, userID string) bool {
	if leaf.GetType() != model.FileTypeFile {
		return false
	}
	if err := svc.fileGuard.Authorize(userID, leaf, model.PermissionEditor); err != nil {
		log.GetLogger().Error(err)
		return false
	}
	snapshot, err := svc.snapshotCache.Get(*leaf.GetSnapshotID())
	if err != nil {
		log.GetLogger().Error(err)
		return false
	}
	if !svc.check(leaf, snapshot) {
		return false
	}
	if err := svc.runPipeline(leaf, snapshot, userID); err != nil {
		log.GetLogger().Error(err)
		return false
	}
	return true
}

func (svc *fileReprocess) check(file model.File, snapshot model.Snapshot) bool {
	if file.GetSnapshotID() == nil {
		// We don't reprocess if there is no active snapshot
		return false
	}
	if snapshot.GetTaskID() != nil {
		task, err := svc.taskCache.Get(*snapshot.GetTaskID())
		if err != nil {
			log.GetLogger().Error(err)
			return false
		}
		if task.GetStatus() == model.TaskStatusWaiting || task.GetStatus() == model.TaskStatusRunning {
			// We don't reprocess if there is a pending task
			return false
		}
	}
	if !snapshot.HasOriginal() {
		// We don't reprocess without an original on the active snapshot
		return false
	}
	return true
}

func (svc *fileReprocess) getTree(file model.File, userID string) ([]model.File, error) {
	var tree []model.File
	var err error
	if file.GetType() == model.FileTypeFolder {
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		tree, err = svc.fileRepo.FindTree(file.GetID())
		if err != nil {
			return nil, err
		}
	} else if file.GetType() == model.FileTypeFile {
		tree = append(tree, file)
	}
	return tree, nil
}

func (svc *fileReprocess) createTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
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
	return res, nil
}

func (svc *fileReprocess) runPipeline(file model.File, snapshot model.Snapshot, userID string) error {
	task, err := svc.createTask(file, userID)
	if err != nil {
		return err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return err
	}
	if err := svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        snapshot.GetOriginal().Key,
	}); err != nil {
		return err
	}
	return nil
}

type fileStore struct {
	fileCache      *cache.FileCache
	fileCoreSvc    *fileCoreService
	fileMapper     *fileMapper
	workspaceCache *cache.WorkspaceCache
	snapshotRepo   *repo.SnapshotRepo
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	taskSvc        *TaskService
	fileIdent      *infra.FileIdentifier
	s3             infra.S3Manager
	pipelineClient conversion_client.PipelineClient
}

func newFileStore() *fileStore {
	return &fileStore{
		fileCache:      cache.NewFileCache(),
		fileCoreSvc:    newFileCoreService(),
		fileMapper:     newFileMapper(),
		workspaceCache: cache.NewWorkspaceCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotSvc:    NewSnapshotService(),
		taskSvc:        NewTaskService(),
		fileIdent:      infra.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		pipelineClient: conversion_client.NewPipelineClient(),
	}
}

type FileStoreOptions struct {
	S3Reference *model.S3Reference
	Path        *string
}

func (svc *fileStore) store(id string, opts FileStoreOptions, userID string) (*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	props, err := svc.getProperties(file, opts)
	if err != nil {
		return nil, err
	}
	if opts.S3Reference == nil {
		if err := svc.performStore(props); err != nil {
			return nil, err
		}
	}
	snapshot, err := svc.createSnapshot(file, props)
	if err != nil {
		return nil, err
	}
	if err := svc.assignSnapshotToFile(file, snapshot); err != nil {
		return nil, err
	}
	if !props.ExceedsProcessingLimit {
		if err := svc.runPipeline(file, snapshot, props, userID); err != nil {
			return nil, err
		}
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type fileStoreProperties struct {
	SnapshotID             string
	Size                   int64
	Path                   string
	Original               model.S3Object
	Bucket                 string
	ContentType            string
	ExceedsProcessingLimit bool
}

func (svc *fileStore) getProperties(file model.File, opts FileStoreOptions) (fileStoreProperties, error) {
	props := fileStoreProperties{}
	if opts.S3Reference == nil {
		var err error
		props, err = svc.getPropertiesFromPath(file, opts)
		if err != nil {
			return fileStoreProperties{}, err
		}
	} else {
		props = svc.getPropertiesFromS3Reference(opts)
	}
	props.ExceedsProcessingLimit = props.Size > helper.MegabyteToByte(svc.fileIdent.GetProcessingLimitMB(props.Path))
	return props, nil
}

func (svc *fileStore) getPropertiesFromPath(file model.File, opts FileStoreOptions) (fileStoreProperties, error) {
	stat, err := os.Stat(*opts.Path)
	if err != nil {
		return fileStoreProperties{}, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return fileStoreProperties{}, err
	}
	snapshotID := helper.NewID()
	return fileStoreProperties{
		SnapshotID: snapshotID,
		Path:       *opts.Path,
		Size:       stat.Size(),
		Original: model.S3Object{
			Bucket: workspace.GetBucket(),
			Key:    snapshotID + "/original" + strings.ToLower(filepath.Ext(*opts.Path)),
			Size:   helper.ToPtr(stat.Size()),
		},
		Bucket:      workspace.GetBucket(),
		ContentType: infra.DetectMIMEFromPath(*opts.Path),
	}, nil
}

func (svc *fileStore) getPropertiesFromS3Reference(opts FileStoreOptions) fileStoreProperties {
	return fileStoreProperties{
		SnapshotID: opts.S3Reference.SnapshotID,
		Path:       opts.S3Reference.Key,
		Size:       opts.S3Reference.Size,
		Original: model.S3Object{
			Bucket: opts.S3Reference.Bucket,
			Key:    opts.S3Reference.Key,
			Size:   helper.ToPtr(opts.S3Reference.Size),
		},
		Bucket:      opts.S3Reference.Bucket,
		ContentType: opts.S3Reference.ContentType,
	}
}

func (svc *fileStore) performStore(props fileStoreProperties) error {
	if err := svc.s3.PutFile(props.Original.Key, props.Path, props.ContentType, props.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (svc *fileStore) createSnapshot(file model.File, props fileStoreProperties) (model.Snapshot, error) {
	res := repo.NewSnapshotModel()
	res.SetID(props.SnapshotID)
	if props.ExceedsProcessingLimit {
		res.SetStatus(model.SnapshotStatusReady)
	} else {
		res.SetStatus(model.SnapshotStatusWaiting)
	}
	latestVersion, err := svc.snapshotRepo.FindLatestVersionForFile(file.GetID())
	if err != nil {
		return nil, err
	}
	res.SetVersion(latestVersion + 1)
	res.SetOriginal(&props.Original)
	if err := svc.snapshotSvc.insertAndSync(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileStore) assignSnapshotToFile(file model.File, snapshot model.Snapshot) error {
	file.SetSnapshotID(helper.ToPtr(snapshot.GetID()))
	if err := svc.fileCoreSvc.saveAndSync(file); err != nil {
		return err
	}
	if err := svc.snapshotRepo.MapWithFile(snapshot.GetID(), file.GetID()); err != nil {
		return err
	}
	return nil
}

func (svc *fileStore) createTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
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
	return res, nil
}

func (svc *fileStore) runPipeline(file model.File, snapshot model.Snapshot, props fileStoreProperties, userID string) error {
	task, err := svc.createTask(file, userID)
	if err != nil {
		return err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return err
	}
	if err := svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     props.Original.Bucket,
		Key:        props.Original.Key,
	}); err != nil {
		return err
	}
	return nil
}

type fileCoreService struct {
	fileRepo   *repo.FileRepo
	fileSearch *search.FileSearch
	fileCache  *cache.FileCache
	fileGuard  *guard.FileGuard
}

func newFileCoreService() *fileCoreService {
	return &fileCoreService{
		fileRepo:   repo.NewFileRepo(),
		fileCache:  cache.NewFileCache(),
		fileSearch: search.NewFileSearch(),
		fileGuard:  guard.NewFileGuard(),
	}
}

func (svc *fileCoreService) getChildWithName(id string, name string) (model.File, error) {
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

func (svc *fileCoreService) sync(file model.File) error {
	if err := svc.fileSearch.Update([]model.File{file}); err != nil {
		return err
	}
	if err := svc.fileCache.Set(file); err != nil {
		return err
	}
	return nil
}

func (svc *fileCoreService) saveAndSync(file model.File) error {
	if err := svc.fileRepo.Save(file); err != nil {
		return err
	}
	if err := svc.sync(file); err != nil {
		return err
	}
	return nil
}

func (svc *fileCoreService) authorize(userID string, files []model.File, permission string) ([]model.File, error) {
	var res []model.File
	for _, f := range files {
		if svc.fileGuard.IsAuthorized(userID, f, permission) {
			res = append(res, f)
		}
	}
	return res, nil
}

func (svc *fileCoreService) authorizeIDs(userID string, ids []string, permission string) ([]model.File, error) {
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
		if svc.fileGuard.IsAuthorized(userID, f, permission) {
			res = append(res, f)
		}
	}
	return res, nil
}

type fileFilterService struct {
	fileRepo   *repo.FileRepo
	fileGuard  *guard.FileGuard
	fileMapper *fileMapper
	fileIdent  *infra.FileIdentifier
}

func newFileFilterService() *fileFilterService {
	return &fileFilterService{
		fileRepo:   repo.NewFileRepo(),
		fileGuard:  guard.NewFileGuard(),
		fileMapper: newFileMapper(),
		fileIdent:  infra.NewFileIdentifier(),
	}
}

func (svc *fileFilterService) filterFolders(data []model.File) []model.File {
	folders, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetType() == model.FileTypeFolder
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range folders {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterFiles(data []model.File) []model.File {
	files, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			return v.(model.File).GetType() == model.FileTypeFile
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range files {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterImages(data []model.File, userID string) []model.File {
	images, _ := rxgo.Just(data)().
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
	var res []model.File
	for _, v := range images {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterPDFs(data []model.File, userID string) []model.File {
	pdfs, _ := rxgo.Just(data)().
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
	var res []model.File
	for _, v := range pdfs {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterDocuments(data []model.File, userID string) []model.File {
	documents, _ := rxgo.Just(data)().
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
	var res []model.File
	for _, v := range documents {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterVideos(data []model.File, userID string) []model.File {
	videos, _ := rxgo.Just(data)().
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
	var res []model.File
	for _, v := range videos {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterTexts(data []model.File, userID string) []model.File {
	texts, _ := rxgo.Just(data)().
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
	var res []model.File
	for _, v := range texts {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterOthers(data []model.File, userID string) []model.File {
	others, _ := rxgo.Just(data)().
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
	for _, v := range others {
		res = append(res, v.(model.File))
	}
	return res
}

func (svc *fileFilterService) filterWithQuery(data []model.File, opts FileQuery, parent model.File) ([]model.File, error) {
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
				return helper.StringToTime(v.(model.File).GetCreateTime()).UnixMilli() >= *opts.CreateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.CreateTimeBefore != nil {
				file := v.(model.File)
				return helper.StringToTime(file.GetCreateTime()).UnixMilli() <= *opts.CreateTimeBefore
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeAfter != nil {
				file := v.(model.File)
				return file.GetUpdateTime() != nil && helper.StringToTime(*file.GetUpdateTime()).UnixMilli() >= *opts.UpdateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeBefore != nil {
				file := v.(model.File)
				return file.GetUpdateTime() != nil && helper.StringToTime(*file.GetUpdateTime()).UnixMilli() <= *opts.UpdateTimeBefore
			} else {
				return true
			}
		}).
		ToSlice(0)
	var res []model.File
	for _, v := range filtered {
		res = append(res, v.(model.File))
	}
	return res, nil
}

type fileSortService struct {
	fileMapper    *fileMapper
	fileFilterSvc *fileFilterService
}

func newFileSortService() *fileSortService {
	return &fileSortService{
		fileMapper:    newFileMapper(),
		fileFilterSvc: newFileFilterService(),
	}
}

func (svc *fileSortService) sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
	if sortBy == FileSortByName {
		return svc.sortByName(data, sortOrder)
	} else if sortBy == FileSortBySize {
		return svc.sortBySize(data, sortOrder, userID)
	} else if sortBy == FileSortByDateCreated {
		return svc.sortByDateCreated(data, sortOrder)
	} else if sortBy == FileSortByDateModified {
		return svc.sortByDateModified(data, sortOrder)
	} else if sortBy == FileSortByKind {
		return svc.sortByKind(data, userID)
	}
	return data
}

func (svc *fileSortService) sortByName(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if sortOrder == FileSortOrderDesc {
			return data[i].GetName() > data[j].GetName()
		} else {
			return data[i].GetName() < data[j].GetName()
		}
	})
	return data
}

func (svc *fileSortService) sortBySize(data []model.File, sortOrder string, userID string) []model.File {
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
		if sortOrder == FileSortOrderDesc {
			return sizeA > sizeB
		} else {
			return sizeA < sizeB
		}
	})
	return data
}

func (svc *fileSortService) sortByDateCreated(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		a := helper.StringToTime(data[i].GetCreateTime())
		b := helper.StringToTime(data[j].GetCreateTime())
		if sortOrder == FileSortOrderDesc {
			return a.UnixMilli() > b.UnixMilli()
		} else {
			return a.UnixMilli() < b.UnixMilli()
		}
	})
	return data
}

func (svc *fileSortService) sortByDateModified(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
			a := helper.StringToTime(*data[i].GetUpdateTime())
			b := helper.StringToTime(*data[j].GetUpdateTime())
			if sortOrder == FileSortOrderDesc {
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

func (svc *fileSortService) sortByKind(data []model.File, userID string) []model.File {
	var res []model.File
	folders := svc.fileFilterSvc.filterFolders(data)
	files := svc.fileFilterSvc.filterFiles(data)
	res = append(res, folders...)
	res = append(res, files...)
	res = append(res, svc.fileFilterSvc.filterImages(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterPDFs(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterDocuments(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterVideos(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterTexts(files, userID)...)
	res = append(res, svc.fileFilterSvc.filterOthers(files, userID)...)
	return res
}

func (svc *fileSortService) isValidSortBy(value string) bool {
	return value == "" ||
		value == FileSortByName ||
		value == FileSortByKind ||
		value == FileSortBySize ||
		value == FileSortByDateCreated ||
		value == FileSortByDateModified
}

func (svc *fileSortService) isValidSortOrder(value string) bool {
	return value == "" || value == FileSortOrderAsc || value == FileSortOrderDesc
}

type fileMapper struct {
	groupCache     *cache.GroupCache
	snapshotMapper *snapshotMapper
	snapshotCache  *cache.SnapshotCache
	snapshotRepo   *repo.SnapshotRepo
}

func newFileMapper() *fileMapper {
	return &fileMapper{
		groupCache:     cache.NewGroupCache(),
		snapshotMapper: newSnapshotMapper(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
	}
}

func (mp *fileMapper) mapOne(m model.File, userID string) (*File, error) {
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

func (mp *fileMapper) mapMany(data []model.File, userID string) ([]*File, error) {
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
