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
	"path/filepath"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/client/mosaic_client"
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
	taskMapper     *taskMapper
	s3             *infra.S3Manager
	mosaicClient   *mosaic_client.MosaicClient
	pipelineClient *conversion_client.PipelineClient
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
		taskMapper:     newTaskMapper(),
		s3:             infra.NewS3Manager(),
		mosaicClient:   mosaic_client.NewMosaicClient(),
		pipelineClient: conversion_client.NewPipelineClient(),
		fileIdent:      infra.NewFileIdentifier(),
	}
}

func (svc *MosaicService) Create(id string, userID string) (*Task, error) {
	file, err := svc.fileCache.Get(id)
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
	snapshot.SetStatus(model.SnapshotStatusWaiting)
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return nil, err
	}
	if err := svc.runPipeline(task, snapshot); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *MosaicService) runPipeline(task model.Task, snapshot model.Snapshot) error {
	if err := svc.pipelineClient.Run(&conversion_client.PipelineRunOptions{
		PipelineID: helper.ToPtr(conversion_client.PipelineMosaic),
		TaskID:     task.GetID(),
		SnapshotID: snapshot.GetID(),
		Bucket:     snapshot.GetPreview().Bucket,
		Key:        snapshot.GetPreview().Key,
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

func (svc *MosaicService) Delete(id string, userID string) (*Task, error) {
	file, err := svc.fileCache.Get(id)
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
	isTaskPending, err := svc.snapshotSvc.isTaskPending(snapshot)
	if err != nil {
		return nil, err
	}
	if isTaskPending {
		return nil, errorpkg.NewSnapshotHasPendingTaskError(nil)
	}
	if !snapshot.HasMosaic() {
		return nil, errorpkg.NewMosaicNotFoundError(nil)
	}
	task, err := svc.createDeleteTask(file, userID)
	if err != nil {
		return nil, err
	}
	snapshot.SetMosaic(nil)
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	snapshot.SetStatus(model.SnapshotStatusProcessing)
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		return nil, err
	}
	go svc.delete(task, snapshot)
	res, err := svc.taskMapper.mapOne(task)
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
	err := svc.mosaicClient.Delete(mosaic_client.MosaicDeleteOptions{
		S3Key:    filepath.FromSlash(snapshot.GetID()),
		S3Bucket: snapshot.GetPreview().Bucket,
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
	if err := svc.snapshotSvc.saveAndSync(snapshot); err != nil {
		log.GetLogger().Error(err)
		return
	}
}

func (svc *MosaicService) ReadInfo(id string, userID string) (*MosaicInfo, error) {
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
	res, err := svc.mosaicClient.GetMetadata(mosaic_client.MosaicGetMetadataOptions{
		S3Key:    filepath.FromSlash(snapshot.GetID()),
		S3Bucket: snapshot.GetPreview().Bucket,
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
	Column    int
	Extension string
}

func (svc *MosaicService) DownloadTileBuffer(id string, opts MosaicDownloadTileOptions, userID string) (*bytes.Buffer, model.Snapshot, error) {
	file, err := svc.fileCache.Get(id)
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
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, nil, err
		}
		if previous == nil {
			return nil, nil, errorpkg.NewMosaicNotFoundError(nil)
		} else {
			snapshot = previous
		}
	}
	res, err := svc.mosaicClient.DownloadTileBuffer(mosaic_client.MosaicDownloadTileOptions{
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
