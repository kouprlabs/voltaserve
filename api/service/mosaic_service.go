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
	"path/filepath"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type MosaicService struct {
	snapshotCache  *cache.SnapshotCache
	snapshotRepo   repo.SnapshotRepo
	snapshotSvc    *SnapshotService
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	taskSvc        *TaskService
	s3             *infra.S3Manager
	mosaicClient   *client.MosaicClient
	pipelineClient *client.PipelineClient
	fileIdent      *infra.FileIdentifier
}

func NewMosaicService() *MosaicService {
	return &MosaicService{
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		snapshotSvc:    NewSnapshotService(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		taskSvc:        NewTaskService(),
		s3:             infra.NewS3Manager(),
		mosaicClient:   client.NewMosaicClient(),
		pipelineClient: client.NewPipelineClient(),
		fileIdent:      infra.NewFileIdentifier(),
	}
}

func (svc *MosaicService) Create(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	isTaskPending, err := svc.snapshotSvc.IsTaskPending(snapshot)
	if err != nil {
		return err
	}
	if *isTaskPending {
		return errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Waiting.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusWaiting,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return err
	}
	snapshot.SetStatus(model.SnapshotStatusWaiting)
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
		return err
	}
	if err := svc.pipelineClient.Run(&client.PipelineRunOptions{
		PipelineID: helper.ToPtr(client.PipelineMosaic),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetOriginal().Bucket,
		Key:        snapshot.GetOriginal().Key,
	}); err != nil {
		return err
	}
	return nil
}

func (svc *MosaicService) Delete(id string, userID string) error {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	isTaskPending, err := svc.snapshotSvc.IsTaskPending(snapshot)
	if err != nil {
		return err
	}
	if *isTaskPending {
		return errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	if !snapshot.HasMosaic() {
		return errorpkg.NewMosaicNotFoundError(nil)
	}
	if svc.fileIdent.IsImage(snapshot.GetOriginal().Key) {
		task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
			ID:              helper.NewID(),
			Name:            "Deleting mosaic.",
			UserID:          userID,
			IsIndeterminate: true,
			Status:          model.TaskStatusRunning,
			Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
		})
		if err != nil {
			return err
		}
		snapshot.SetTaskID(helper.ToPtr(task.GetID()))
		snapshot.SetStatus(model.SnapshotStatusProcessing)
		if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
			return err
		}
		go func(task model.Task, snapshot model.Snapshot) {
			err = svc.mosaicClient.Delete(client.MosaicDeleteOptions{
				S3Key:    filepath.FromSlash(snapshot.GetID()),
				S3Bucket: snapshot.GetOriginal().Bucket,
			})
			if err != nil {
				value := err.Error()
				task.SetError(&value)
				if err := svc.taskSvc.saveAndSync(task); err != nil {
					log.GetLogger().Error(err)
					return
				}
			} else {
				if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
					log.GetLogger().Error(err)
					return
				}
			}
			snapshot.SetMosaic(nil)
			snapshot.SetTaskID(nil)
			snapshot.SetStatus(model.SnapshotStatusReady)
			if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
				log.GetLogger().Error(err)
				return
			}
		}(task, snapshot)
	}
	return nil
}

type MosaicInfo struct {
	IsAvailable bool                   `json:"isAvailable"`
	IsOutdated  bool                   `json:"isOutdated"`
	Snapshot    *Snapshot              `json:"snapshot,omitempty"`
	Metadata    *client.MosaicMetadata `json:"metadata,omitempty"`
}

func (svc *MosaicService) GetInfo(id string, userID string) (*MosaicInfo, error) {
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
	isOutdated := false
	if !snapshot.HasMosaic() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous == nil {
			return &MosaicInfo{IsAvailable: false}, nil
		} else {
			snapshot = previous
			isOutdated = true
		}
	}
	res, err := svc.mosaicClient.GetMetadata(client.MosaicGetMetadataOptions{
		S3Key:    filepath.FromSlash(snapshot.GetID()),
		S3Bucket: snapshot.GetOriginal().Bucket,
	})
	if err != nil {
		return nil, err
	}
	return &MosaicInfo{
		IsAvailable: true,
		IsOutdated:  isOutdated,
		Snapshot:    svc.snapshotSvc.snapshotMapper.mapOne(snapshot),
		Metadata:    res,
	}, nil
}

type MosaicDownloadTileOptions struct {
	ZoomLevel int
	Row       int
	Col       int
	Ext       string
}

func (svc *MosaicService) DownloadTileBuffer(id string, opts MosaicDownloadTileOptions, userID string) (*bytes.Buffer, error) {
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
	if !snapshot.HasMosaic() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous == nil {
			return nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	res, err := svc.mosaicClient.DownloadTileBuffer(client.MosaicDownloadTileOptions{
		S3Key:     filepath.FromSlash(snapshot.GetID()),
		S3Bucket:  snapshot.GetOriginal().Bucket,
		ZoomLevel: opts.ZoomLevel,
		Row:       opts.Row,
		Col:       opts.Col,
		Ext:       opts.Ext,
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (svc *MosaicService) getPreviousSnapshot(fileID string, version int64) (model.Snapshot, error) {
	snapshots, err := svc.snapshotRepo.FindAllPrevious(fileID, version)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots {
		if snapshot.HasMosaic() {
			return snapshot, nil
		}
	}
	return nil, nil
}
