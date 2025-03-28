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
	"sort"
	"strings"

	"github.com/gosimple/slug"
	"github.com/minio/minio-go/v7"
	"github.com/reactivex/rxgo/v2"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/guard"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/mapper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
	"github.com/kouprlabs/voltaserve/shared/search"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type FileService struct {
	fileCreate      *fileCreate
	fileStore       *fileStore
	fileDelete      *FileDelete
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
		fileDelete:      NewFileDelete(),
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

func (svc *FileService) Create(opts FileCreateOptions, userID string) (*dto.File, error) {
	return svc.fileCreate.create(opts, userID)
}

func (svc *FileService) Find(ids []string, userID string) ([]*dto.File, error) {
	return svc.fileFetch.find(ids, userID)
}

func (svc *FileService) FindByPath(path string, userID string) (*dto.File, error) {
	return svc.fileFetch.findByPath(path, userID)
}

func (svc *FileService) ListByPath(path string, userID string) ([]*dto.File, error) {
	return svc.fileFetch.listByPath(path, userID)
}

func (svc *FileService) FindPath(id string, userID string) ([]*dto.File, error) {
	return svc.fileFetch.findPath(id, userID)
}

func (svc *FileService) GetPathString(files []*dto.File) string {
	return svc.fileFetch.getPathString(files)
}

func (svc *FileService) GetPathStringWithoutWorkspace(files []*dto.File) string {
	return svc.fileFetch.getPathStringWithoutWorkspace(files)
}

type FileListOptions struct {
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
	Query     *dto.FileQuery
}

func (svc *FileService) Probe(id string, opts FileListOptions, userID string) (*dto.FileProbe, error) {
	return svc.fileList.probe(id, opts, userID)
}

func (svc *FileService) List(id string, opts FileListOptions, userID string) (*dto.FileList, error) {
	return svc.fileList.list(id, opts, userID)
}

func (svc *FileService) IsValidSortBy(value string) bool {
	return svc.fileSortService.isValidSortBy(value)
}

func (svc *FileService) IsValidSortOrder(value string) bool {
	return svc.fileSortService.isValidSortOrder(value)
}

func (svc *FileService) GetSize(id string, userID string) (*int64, error) {
	return svc.fileCompute.getSize(id, userID)
}

func (svc *FileService) GetCount(id string, userID string) (*int64, error) {
	return svc.fileCompute.getCount(id, userID)
}

func (svc *FileService) Copy(sourceID string, targetID string, userID string) (*dto.File, error) {
	return svc.fileCopy.copy(sourceID, targetID, userID)
}

func (svc *FileService) CopyMany(opts dto.FileCopyManyOptions, userID string) (*dto.FileCopyManyResult, error) {
	return svc.fileCopy.copyMany(opts, userID)
}

func (svc *FileService) Delete(id string, userID string) error {
	return svc.fileDelete.delete(id, userID)
}

func (svc *FileService) DeleteMany(opts dto.FileDeleteManyOptions, userID string) (*dto.FileDeleteManyResult, error) {
	return svc.fileDelete.deleteMany(opts, userID)
}

func (svc *FileService) DownloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.fileDownload.downloadOriginalBuffer(id, rangeHeader, buf, userID)
}

func (svc *FileService) DownloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.fileDownload.downloadPreviewBuffer(id, rangeHeader, buf, userID)
}

func (svc *FileService) DownloadTextBuffer(id string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.fileDownload.downloadTextBuffer(id, buf, userID)
}

func (svc *FileService) DownloadOCRBuffer(id string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.fileDownload.downloadOCRBuffer(id, buf, userID)
}

func (svc *FileService) DownloadThumbnailBuffer(id string, buf *bytes.Buffer, userID string) (model.Snapshot, error) {
	return svc.fileDownload.downloadThumbnailBuffer(id, buf, userID)
}

func (svc *FileService) Move(sourceID string, targetID string, userID string) (*dto.File, error) {
	return svc.fileMove.move(sourceID, targetID, userID)
}

func (svc *FileService) MoveMany(opts dto.FileMoveManyOptions, userID string) (*dto.FileMoveManyResult, error) {
	return svc.fileMove.moveMany(opts, userID)
}

func (svc *FileService) PatchName(id string, name string, userID string) (*dto.File, error) {
	return svc.filePatch.patchName(id, name, userID)
}

func (svc *FileService) GrantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	return svc.filePermission.grantUserPermissions(ids, assigneeID, permission, userID)
}

func (svc *FileService) RevokeUserPermission(ids []string, assigneeID string, userID string) error {
	return svc.filePermission.revokeUserPermissions(ids, assigneeID, userID)
}

func (svc *FileService) GrantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	return svc.filePermission.grantGroupPermissions(ids, groupID, permission, userID)
}

func (svc *FileService) RevokeGroupPermission(ids []string, groupID string, userID string) error {
	return svc.filePermission.revokeGroupPermissions(ids, groupID, userID)
}

func (svc *FileService) FindUserPermissions(id string, userID string) ([]*dto.UserPermission, error) {
	return svc.filePermission.findUserPermissions(id, userID)
}

func (svc *FileService) FindGroupPermissions(id string, userID string) ([]*dto.GroupPermission, error) {
	return svc.filePermission.findGroupPermissions(id, userID)
}

func (svc *FileService) Reprocess(id string, userID string) (*dto.FileReprocessResult, error) {
	return svc.fileReprocess.reprocess(id, userID)
}

func (svc *FileService) Store(id string, opts FileStoreOptions, userID string) (*dto.File, error) {
	return svc.fileStore.store(id, opts, userID)
}

type fileCreate struct {
	fileRepo       *repo.FileRepo
	fileSearch     *search.FileSearch
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	fileMapper     *mapper.FileMapper
	fileCoreSvc    *fileCoreService
	workspaceCache *cache.WorkspaceCache
	workspaceGuard *guard.WorkspaceGuard
}

func newFileCreate() *fileCreate {
	return &fileCreate{
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceGuard: guard.NewWorkspaceGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

type FileCreateOptions struct {
	WorkspaceID string `json:"workspaceId" validate:"required"`
	Name        string `json:"name"        validate:"required,max=255"`
	Type        string `json:"type"        validate:"required,oneof=file folder"`
	ParentID    string `json:"parentId"    validate:"required"`
}

func (svc *fileCreate) create(opts FileCreateOptions, userID string) (*dto.File, error) {
	workspace, err := svc.workspaceCache.Get(opts.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if err = svc.workspaceGuard.Authorize(userID, workspace, model.PermissionViewer); err != nil {
		return nil, err
	}
	if len(opts.ParentID) > 0 {
		if err := svc.validateParent(opts.ParentID, userID); err != nil {
			return nil, err
		}
	} else {
		if err := svc.validateParent(workspace.GetRootID(), userID); err != nil {
			return nil, err
		}
	}
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

func (svc *fileCreate) performCreate(opts FileCreateOptions, failOnDuplicateName bool, userID string) (*dto.File, error) {
	if len(opts.ParentID) > 0 {
		existing, err := svc.fileCoreSvc.getChildWithName(opts.ParentID, opts.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			if failOnDuplicateName {
				return nil, errorpkg.NewFileWithSimilarNameExistsError()
			} else {
				res, err := svc.fileMapper.Map(existing, userID)
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
	res, err := svc.fileMapper.Map(file, userID)
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
	fileCache       *cache.FileCache
	fileRepo        *repo.FileRepo
	fileSearch      *search.FileSearch
	fileGuard       *guard.FileGuard
	fileMapper      *mapper.FileMapper
	fileCoreSvc     *fileCoreService
	fileIdent       *infra.FileIdentifier
	userRepo        *repo.UserRepo
	workspaceRepo   *repo.WorkspaceRepo
	workspaceSvc    *WorkspaceService
	workspaceGuard  *guard.WorkspaceGuard
	workspaceMapper *mapper.WorkspaceMapper
}

func newFileFetch() *fileFetch {
	return &fileFetch{
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		fileIdent:   infra.NewFileIdentifier(),
		userRepo: repo.NewUserRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceSvc: NewWorkspaceService(),
		workspaceGuard: guard.NewWorkspaceGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceMapper: mapper.NewWorkspaceMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

func (svc *fileFetch) find(ids []string, userID string) ([]*dto.File, error) {
	var res []*dto.File
	for _, id := range ids {
		file, err := svc.fileCache.Get(id)
		if err != nil {
			continue
		}
		if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
			return nil, err
		}
		mapped, err := svc.fileMapper.Map(file, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, mapped)
	}
	return res, nil
}

func (svc *fileFetch) findByPath(path string, userID string) (*dto.File, error) {
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
	res, err := svc.fileMapper.Map(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *fileFetch) listByPath(path string, userID string) ([]*dto.File, error) {
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
		res, err := svc.fileMapper.MapMany(children, file.GetWorkspaceID(), userID)
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

func (svc *fileFetch) findPath(id string, userID string) ([]*dto.File, error) {
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
	res := make([]*dto.File, 0)
	for _, leaf := range path {
		if err = svc.fileGuard.Authorize(userID, leaf, model.PermissionViewer); err != nil {
			return nil, err
		}
		f, err := svc.fileMapper.Map(leaf, userID)
		if err != nil {
			return nil, err
		}
		res = append([]*dto.File{f}, res...)
	}
	return res, nil
}

func (svc *fileFetch) getPathString(files []*dto.File) string {
	return strings.Join(svc.getPathStrings(files), "/")
}

func (svc *fileFetch) getPathStringWithoutWorkspace(files []*dto.File) string {
	return strings.Join(svc.getPathStrings(files)[1:], "/")
}

func (svc *fileFetch) getPathStrings(files []*dto.File) []string {
	var components []string
	for _, f := range files {
		components = append(components, f.Name)
	}
	return components
}

func (svc *fileFetch) getWorkspacesAsFiles(userID string) ([]*dto.File, error) {
	loaded, err := svc.workspaceSvc.load(userID)
	if err != nil {
		return nil, err
	}
	mapped, err := svc.workspaceMapper.MapMany(loaded, userID)
	if err != nil {
		return nil, err
	}
	res := make([]*dto.File, 0)
	for _, w := range mapped {
		res = append(res, svc.getWorkspaceAsFile(w))
	}
	return res, nil
}

func (svc *fileFetch) getWorkspaceAsFile(workspace *dto.Workspace) *dto.File {
	return &dto.File{
		ID:         workspace.RootID,
		Workspace:  *workspace,
		Name:       svc.slugFromWorkspace(workspace.ID, workspace.Name),
		Type:       model.FileTypeFolder,
		Permission: workspace.Permission,
		CreateTime: workspace.CreateTime,
		UpdateTime: workspace.UpdateTime,
	}
}

func (svc *fileFetch) getUserAsFile(user model.User) *dto.File {
	return &dto.File{
		ID: user.GetID(),
		Workspace: dto.Workspace{
			ID:         user.GetID(),
			RootID:     user.GetID(),
			Name:       "root",
			Permission: model.PermissionOwner,
			CreateTime: user.GetCreateTime(),
		},
		Name:       "/",
		Type:       model.FileTypeFolder,
		Permission: model.PermissionOwner,
		CreateTime: user.GetCreateTime(),
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
	fileMapper     *mapper.FileMapper
	workspaceRepo  *repo.WorkspaceRepo
	workspaceGuard *guard.WorkspaceGuard
}

func newFileList() *fileList {
	return &fileList{
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc:   newFileCoreService(),
		fileFilterSvc: newFileFilterService(),
		fileSortSvc:   newFileSortService(),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceGuard: guard.NewWorkspaceGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

func (svc *fileList) probe(id string, opts FileListOptions, userID string) (*dto.FileProbe, error) {
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
	data, err := svc.getChildren(id)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.fileCoreSvc.authorize(userID, data, model.PermissionViewer)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(authorized))
	return &dto.FileProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *fileList) list(id string, opts FileListOptions, userID string) (*dto.FileList, error) {
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

func (svc *fileList) search(query *dto.FileQuery, workspace model.Workspace) ([]model.File, error) {
	var res []model.File
	count, err := svc.fileRepo.Count()
	if err != nil {
		return nil, err
	}
	filter := fmt.Sprintf("workspaceId=\"%s\"", workspace.GetID())
	if query.Type != nil {
		filter += fmt.Sprintf(" AND type=\"%s\"", *query.Type)
	}
	hits, err := svc.fileSearch.Query(*query.Text, infra.SearchQueryOptions{
		Limit:  count,
		Filter: filter,
	})
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

func (svc *fileList) createList(data []model.File, parent model.File, opts FileListOptions, userID string) (*dto.FileList, error) {
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
	if opts.SortBy == "" {
		opts.SortBy = dto.FileSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.FileSortOrderAsc
	}
	sorted := svc.fileSortSvc.sort(authorized, opts.SortBy, opts.SortOrder, userID)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.fileMapper.MapMany(paged, parent.GetWorkspaceID(), userID)
	if err != nil {
		return nil, err
	}
	res := &dto.FileList{
		Data:          mapped,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
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
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

func (svc *fileCompute) getSize(id string, userID string) (*int64, error) {
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

func (svc *fileCompute) getCount(id string, userID string) (*int64, error) {
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
	fileMapper   *mapper.FileMapper
	fileCoreSvc  *fileCoreService
	taskSvc      *TaskService
	snapshotRepo *repo.SnapshotRepo
}

func newFileCopy() *fileCopy {
	return &fileCopy{
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		taskSvc:     NewTaskService(),
		snapshotRepo: repo.NewSnapshotRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
	}
}

func (svc *fileCopy) copy(sourceID string, targetID string, userID string) (*dto.File, error) {
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
			logger.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(source, target, userID); err != nil {
		return nil, err
	}
	return svc.performCopy(source, target, userID)
}

func (svc *fileCopy) copyMany(opts dto.FileCopyManyOptions, userID string) (*dto.FileCopyManyResult, error) {
	res := &dto.FileCopyManyResult{
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

func (svc *fileCopy) performCopy(source model.File, target model.File, userID string) (*dto.File, error) {
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
	res, err := svc.fileMapper.Map(cloneResult.Root, userID)
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
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *fileCopy) index(clones []model.File) {
	if err := svc.fileSearch.Index(clones); err != nil {
		logger.GetLogger().Error(err)
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

type FileDelete struct {
	fileRepo       *repo.FileRepo
	fileSearch     *search.FileSearch
	fileGuard      *guard.FileGuard
	fileCache      *cache.FileCache
	workspaceCache *cache.WorkspaceCache
	taskSvc        *TaskService
	snapshotRepo   *repo.SnapshotRepo
	snapshotSvc    *SnapshotService
}

func NewFileDelete() *FileDelete {
	return &FileDelete{
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		taskSvc: NewTaskService(),
		snapshotRepo: repo.NewSnapshotRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		snapshotSvc: NewSnapshotService(),
	}
}

func (svc *FileDelete) delete(id string, userID string) error {
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
			logger.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(file); err != nil {
		return err
	}
	return svc.performDelete(file)
}

func (svc *FileDelete) deleteMany(opts dto.FileDeleteManyOptions, userID string) (*dto.FileDeleteManyResult, error) {
	res := &dto.FileDeleteManyResult{
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

func (svc *FileDelete) performDelete(file model.File) error {
	if file.GetType() == model.FileTypeFolder {
		return svc.DeleteFolder(file.GetID())
	} else if file.GetType() == model.FileTypeFile {
		return svc.deleteFile(file.GetID())
	}
	return nil
}

func (svc *FileDelete) check(file model.File) error {
	if file.GetParentID() == nil {
		workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		return errorpkg.NewCannotDeleteWorkspaceRootError(file, workspace)
	}
	return nil
}

func (svc *FileDelete) createTask(file model.File, userID string) (model.Task, error) {
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

func (svc *FileDelete) DeleteFolder(id string) error {
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

func (svc *FileDelete) deleteFile(id string) error {
	if err := svc.snapshotRepo.DeleteMappingsForTree(id); err != nil {
		logger.GetLogger().Error(err)
	}
	if err := svc.snapshotSvc.deleteForFile(id); err != nil {
		logger.GetLogger().Error(err)
	}
	if err := svc.fileCache.Delete(id); err != nil {
		logger.GetLogger().Error(err)
	}
	if err := svc.fileRepo.Delete(id); err != nil {
		logger.GetLogger().Error(err)
	}
	if err := svc.fileSearch.Delete([]string{id}); err != nil {
		logger.GetLogger().Error(err)
	}
	return nil
}

func (svc *FileDelete) deleteSnapshots(ids []string) {
	for _, id := range ids {
		if err := svc.snapshotSvc.deleteForFile(id); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *FileDelete) deleteFromRepo(ids []string) {
	const ChunkSize = 1000
	for i := 0; i < len(ids); i += ChunkSize {
		end := i + ChunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunk := ids[i:end]
		if err := svc.fileRepo.DeleteChunk(chunk); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *FileDelete) deleteFromCache(ids []string) {
	for _, id := range ids {
		if err := svc.fileCache.Delete(id); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

func (svc *FileDelete) deleteFromSearch(ids []string) {
	if err := svc.fileSearch.Delete(ids); err != nil {
		logger.GetLogger().Error(err)
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
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		snapshotCache: cache.NewSnapshotCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		s3: infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
	}
}

type DownloadResult struct {
	File          model.File
	Snapshot      model.Snapshot
	RangeInterval *helper.RangeInterval
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

func (svc *fileDownload) downloadTextBuffer(id string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasText() {
		rangeInterval, err := svc.downloadS3Object(snapshot.GetText(), "", buf)
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

func (svc *fileDownload) downloadOCRBuffer(id string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot.HasOCR() {
		rangeInterval, err := svc.downloadS3Object(snapshot.GetOCR(), "", buf)
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

func (svc *fileDownload) downloadS3Object(s3Object *model.S3Object, rangeHeader string, buf *bytes.Buffer) (*helper.RangeInterval, error) {
	objectInfo, err := svc.s3.StatObject(s3Object.Key, s3Object.Bucket, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	opts := minio.GetObjectOptions{}
	var rangeInterval *helper.RangeInterval
	if rangeHeader != "" {
		rangeInterval = helper.NewRangeInterval(rangeHeader, objectInfo.Size)
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
	fileMapper  *mapper.FileMapper
	fileCoreSvc *fileCoreService
	taskSvc     *TaskService
}

func newFileMove() *fileMove {
	return &fileMove{
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		taskSvc:     NewTaskService(),
	}
}

func (svc *fileMove) move(sourceID string, targetID string, userID string) (*dto.File, error) {
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
			logger.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(source, target, userID); err != nil {
		return nil, err
	}
	return svc.performMove(source, target, userID)
}

func (svc *fileMove) moveMany(opts dto.FileMoveManyOptions, userID string) (*dto.FileMoveManyResult, error) {
	res := &dto.FileMoveManyResult{
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

func (svc *fileMove) performMove(source model.File, target model.File, userID string) (*dto.File, error) {
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
	res, err := svc.fileMapper.Map(source, userID)
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
	fileMapper  *mapper.FileMapper
}

func newFilePatch() *filePatch {
	return &filePatch{
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

func (svc *filePatch) patchName(id string, name string, userID string) (*dto.File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
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
	file.SetName(name)
	if err = svc.fileRepo.Save(file); err != nil {
		return nil, err
	}
	if err := svc.fileCoreSvc.sync(file); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.Map(file, userID)
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
	userMapper     *mapper.UserMapper
	workspaceRepo  *repo.WorkspaceRepo
	workspaceCache *cache.WorkspaceCache
	groupCache     *cache.GroupCache
	groupGuard     *guard.GroupGuard
	groupMapper    *mapper.GroupMapper
	permissionRepo *repo.PermissionRepo
}

func newFilePermission() *filePermission {
	return &filePermission{
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		userRepo: repo.NewUserRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		userMapper: mapper.NewUserMapper(),
		workspaceRepo: repo.NewWorkspaceRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupCache: cache.NewGroupCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupGuard: guard.NewGroupGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		groupMapper: mapper.NewGroupMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		permissionRepo: repo.NewPermissionRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
	}
}

func (svc *filePermission) grantUserPermissions(ids []string, assigneeID string, permission string, userID string) error {
	for _, id := range ids {
		if err := svc.grantUserPermission(id, assigneeID, permission, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) grantUserPermission(id string, assigneeID string, permission string, userID string) error {
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

func (svc *filePermission) revokeUserPermissions(ids []string, assigneeID string, userID string) error {
	for _, id := range ids {
		if err := svc.revokeUserPermission(id, assigneeID, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) revokeUserPermission(id string, assigneeID string, userID string) error {
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

func (svc *filePermission) grantGroupPermissions(ids []string, groupID string, permission string, userID string) error {
	for _, id := range ids {
		if err := svc.grantGroupPermission(id, groupID, permission, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) grantGroupPermission(id string, groupID string, permission string, userID string) error {
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

func (svc *filePermission) revokeGroupPermissions(ids []string, groupID string, userID string) error {
	for _, id := range ids {
		if err := svc.revokeGroupPermission(id, groupID, userID); err != nil {
			return err
		}
	}
	return nil
}

func (svc *filePermission) revokeGroupPermission(id string, groupID string, userID string) error {
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

func (svc *filePermission) findUserPermissions(id string, userID string) ([]*dto.UserPermission, error) {
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
	res := make([]*dto.UserPermission, 0)
	for _, p := range permissions {
		if p.GetUserID() == userID {
			continue
		}
		u, err := svc.userRepo.Find(p.GetUserID())
		if err != nil {
			return nil, err
		}
		res = append(res, &dto.UserPermission{
			ID:         p.GetID(),
			User:       svc.userMapper.Map(u),
			Permission: p.GetPermission(),
		})
	}
	return res, nil
}

func (svc *filePermission) findGroupPermissions(id string, userID string) ([]*dto.GroupPermission, error) {
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
	res := make([]*dto.GroupPermission, 0)
	for _, p := range permissions {
		m, err := svc.groupCache.Get(p.GetGroupID())
		if err != nil {
			return nil, err
		}
		g, err := svc.groupMapper.Map(m, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, &dto.GroupPermission{
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
	fileCoreSvc    *fileCoreService
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	taskCache      *cache.TaskCache
	taskSvc        *TaskService
	fileIdent      *infra.FileIdentifier
	pipelineClient client.PipelineClient
}

func newFileReprocess() *fileReprocess {
	return &fileReprocess{
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		snapshotCache: cache.NewSnapshotCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		snapshotSvc: NewSnapshotService(),
		taskCache: cache.NewTaskCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		taskSvc:   NewTaskService(),
		fileIdent: infra.NewFileIdentifier(),
		pipelineClient: client.NewPipelineClient(
			config.GetConfig().ConversionURL,
			config.GetConfig().Environment.IsTest,
		),
	}
}

func (svc *fileReprocess) reprocess(id string, userID string) (*dto.FileReprocessResult, error) {
	resp := &dto.FileReprocessResult{
		// We intend to send an empty array to the caller, better than nil
		Accepted: []string{},
		Rejected: []string{},
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
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
		if snapshot.GetOriginal().Size <= helper.MegabyteToByte(svc.fileCoreSvc.getProcessingLimitMB(leaf.GetName())) {
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
		logger.GetLogger().Error(err)
		return false
	}
	snapshot, err := svc.snapshotCache.Get(*leaf.GetSnapshotID())
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}
	if !svc.check(leaf, snapshot) {
		return false
	}
	if err := svc.runPipeline(leaf, snapshot, userID); err != nil {
		logger.GetLogger().Error(err)
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
			logger.GetLogger().Error(err)
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
	if err := svc.pipelineClient.Run(&dto.PipelineRunOptions{
		TaskID:     helper.ToPtr(task.GetID()),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        snapshot.GetOriginal().Key,
		Intent:     snapshot.GetIntent(),
		Language:   snapshot.GetLanguage(),
	}); err != nil {
		return err
	}
	return nil
}

type fileStore struct {
	fileCache             *cache.FileCache
	fileCoreSvc           *fileCoreService
	fileMapper            *mapper.FileMapper
	workspaceCache        *cache.WorkspaceCache
	snapshotRepo          *repo.SnapshotRepo
	snapshotCache         *cache.SnapshotCache
	snapshotSvc           *SnapshotService
	snapshotMapper        *mapper.SnapshotMapper
	snapshotWebhookClient *client.SnapshotWebhookClient
	taskSvc               *TaskService
	fileIdent             *infra.FileIdentifier
	s3                    infra.S3Manager
	pipelineClient        client.PipelineClient
	config                *config.Config
}

func newFileStore() *fileStore {
	return &fileStore{
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileCoreSvc: newFileCoreService(),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		workspaceCache: cache.NewWorkspaceCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		snapshotRepo: repo.NewSnapshotRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		snapshotCache: cache.NewSnapshotCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		snapshotSvc: NewSnapshotService(),
		snapshotMapper: mapper.NewSnapshotMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		snapshotWebhookClient: client.NewSnapshotWebhookClient(
			config.GetConfig().Security,
		),
		taskSvc:   NewTaskService(),
		fileIdent: infra.NewFileIdentifier(),
		s3:        infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		pipelineClient: client.NewPipelineClient(
			config.GetConfig().ConversionURL,
			config.GetConfig().Environment.IsTest,
		),
		config: config.GetConfig(),
	}
}

type FileStoreOptions struct {
	S3Reference *model.S3Reference
	Path        *string
}

func (svc *fileStore) store(id string, opts FileStoreOptions, userID string) (*dto.File, error) {
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
	snapshot, err = svc.callSnapshotHookWithCreateEvent(snapshot)
	if err != nil {
		return nil, err
	}
	if !props.ExceedsProcessingLimit {
		if err := svc.runPipeline(file, snapshot, props, userID); err != nil {
			return nil, err
		}
	}
	res, err := svc.fileMapper.Map(file, userID)
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
	props.ExceedsProcessingLimit = props.Size > helper.MegabyteToByte(svc.fileCoreSvc.getProcessingLimitMB(props.Path))
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
			Size:   stat.Size(),
		},
		Bucket:      workspace.GetBucket(),
		ContentType: helper.DetectMIMEFromPath(*opts.Path),
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
			Size:   opts.S3Reference.Size,
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

func (svc *fileStore) callSnapshotHookWithCreateEvent(snapshot model.Snapshot) (model.Snapshot, error) {
	if svc.config.SnapshotWebhook != "" {
		if err := svc.snapshotWebhookClient.Call(config.GetConfig().SnapshotWebhook, dto.SnapshotWebhookOptions{
			EventType: dto.SnapshotWebhookEventTypeCreate,
			Snapshot:  svc.snapshotMapper.MapWithS3Objects(snapshot),
		}); err != nil {
			logger.GetLogger().Error(err)
		} else {
			snapshot, err = svc.snapshotCache.Get(snapshot.GetID())
			if err != nil {
				return nil, err
			}
		}
	}
	return snapshot, nil
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
	if err := svc.pipelineClient.Run(&dto.PipelineRunOptions{
		TaskID:     helper.ToPtr(task.GetID()),
		SnapshotID: snapshot.GetID(),
		Bucket:     props.Original.Bucket,
		Key:        props.Original.Key,
		Intent:     snapshot.GetIntent(),
		Language:   snapshot.GetLanguage(),
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
	fileIdent  *infra.FileIdentifier
	config     *config.Config
}

func newFileCoreService() *fileCoreService {
	return &fileCoreService{
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileSearch: search.NewFileSearch(
			config.GetConfig().Postgres,
			config.GetConfig().Search,
			config.GetConfig().S3,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileIdent: infra.NewFileIdentifier(),
		config:    config.GetConfig(),
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

func (svc *fileCoreService) getProcessingLimitMB(path string) int {
	var res int
	if svc.fileIdent.IsAudio(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeAudio)
	} else if svc.fileIdent.IsImage(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeImage)
	} else if svc.fileIdent.IsOffice(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeOffice)
	} else if svc.fileIdent.IsPDF(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypePDF)
	} else if svc.fileIdent.IsPlainText(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypePlainText)
	} else if svc.fileIdent.IsVideo(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeVideo)
	} else if svc.fileIdent.IsGLB(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeGLB)
	} else if svc.fileIdent.IsZIP(path) {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeZIP)
	} else if ok, err := svc.fileIdent.IsGLTF(path); ok && err != nil {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeGLTF)
	} else {
		res = svc.config.Limits.GetFileProcessingMB(config.FileTypeEverythingElse)
	}
	return res
}

type fileFilterService struct {
	fileRepo   *repo.FileRepo
	fileGuard  *guard.FileGuard
	fileMapper *mapper.FileMapper
	fileIdent  *infra.FileIdentifier
}

func newFileFilterService() *fileFilterService {
	return &fileFilterService{
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileGuard: guard.NewFileGuard(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileIdent: infra.NewFileIdentifier(),
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
			f, err := svc.fileMapper.Map(file.(model.File), userID)
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
			f, err := svc.fileMapper.Map(file.(model.File), userID)
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
			f, err := svc.fileMapper.Map(file.(model.File), userID)
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
			f, err := svc.fileMapper.Map(file.(model.File), userID)
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
			f, err := svc.fileMapper.Map(file.(model.File), userID)
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
			f, err := svc.fileMapper.Map(file.(model.File), userID)
			if err != nil {
				return false
			}
			if f.Snapshot == nil {
				return true
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

func (svc *fileFilterService) filterWithQuery(data []model.File, opts dto.FileQuery, parent model.File) ([]model.File, error) {
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
				return helper.StringToTimestamp(v.(model.File).GetCreateTime()) >= *opts.CreateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.CreateTimeBefore != nil {
				return helper.StringToTimestamp(v.(model.File).GetCreateTime()) <= *opts.CreateTimeBefore
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeAfter != nil {
				file := v.(model.File)
				return file.GetUpdateTime() != nil && helper.StringToTimestamp(*file.GetUpdateTime()) >= *opts.UpdateTimeAfter
			} else {
				return true
			}
		}).
		Filter(func(v interface{}) bool {
			if opts.UpdateTimeBefore != nil {
				file := v.(model.File)
				return file.GetUpdateTime() != nil && helper.StringToTimestamp(*file.GetUpdateTime()) <= *opts.UpdateTimeBefore
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
	fileMapper    *mapper.FileMapper
	fileFilterSvc *fileFilterService
}

func newFileSortService() *fileSortService {
	return &fileSortService{
		fileMapper: mapper.NewFileMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		fileFilterSvc: newFileFilterService(),
	}
}

func (svc *fileSortService) sort(data []model.File, sortBy string, sortOrder string, userID string) []model.File {
	if sortBy == dto.FileSortByName {
		return svc.sortByName(data, sortOrder)
	} else if sortBy == dto.FileSortBySize {
		return svc.sortBySize(data, sortOrder, userID)
	} else if sortBy == dto.FileSortByDateCreated {
		return svc.sortByDateCreated(data, sortOrder)
	} else if sortBy == dto.FileSortByDateModified {
		return svc.sortByDateModified(data, sortOrder)
	} else if sortBy == dto.FileSortByKind {
		return svc.sortByKind(data, userID)
	}
	return data
}

func (svc *fileSortService) sortByName(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if sortOrder == dto.FileSortOrderDesc {
			return data[i].GetName() > data[j].GetName()
		} else {
			return data[i].GetName() < data[j].GetName()
		}
	})
	return data
}

func (svc *fileSortService) sortBySize(data []model.File, sortOrder string, userID string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		fileA, err := svc.fileMapper.Map(data[i], userID)
		if err != nil {
			return false
		}
		fileB, err := svc.fileMapper.Map(data[j], userID)
		if err != nil {
			return false
		}
		var sizeA int64
		if fileA.Snapshot != nil && fileA.Snapshot.Original != nil {
			sizeA = fileA.Snapshot.Original.Size
		}
		var sizeB int64
		if fileB.Snapshot != nil && fileB.Snapshot.Original != nil {
			sizeB = fileB.Snapshot.Original.Size
		}
		if sortOrder == dto.FileSortOrderDesc {
			return sizeA > sizeB
		} else {
			return sizeA < sizeB
		}
	})
	return data
}

func (svc *fileSortService) sortByDateCreated(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		a := helper.StringToTimestamp(data[i].GetCreateTime())
		b := helper.StringToTimestamp(data[j].GetCreateTime())
		if sortOrder == dto.FileSortOrderDesc {
			return a > b
		} else {
			return a < b
		}
	})
	return data
}

func (svc *fileSortService) sortByDateModified(data []model.File, sortOrder string) []model.File {
	sort.Slice(data, func(i, j int) bool {
		if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
			a := helper.StringToTimestamp(*data[i].GetUpdateTime())
			b := helper.StringToTimestamp(*data[j].GetUpdateTime())
			if sortOrder == dto.FileSortOrderDesc {
				return a > b
			} else {
				return a < b
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
		value == dto.FileSortByName ||
		value == dto.FileSortByKind ||
		value == dto.FileSortBySize ||
		value == dto.FileSortByDateCreated ||
		value == dto.FileSortByDateModified
}

func (svc *fileSortService) isValidSortOrder(value string) bool {
	return value == "" || value == dto.FileSortOrderAsc || value == dto.FileSortOrderDesc
}
