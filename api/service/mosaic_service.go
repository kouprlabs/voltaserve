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
	"path/filepath"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/guard"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/mapper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type MosaicService struct {
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	taskSvc        *TaskService
	taskMapper     *mapper.TaskMapper
	mosaicClient   *client.MosaicClient
	pipelineClient client.PipelineClient
}

func NewMosaicService() *MosaicService {
	return &MosaicService{
		snapshotCache: cache.NewSnapshotCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		snapshotSvc: NewSnapshotService(),
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
		taskSvc: NewTaskService(),
		taskMapper: mapper.NewTaskMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		mosaicClient: client.NewMosaicClient(config.GetConfig().MosaicURL),
		pipelineClient: client.NewPipelineClient(
			config.GetConfig().ConversionURL,
			config.GetConfig().Environment.IsTest,
		),
	}
}

func (svc *MosaicService) Create(fileID string, userID string) (*dto.Task, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionEditor); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	isTaskPending, err := svc.snapshotSvc.isTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	task, err := svc.createWaitingTask(file, userID)
	if err != nil {
		return nil, err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return nil, err
	}
	if err := svc.runPipeline(snapshot, task); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.MapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *MosaicService) Delete(fileID string, userID string) (*dto.Task, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionOwner); err != nil {
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
		return nil, errorpkg.NewMosaicNotFoundError(nil)
	}
	isTaskPending, err := svc.snapshotSvc.isTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	task, err := svc.createDeleteTask(file, userID)
	if err != nil {
		return nil, err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return nil, err
	}
	go svc.delete(task, snapshot)
	res, err := svc.taskMapper.MapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *MosaicService) GetMetadata(fileID string, userID string) (*dto.MosaicMetadata, error) {
	file, err := svc.fileCache.Get(fileID)
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
		return nil, errorpkg.NewMosaicNotFoundError(nil)
	}
	res, err := svc.mosaicClient.GetMetadata(client.MosaicGetMetadataOptions{
		S3Key:    filepath.FromSlash(snapshot.GetID()),
		S3Bucket: snapshot.GetPreview().Bucket,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

type MosaicDownloadTileOptions struct {
	ZoomLevel int
	Row       int
	Column    int
	Extension string
}

func (svc *MosaicService) DownloadTileBuffer(fileID string, opts MosaicDownloadTileOptions, userID string) ([]byte, model.Snapshot, error) {
	file, err := svc.fileCache.Get(fileID)
	if err != nil {
		return nil, nil, err
	}
	if err = svc.fileGuard.Authorize(userID, file, model.PermissionViewer); err != nil {
		return nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotCache.Get(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, err
	}
	if !snapshot.HasMosaic() {
		return nil, nil, errorpkg.NewMosaicNotFoundError(nil)
	}
	res, err := svc.mosaicClient.DownloadTileBuffer(client.MosaicDownloadTileOptions{
		S3Key:     filepath.FromSlash(snapshot.GetID()),
		S3Bucket:  snapshot.GetPreview().Bucket,
		ZoomLevel: opts.ZoomLevel,
		Row:       opts.Row,
		Column:    opts.Column,
		Extension: opts.Extension,
	})
	if err != nil {
		return nil, nil, err
	}
	return res, snapshot, err
}

func (svc *MosaicService) runPipeline(snapshot model.Snapshot, task model.Task) error {
	if err := svc.pipelineClient.Run(&dto.PipelineRunOptions{
		PipelineID: helper.ToPtr(dto.PipelineMosaic),
		TaskID:     helper.ToPtr(task.GetID()),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetPreview().Bucket,
		Key:        snapshot.GetPreview().Key,
		Intent:     snapshot.GetIntent(),
		Language:   snapshot.GetLanguage(),
	}); err != nil {
		return err
	}
	return nil
}

func (svc *MosaicService) createWaitingTask(file model.File, userID string) (model.Task, error) {
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

func (svc *MosaicService) createDeleteTask(file model.File, userID string) (model.Task, error) {
	res, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Deleting mosaic.",
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

func (svc *MosaicService) delete(task model.Task, snapshot model.Snapshot) {
	err := svc.mosaicClient.Delete(client.MosaicDeleteOptions{
		S3Key:    filepath.FromSlash(snapshot.GetID()),
		S3Bucket: snapshot.GetPreview().Bucket,
	})
	if err != nil {
		value := err.Error()
		task.SetError(&value)
		if err := svc.taskSvc.saveAndSync(task); err != nil {
			logger.GetLogger().Error(err)
			return
		}
	} else {
		if err := svc.taskSvc.deleteAndSync(task.GetID()); err != nil {
			logger.GetLogger().Error(err)
			return
		}
	}
	snapshot.SetMosaic(nil)
	snapshot.SetTaskID(nil)
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		logger.GetLogger().Error(err)
		return
	}
}
