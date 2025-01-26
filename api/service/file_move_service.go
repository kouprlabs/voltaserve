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
	"time"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type FileMoveService struct {
	fileRepo    repo.FileRepo
	fileSearch  *search.FileSearch
	fileCache   *cache.FileCache
	fileGuard   *guard.FileGuard
	fileMapper  *FileMapper
	fileCoreSvc *FileCoreService
	taskSvc     *TaskService
}

func NewFileMoveService() *FileMoveService {
	return &FileMoveService{
		fileRepo:    repo.NewFileRepo(),
		fileSearch:  search.NewFileSearch(),
		fileCache:   cache.NewFileCache(),
		fileGuard:   guard.NewFileGuard(),
		fileMapper:  NewFileMapper(),
		fileCoreSvc: NewFileCoreService(),
		taskSvc:     NewTaskService(),
	}
}

func (svc *FileMoveService) MoveOne(sourceID string, targetID string, userID string) (*File, error) {
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
	return svc.move(source, target, userID)
}

func (svc *FileMoveService) move(source model.File, target model.File, userID string) (*File, error) {
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

func (svc *FileMoveService) createTask(file model.File, userID string) (model.Task, error) {
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
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
	return task, nil
}

func (svc *FileMoveService) check(source model.File, target model.File, userID string) error {
	if source.GetParentID() != nil {
		existing, err := svc.fileCoreSvc.GetChildWithName(target.GetID(), source.GetName())
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

func (svc *FileMoveService) refreshUpdateAndCreateTime(source model.File, target model.File) error {
	now := time.Now().UTC().Format(time.RFC3339)
	source.SetUpdateTime(&now)
	if err := svc.fileRepo.Save(source); err != nil {
		return err
	}
	if err := svc.fileCoreSvc.Sync(source); err != nil {
		return err
	}
	target.SetUpdateTime(&now)
	if err := svc.fileRepo.Save(target); err != nil {
		return err
	}
	if err := svc.fileCoreSvc.Sync(target); err != nil {
		return err
	}
	return nil
}

type FileMoveManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileMoveManyResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *FileMoveService) MoveMany(opts FileMoveManyOptions, userID string) (*FileMoveManyResult, error) {
	res := &FileMoveManyResult{
		Failed:    make([]string, 0),
		Succeeded: make([]string, 0),
	}
	for _, id := range opts.SourceIDs {
		if _, err := svc.MoveOne(id, opts.TargetID, userID); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}
