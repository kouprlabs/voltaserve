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

	"github.com/kouprlabs/voltaserve/api/model"
)

type FileService struct {
	coreSvc       *FileCore
	createSvc     *FileCreate
	storeSvc      *FileStore
	deleteSvc     *FileDelete
	moveSvc       *FileMove
	copySvc       *FileCopy
	downloadSvc   *FileDownload
	fetchSvc      *FileFetch
	reprocessSvc  *FileReprocess
	permissionSvc *FilePermission
	computeSvc    *FileCompute
	patchSvc      *FilePatch
}

func NewFileService() *FileService {
	return &FileService{
		coreSvc:       NewFileCore(),
		createSvc:     NewFileCreate(),
		storeSvc:      NewFileStore(),
		deleteSvc:     NewFileDelete(),
		moveSvc:       NewFileMove(),
		copySvc:       NewFileCopy(),
		downloadSvc:   NewFileDownload(),
		fetchSvc:      NewFileFind(),
		reprocessSvc:  NewFileReprocess(),
		permissionSvc: NewFilePermission(),
		computeSvc:    NewFileCompute(),
		patchSvc:      NewFilePatch(),
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

func (svc *FileService) Create(opts FileCreateOptions, userID string) (*File, error) {
	return svc.createSvc.Create(opts, userID)
}

func (svc *FileService) Store(id string, opts StoreOptions, userID string) (*File, error) {
	return svc.storeSvc.Store(id, opts, userID)
}

func (svc *FileService) DownloadOriginalBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.downloadSvc.DownloadOriginalBuffer(id, rangeHeader, buf, userID)
}

func (svc *FileService) DownloadPreviewBuffer(id string, rangeHeader string, buf *bytes.Buffer, userID string) (*DownloadResult, error) {
	return svc.downloadSvc.DownloadPreviewBuffer(id, rangeHeader, buf, userID)
}

func (svc *FileService) DownloadThumbnailBuffer(id string, buf *bytes.Buffer, userID string) (model.Snapshot, error) {
	return svc.downloadSvc.DownloadThumbnailBuffer(id, buf, userID)
}

func (svc *FileService) Find(ids []string, userID string) ([]*File, error) {
	return svc.fetchSvc.Find(ids, userID)
}

func (svc *FileService) FindByPath(path string, userID string) (*File, error) {
	return svc.fetchSvc.FindByPath(path, userID)
}

func (svc *FileService) ListByPath(path string, userID string) ([]*File, error) {
	return svc.fetchSvc.ListByPath(path, userID)
}

func (svc *FileService) Probe(id string, opts FileListOptions, userID string) (*FileProbe, error) {
	return svc.fetchSvc.Probe(id, opts, userID)
}

func (svc *FileService) List(id string, opts FileListOptions, userID string) (*FileList, error) {
	return svc.fetchSvc.List(id, opts, userID)
}

func (svc *FileService) FindPath(id string, userID string) ([]*File, error) {
	return svc.fetchSvc.FindPath(id, userID)
}

func (svc *FileService) CopyOne(sourceID string, targetID string, userID string) (*File, error) {
	return svc.copySvc.CopyOne(sourceID, targetID, userID)
}

func (svc *FileService) CopyMany(opts FileCopyManyOptions, userID string) (*FileCopyManyResult, error) {
	return svc.copySvc.CopyMany(opts, userID)
}

func (svc *FileService) MoveOne(sourceID string, targetID string, userID string) (*File, error) {
	return svc.moveSvc.MoveOne(sourceID, targetID, userID)
}

func (svc *FileService) MoveMany(opts FileMoveManyOptions, userID string) (*FileMoveManyResult, error) {
	return svc.moveSvc.MoveMany(opts, userID)
}

func (svc *FileService) PatchName(id string, name string, userID string) (*File, error) {
	return svc.patchSvc.PatchName(id, name, userID)
}

func (svc *FileService) Reprocess(id string, userID string) (res *ReprocessResponse, err error) {
	return svc.reprocessSvc.Reprocess(id, userID)
}

func (svc *FileService) DeleteOne(id string, userID string) error {
	return svc.deleteSvc.DeleteOne(id, userID)
}

func (svc *FileService) DeleteMany(opts FileDeleteManyOptions, userID string) (*FileDeleteManyResult, error) {
	return svc.deleteSvc.DeleteMany(opts, userID)
}

func (svc *FileService) ComputeSize(id string, userID string) (*int64, error) {
	return svc.computeSvc.ComputeSize(id, userID)
}

func (svc *FileService) Count(id string, userID string) (*int64, error) {
	return svc.computeSvc.Count(id, userID)
}

func (svc *FileService) GrantUserPermission(ids []string, assigneeID string, permission string, userID string) error {
	return svc.permissionSvc.GrantUserPermission(ids, assigneeID, permission, userID)
}

func (svc *FileService) RevokeUserPermission(ids []string, assigneeID string, userID string) error {
	return svc.permissionSvc.RevokeUserPermission(ids, assigneeID, userID)
}

func (svc *FileService) GrantGroupPermission(ids []string, groupID string, permission string, userID string) error {
	return svc.permissionSvc.GrantGroupPermission(ids, groupID, permission, userID)
}

func (svc *FileService) RevokeGroupPermission(ids []string, groupID string, userID string) error {
	return svc.permissionSvc.RevokeGroupPermission(ids, groupID, userID)
}

func (svc *FileService) FindUserPermissions(id string, userID string) ([]*UserPermission, error) {
	return svc.permissionSvc.FindUserPermissions(id, userID)
}

func (svc *FileService) FindGroupPermissions(id string, userID string) ([]*GroupPermission, error) {
	return svc.permissionSvc.FindGroupPermissions(id, userID)
}
