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
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileDeleteService struct {
	fileRepo       repo.FileRepo
	fileSearch     search.FileSearch
	fileGuard      guard.FileGuard
	fileCache      cache.FileCache
	workspaceCache cache.WorkspaceCache
	taskSvc        *TaskService
	snapshotRepo   repo.SnapshotRepo
	snapshotSvc    *SnapshotService
}

func NewFileDeleteService() *FileDeleteService {
	return &FileDeleteService{
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

func (svc *FileDeleteService) DeleteOne(id string, userID string) error {
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
	return svc.delete(file)
}

func (svc *FileDeleteService) delete(file model.File) error {
	if file.GetType() == model.FileTypeFolder {
		return svc.deleteFolder(file.GetID())
	} else if file.GetType() == model.FileTypeFile {
		return svc.deleteFile(file.GetID())
	}
	return nil
}

func (svc *FileDeleteService) check(file model.File) error {
	if file.GetParentID() == nil {
		workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
		if err != nil {
			return err
		}
		return errorpkg.NewCannotDeleteWorkspaceRootError(file, workspace)
	}
	return nil
}

func (svc *FileDeleteService) createTask(file model.File, userID string) (model.Task, error) {
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

func (svc *FileDeleteService) deleteFolder(id string) error {
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

func (svc *FileDeleteService) deleteFile(id string) error {
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

func (svc *FileDeleteService) deleteSnapshots(ids []string) {
	for _, id := range ids {
		if err := svc.snapshotSvc.deleteForFile(id); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *FileDeleteService) deleteFromRepo(ids []string) {
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

func (svc *FileDeleteService) deleteFromCache(ids []string) {
	for _, id := range ids {
		if err := svc.fileCache.Delete(id); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *FileDeleteService) deleteFromSearch(ids []string) {
	if err := svc.fileSearch.Delete(ids); err != nil {
		log.GetLogger().Error(err)
	}
}

type FileDeleteManyOptions struct {
	IDs []string `json:"ids" validate:"required"`
}

type FileDeleteManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *FileDeleteService) DeleteMany(opts FileDeleteManyOptions, userID string) (*FileDeleteManyResult, error) {
	res := &FileDeleteManyResult{
		Failed:    make([]string, 0),
		Succeeded: make([]string, 0),
	}
	for _, id := range opts.IDs {
		if err := svc.DeleteOne(id, userID); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}
