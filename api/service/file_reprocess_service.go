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
	"slices"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/conversion_client"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type FileReprocessService struct {
	fileCache      *cache.FileCache
	fileRepo       repo.FileRepo
	fileGuard      *guard.FileGuard
	snapshotCache  *cache.SnapshotCache
	snapshotSvc    *SnapshotService
	taskCache      *cache.TaskCache
	taskSvc        *TaskService
	pipelineClient *conversion_client.PipelineClient
}

func NewFileReprocessService() *FileReprocessService {
	return &FileReprocessService{
		fileCache:      cache.NewFileCache(),
		fileRepo:       repo.NewFileRepo(),
		fileGuard:      guard.NewFileGuard(),
		snapshotCache:  cache.NewSnapshotCache(),
		snapshotSvc:    NewSnapshotService(),
		taskCache:      cache.NewTaskCache(),
		taskSvc:        NewTaskService(),
		pipelineClient: conversion_client.NewPipelineClient(),
	}
}

type FileReprocessResponse struct {
	Accepted []string `json:"accepted"`
	Rejected []string `json:"rejected"`
}

func (r *FileReprocessResponse) AppendAccepted(id string) {
	if !slices.Contains(r.Accepted, id) {
		r.Accepted = append(r.Accepted, id)
	}
}

func (r *FileReprocessResponse) AppendRejected(id string) {
	if !slices.Contains(r.Rejected, id) {
		r.Rejected = append(r.Rejected, id)
	}
}

func (svc *FileReprocessService) Reprocess(id string, userID string) (*FileReprocessResponse, error) {
	resp := &FileReprocessResponse{
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
		if svc.reprocess(leaf, userID) {
			resp.AppendAccepted(leaf.GetID())
		} else {
			resp.AppendRejected(leaf.GetID())
		}
	}
	return resp, nil
}

func (svc *FileReprocessService) reprocess(leaf model.File, userID string) bool {
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

func (svc *FileReprocessService) check(file model.File, snapshot model.Snapshot) bool {
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

func (svc *FileReprocessService) getTree(file model.File, userID string) ([]model.File, error) {
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

func (svc *FileReprocessService) createTask(file model.File, userID string) (model.Task, error) {
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

func (svc *FileReprocessService) runPipeline(file model.File, snapshot model.Snapshot, userID string) error {
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
