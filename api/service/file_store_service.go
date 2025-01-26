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
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FileStoreService struct {
	fileCache      *cache.FileCache
	fileCoreSvc    *FileCoreService
	fileMapper     *FileMapper
	workspaceCache *cache.WorkspaceCache
	snapshotRepo   repo.SnapshotRepo
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	taskSvc        *TaskService
	fileIdent      *infra.FileIdentifier
	s3             *infra.S3Manager
	pipelineClient *conversion_client.PipelineClient
}

func NewFileStoreService() *FileStoreService {
	return &FileStoreService{
		fileCache:      cache.NewFileCache(),
		fileCoreSvc:    NewFileCoreService(),
		fileMapper:     NewFileMapper(),
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

type StoreOptions struct {
	S3Reference *model.S3Reference
	Path        *string
}

func (svc *FileStoreService) Store(id string, opts StoreOptions, userID string) (*File, error) {
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	props, err := svc.getProperties(file, opts)
	if err != nil {
		return nil, err
	}
	if opts.S3Reference == nil {
		if err := svc.store(props); err != nil {
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
		if err := svc.process(file, snapshot, props, userID); err != nil {
			return nil, err
		}
	}
	res, err := svc.fileMapper.mapOne(file, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type storeProperties struct {
	SnapshotID             string
	Size                   int64
	Path                   string
	Original               model.S3Object
	Bucket                 string
	ContentType            string
	ExceedsProcessingLimit bool
}

func (svc *FileStoreService) getProperties(file model.File, opts StoreOptions) (storeProperties, error) {
	props := storeProperties{}
	if opts.S3Reference == nil {
		var err error
		props, err = svc.getPropertiesFromPath(file, opts)
		if err != nil {
			return storeProperties{}, err
		}
	} else {
		props = svc.getPropertiesFromS3Reference(opts)
	}
	props.ExceedsProcessingLimit = props.Size > helper.MegabyteToByte(svc.fileIdent.GetProcessingLimitMB(props.Path))
	return props, nil
}

func (svc *FileStoreService) getPropertiesFromPath(file model.File, opts StoreOptions) (storeProperties, error) {
	stat, err := os.Stat(*opts.Path)
	if err != nil {
		return storeProperties{}, err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return storeProperties{}, err
	}
	snapshotID := helper.NewID()
	return storeProperties{
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

func (svc *FileStoreService) getPropertiesFromS3Reference(opts StoreOptions) storeProperties {
	return storeProperties{
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

func (svc *FileStoreService) store(props storeProperties) error {
	if err := svc.s3.PutFile(props.Original.Key, props.Path, props.ContentType, props.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (svc *FileStoreService) createSnapshot(file model.File, props storeProperties) (model.Snapshot, error) {
	s := repo.NewSnapshot()
	s.SetID(props.SnapshotID)
	if props.ExceedsProcessingLimit {
		s.SetStatus(model.SnapshotStatusReady)
	} else {
		s.SetStatus(model.SnapshotStatusWaiting)
	}
	latestVersion, err := svc.snapshotRepo.FindLatestVersionForFile(file.GetID())
	if err != nil {
		return nil, err
	}
	s.SetVersion(latestVersion + 1)
	s.SetOriginal(&props.Original)
	if err := svc.snapshotSvc.SaveAndSync(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (svc *FileStoreService) assignSnapshotToFile(file model.File, snapshot model.Snapshot) error {
	file.SetSnapshotID(helper.ToPtr(snapshot.GetID()))
	if err := svc.fileCoreSvc.SaveAndSync(file); err != nil {
		return err
	}
	if err := svc.snapshotRepo.MapWithFile(snapshot.GetID(), file.GetID()); err != nil {
		return err
	}
	return nil
}

func (svc *FileStoreService) createTask(file model.File, userID string) (model.Task, error) {
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
	return task, nil
}

func (svc *FileStoreService) process(file model.File, snapshot model.Snapshot, props storeProperties, userID string) error {
	task, err := svc.createTask(file, userID)
	if err != nil {
		return err
	}
	snapshot.SetTaskID(helper.ToPtr(task.GetID()))
	if err := svc.snapshotSvc.SaveAndSync(snapshot); err != nil {
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
