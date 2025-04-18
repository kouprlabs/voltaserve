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
	"errors"
	"slices"
	"sort"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/mapper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
	"github.com/kouprlabs/voltaserve/shared/search"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type TaskService struct {
	taskMapper    *mapper.TaskMapper
	taskCache     *cache.TaskCache
	taskSearch    *search.TaskSearch
	taskRepo      *repo.TaskRepo
	snapshotRepo  *repo.SnapshotRepo
	snapshotCache *cache.SnapshotCache
	fileRepo      *repo.FileRepo
	fileCache     *cache.FileCache
}

func NewTaskService() *TaskService {
	return &TaskService{
		taskMapper: mapper.NewTaskMapper(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		taskCache: cache.NewTaskCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
		taskSearch: search.NewTaskSearch(
			config.GetConfig().Search,
			config.GetConfig().Environment,
		),
		taskRepo: repo.NewTaskRepo(
			config.GetConfig().Postgres,
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
		fileRepo: repo.NewFileRepo(
			config.GetConfig().Postgres,
			config.GetConfig().Environment,
		),
		fileCache: cache.NewFileCache(
			config.GetConfig().Postgres,
			config.GetConfig().Redis,
			config.GetConfig().Environment,
		),
	}
}

func (svc *TaskService) Create(opts dto.TaskCreateOptions) (*dto.Task, error) {
	task, err := svc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            opts.Name,
		Status:          opts.Status,
		Error:           opts.Error,
		Percentage:      opts.Percentage,
		IsIndeterminate: opts.IsIndeterminate,
		UserID:          opts.UserID,
		Payload:         opts.Payload,
	})
	if err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.Map(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *TaskService) Patch(id string, opts dto.TaskPatchOptions) (*dto.Task, error) {
	task, err := svc.taskCache.Get(id)
	if err != nil {
		return nil, err
	}
	if slices.Contains(opts.Fields, model.TaskFieldName) {
		task.SetName(*opts.Name)
	}
	if slices.Contains(opts.Fields, model.TaskFieldError) {
		task.SetError(opts.Error)
	}
	if slices.Contains(opts.Fields, model.TaskFieldPercentage) {
		task.SetPercentage(opts.Percentage)
	}
	if slices.Contains(opts.Fields, model.TaskFieldIsIndeterminate) {
		task.SetIsIndeterminate(true)
	}
	if slices.Contains(opts.Fields, model.TaskFieldUserID) {
		task.SetUserID(*opts.UserID)
	}
	if slices.Contains(opts.Fields, model.TaskFieldStatus) {
		task.SetStatus(*opts.Status)
	}
	if slices.Contains(opts.Fields, model.TaskFieldPayload) {
		task.SetPayload(opts.Payload)
	}
	if err := svc.saveAndSync(task); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.Map(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *TaskService) Find(id string, userID string) (*dto.Task, error) {
	task, err := svc.taskCache.Get(id)
	if err != nil {
		return nil, err
	}
	if task.GetUserID() != userID {
		return nil, errorpkg.NewTaskNotFoundError(nil)
	}
	res, err := svc.taskMapper.Map(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type TaskListOptions struct {
	Query     string
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
}

func (svc *TaskService) List(opts TaskListOptions, userID string) (*dto.TaskList, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.TaskSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.TaskSortOrderAsc
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.taskMapper.MapMany(paged)
	if err != nil {
		return nil, err
	}
	return &dto.TaskList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *TaskService) Probe(opts TaskListOptions, userID string) (*dto.TaskProbe, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &dto.TaskProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *TaskService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.TaskSortByName ||
		value == dto.TaskSortByStatus ||
		value == dto.TaskSortByDateCreated ||
		value == dto.TaskSortByDateModified
}

func (svc *TaskService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.TaskSortOrderAsc || value == dto.TaskSortOrderDesc
}

func (svc *TaskService) Count(userID string) (*int64, error) {
	var res int64
	var err error
	if res, err = svc.taskRepo.CountByUserID(userID); err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *TaskService) Dismiss(id string, userID string) error {
	task, err := svc.taskCache.Get(id)
	if err != nil {
		return err
	}
	if task.GetUserID() != userID {
		return errorpkg.NewTaskBelongsToAnotherUserError(nil)
	}
	if task.GetStatus() != model.TaskStatusSuccess && task.GetStatus() != model.TaskStatusError {
		return errorpkg.NewTaskIsRunningError(nil)
	}
	return svc.deleteAndSync(id)
}

func (svc *TaskService) DismissAll(userID string) (*dto.TaskDismissAllResult, error) {
	ids, err := svc.taskRepo.FindIDs(userID)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.authorizeIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	res := dto.TaskDismissAllResult{
		Succeeded: make([]string, 0),
		Failed:    make([]string, 0),
	}
	for _, t := range authorized {
		if t.GetStatus() == model.TaskStatusSuccess || t.GetStatus() == model.TaskStatusError {
			if err := svc.deleteAndSync(t.GetID()); err != nil {
				res.Failed = append(res.Failed, t.GetID())
			} else {
				res.Succeeded = append(res.Succeeded, t.GetID())
			}
		}
	}
	return &res, nil
}

func (svc *TaskService) Delete(id string) error {
	return svc.deleteAndSync(id)
}

func (svc *TaskService) findAll(opts TaskListOptions, userID string) ([]model.Task, error) {
	var res []model.Task
	var err error
	if opts.Query == "" {
		res, err = svc.load(userID)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = svc.search(opts, userID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (svc *TaskService) load(userID string) ([]model.Task, error) {
	var res []model.Task
	ids, err := svc.taskRepo.FindIDs(userID)
	if err != nil {
		return nil, err
	}
	res, err = svc.authorizeIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *TaskService) search(opts TaskListOptions, userID string) ([]model.Task, error) {
	var res []model.Task
	count, err := svc.taskRepo.Count()
	if err != nil {
		return nil, err
	}
	hits, err := svc.taskSearch.Query(opts.Query, infra.SearchQueryOptions{Limit: count})
	if err != nil {
		return nil, err
	}
	var tasks []model.Task
	for _, hit := range hits {
		task, err := svc.taskCache.Get(hit.GetID())
		if err != nil {
			var e *errorpkg.ErrorResponse
			// We don't want to break if the search engine contains tasks that shouldn't be there
			if errors.As(err, &e) && e.Code == errorpkg.NewTaskNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		tasks = append(tasks, task)
	}
	res, err = svc.authorize(tasks, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *TaskService) authorize(data []model.Task, userID string) ([]model.Task, error) {
	var res []model.Task
	for _, t := range data {
		if t.GetUserID() == userID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (svc *TaskService) authorizeIDs(ids []string, userID string) ([]model.Task, error) {
	var res []model.Task
	for _, id := range ids {
		t, err := svc.taskCache.Get(id)
		if err != nil {
			var e *errorpkg.ErrorResponse
			if errors.As(err, &e) && e.Code == errorpkg.NewTaskNotFoundError(nil).Code {
				continue
			} else {
				return nil, err
			}
		}
		if t.GetUserID() == userID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (svc *TaskService) sort(data []model.Task, sortBy string, sortOrder string) []model.Task {
	if sortBy == dto.TaskSortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.TaskSortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == dto.TaskSortByStatus {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.TaskSortOrderDesc {
				return data[i].GetStatus() == model.TaskStatusRunning && data[j].GetStatus() != model.TaskStatusRunning
			} else {
				return data[i].GetStatus() == model.TaskStatusWaiting && data[j].GetStatus() != model.TaskStatusWaiting
			}
		})
		return data
	} else if sortBy == dto.TaskSortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a := helper.StringToTime(data[i].GetCreateTime())
			b := helper.StringToTime(data[j].GetCreateTime())
			if sortOrder == dto.TaskSortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == dto.TaskSortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a := helper.StringToTime(*data[i].GetUpdateTime())
				b := helper.StringToTime(*data[j].GetUpdateTime())
				if sortOrder == dto.TaskSortOrderDesc {
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
	return data
}

func (svc *TaskService) paginate(data []model.Task, page, size uint64) (pageData []model.Task, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Task{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
}

func (svc *TaskService) insertAndSync(opts repo.TaskInsertOptions) (model.Task, error) {
	task, err := svc.taskRepo.Insert(opts)
	if err != nil {
		return nil, err
	}
	if err := svc.taskCache.Set(task); err != nil {
		return nil, err
	}
	if err := svc.taskSearch.Index([]model.Task{task}); err != nil {
		return nil, err
	}
	return task, nil
}

func (svc *TaskService) saveAndSync(task model.Task) error {
	if err := svc.taskRepo.Save(task); err != nil {
		return nil
	}
	if err := svc.taskCache.Set(task); err != nil {
		return nil
	}
	if err := svc.taskSearch.Update([]model.Task{task}); err != nil {
		return nil
	}
	return nil
}

func (svc *TaskService) deleteAndSync(id string) error {
	snapshots, err := svc.snapshotRepo.FindAllForTask(id)
	if err != nil {
		return err
	}
	// Clear task ID field from all snapshots and files in both repo and cache
	for _, snapshot := range snapshots {
		snapshot.SetTaskID(nil)
		if err = svc.snapshotRepo.Save(snapshot); err != nil {
			logger.GetLogger().Error(err)
		}
		if _, err = svc.snapshotCache.Refresh(snapshot.GetID()); err != nil {
			logger.GetLogger().Error(err)
		}
		var filesIDs []string
		filesIDs, err = svc.fileRepo.FindIDsBySnapshot(snapshot.GetID())
		if err == nil {
			for _, fileID := range filesIDs {
				if _, err = svc.fileCache.Refresh(fileID); err != nil {
					logger.GetLogger().Error(err)
				}
			}
		} else {
			logger.GetLogger().Error(err)
		}
	}
	// Proceed with deleting the task
	return svc.delete(id)
}

func (svc *TaskService) delete(id string) error {
	if err := svc.taskRepo.Delete(id); err != nil {
		return err
	}
	if err := svc.taskCache.Delete(id); err != nil {
		return err
	}
	if err := svc.taskSearch.Delete([]string{id}); err != nil {
		return err
	}
	return nil
}
