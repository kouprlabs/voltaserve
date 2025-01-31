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
	"time"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
)

type TaskService struct {
	taskMapper    *taskMapper
	taskCache     cache.TaskCache
	taskSearch    search.TaskSearch
	taskRepo      repo.TaskRepo
	snapshotRepo  repo.SnapshotRepo
	snapshotCache cache.SnapshotCache
	fileRepo      repo.FileRepo
	fileCache     cache.FileCache
}

func NewTaskService() *TaskService {
	return &TaskService{
		taskMapper:    newTaskMapper(),
		taskCache:     cache.NewTaskCache(),
		taskSearch:    search.NewTaskSearch(),
		taskRepo:      repo.NewTaskRepo(),
		snapshotRepo:  repo.NewSnapshotRepo(),
		snapshotCache: cache.NewSnapshotCache(),
		fileRepo:      repo.NewFileRepo(),
		fileCache:     cache.NewFileCache(),
	}
}

type TaskCreateOptions struct {
	Name            string            `json:"name"`
	Error           *string           `json:"error,omitempty"`
	Percentage      *int              `json:"percentage,omitempty"`
	IsIndeterminate bool              `json:"isIndeterminate"`
	UserID          string            `json:"userId"`
	Status          string            `json:"status"`
	Payload         map[string]string `json:"payload,omitempty"`
}

func (svc *TaskService) Create(opts TaskCreateOptions) (*Task, error) {
	task, err := svc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            opts.Name,
		Error:           opts.Error,
		Percentage:      opts.Percentage,
		IsIndeterminate: opts.IsIndeterminate,
		UserID:          opts.UserID,
		Payload:         opts.Payload,
	})
	if err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type TaskPatchOptions struct {
	Fields          []string          `json:"fields"`
	Name            *string           `json:"name"`
	Error           *string           `json:"error"`
	Percentage      *int              `json:"percentage"`
	IsIndeterminate *bool             `json:"isIndeterminate"`
	UserID          *string           `json:"userId"`
	Status          *string           `json:"status"`
	Payload         map[string]string `json:"payload"`
}

const (
	TaskFieldName            = "name"
	TaskFieldError           = "error"
	TaskFieldPercentage      = "percentage"
	TaskFieldIsIndeterminate = "isIndeterminate"
	TaskFieldUserID          = "userId"
	TaskFieldStatus          = "status"
	TaskFieldPayload         = "payload"
)

func (svc *TaskService) Patch(id string, opts TaskPatchOptions) (*Task, error) {
	task, err := svc.taskCache.Get(id)
	if err != nil {
		return nil, err
	}
	if slices.Contains(opts.Fields, TaskFieldName) {
		task.SetName(*opts.Name)
	}
	if slices.Contains(opts.Fields, TaskFieldError) {
		task.SetError(opts.Error)
	}
	if slices.Contains(opts.Fields, TaskFieldPercentage) {
		task.SetPercentage(opts.Percentage)
	}
	if slices.Contains(opts.Fields, TaskFieldIsIndeterminate) {
		task.SetIsIndeterminate(true)
	}
	if slices.Contains(opts.Fields, TaskFieldUserID) {
		task.SetUserID(*opts.UserID)
	}
	if slices.Contains(opts.Fields, TaskFieldStatus) {
		task.SetStatus(*opts.Status)
	}
	if slices.Contains(opts.Fields, TaskFieldPayload) {
		task.SetPayload(opts.Payload)
	}
	if err := svc.saveAndSync(task); err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *TaskService) Find(id string, userID string) (*Task, error) {
	task, err := svc.taskCache.Get(id)
	if err != nil {
		return nil, err
	}
	if task.GetUserID() != userID {
		return nil, errorpkg.NewTaskNotFoundError(nil)
	}
	res, err := svc.taskMapper.mapOne(task)
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

func (svc *TaskService) List(opts TaskListOptions, userID string) (*TaskList, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.taskMapper.mapMany(paged)
	if err != nil {
		return nil, err
	}
	return &TaskList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *TaskService) Probe(opts TaskListOptions, userID string) (*TaskProbe, error) {
	all, err := svc.findAll(opts, userID)
	if err != nil {
		return nil, err
	}
	totalElements := uint64(len(all))
	return &TaskProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
	}, nil
}

func (svc *TaskService) findAll(opts TaskListOptions, userID string) ([]model.Task, error) {
	var res []model.Task
	if opts.Query == "" {
		ids, err := svc.taskRepo.FindIDs(userID)
		if err != nil {
			return nil, err
		}
		res, err = svc.authorizeIDs(ids, userID)
		if err != nil {
			return nil, err
		}
	} else {
		hits, err := svc.taskSearch.Query(opts.Query, infra.QueryOptions{})
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
	}
	return res, nil
}

func (svc *TaskService) Count(userID string) (*int64, error) {
	var res int64
	var err error
	if res, err = svc.taskRepo.CountByEmail(userID); err != nil {
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
	if !task.HasError() {
		return errorpkg.NewTaskIsRunningError(nil)
	}
	return svc.deleteAndSync(id)
}

type TaskDismissAllResult struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *TaskService) DismissAll(userID string) (*TaskDismissAllResult, error) {
	ids, err := svc.taskRepo.FindIDs(userID)
	if err != nil {
		return nil, err
	}
	authorized, err := svc.authorizeIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	res := TaskDismissAllResult{
		Succeeded: make([]string, 0),
		Failed:    make([]string, 0),
	}
	for _, t := range authorized {
		if t.HasError() {
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
		var t model.Task
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
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
			}
		})
		return data
	} else if sortBy == SortByStatus {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetStatus() == model.TaskStatusRunning && data[j].GetStatus() != model.TaskStatusRunning
			} else {
				return data[i].GetStatus() == model.TaskStatusWaiting && data[j].GetStatus() != model.TaskStatusWaiting
			}
		})
		return data
	} else if sortBy == SortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a, _ := time.Parse(time.RFC3339, data[i].GetCreateTime())
			b, _ := time.Parse(time.RFC3339, data[j].GetCreateTime())
			if sortOrder == SortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == SortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a, _ := time.Parse(time.RFC3339, *data[i].GetUpdateTime())
				b, _ := time.Parse(time.RFC3339, *data[j].GetUpdateTime())
				if sortOrder == SortOrderDesc {
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
			log.GetLogger().Error(err)
		}
		if _, err = svc.snapshotCache.Refresh(snapshot.GetID()); err != nil {
			log.GetLogger().Error(err)
		}
		var filesIDs []string
		filesIDs, err = svc.fileRepo.FindIDsBySnapshot(snapshot.GetID())
		if err == nil {
			for _, fileID := range filesIDs {
				if _, err = svc.fileCache.Refresh(fileID); err != nil {
					log.GetLogger().Error(err)
				}
			}
		} else {
			log.GetLogger().Error(err)
		}
	}
	// Proceed with deleting the task
	if err = svc.taskRepo.Delete(id); err != nil {
		return err
	}
	if err = svc.taskCache.Delete(id); err != nil {
		return err
	}
	if err = svc.taskSearch.Delete([]string{id}); err != nil {
		return err
	}
	return nil
}
