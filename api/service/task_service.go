package service

import (
	"sort"
	"time"
	"voltaserve/cache"
	"voltaserve/errorpkg"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"
)

type Task struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Error           *string `json:"error,omitempty"`
	Percentage      *int    `json:"percentage,omitempty"`
	IsComplete      bool    `json:"isComplete"`
	IsIndeterminate bool    `json:"isIndeterminate"`
	UserID          string  `json:"userId"`
}

type TaskService struct {
	taskMapper *taskMapper
	taskCache  *cache.TaskCache
	taskSearch *search.TaskSearch
	taskRepo   repo.TaskRepo
}

func NewTaskService() *TaskService {
	return &TaskService{
		taskMapper: newTaskMapper(),
		taskCache:  cache.NewTaskCache(),
		taskSearch: search.NewTaskSearch(),
		taskRepo:   repo.NewTaskRepo(),
	}
}

func (svc *TaskService) Find(id string, userID string) (*Task, error) {
	task, err := svc.taskCache.Get(id)
	if err != nil {
		return nil, err
	}
	res, err := svc.taskMapper.mapOne(task)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type TaskListOptions struct {
	Query     string
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type TaskList struct {
	Data          []*Task `json:"data"`
	TotalPages    uint    `json:"totalPages"`
	TotalElements uint    `json:"totalElements"`
	Page          uint    `json:"page"`
	Size          uint    `json:"size"`
}

func (svc *TaskService) List(opts TaskListOptions, userID string) (*TaskList, error) {
	var authorized []model.Task
	if opts.Query == "" {
		ids, err := svc.taskRepo.GetIDs()
		if err != nil {
			return nil, err
		}
		authorized, err = svc.doAuthorizationByIDs(ids, userID)
		if err != nil {
			return nil, err
		}
	} else {
		orgs, err := svc.taskSearch.Query(opts.Query)
		if err != nil {
			return nil, err
		}
		authorized, err = svc.doAuthorization(orgs, userID)
		if err != nil {
			return nil, err
		}
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(authorized, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.taskMapper.mapMany(paged)
	if err != nil {
		return nil, err
	}
	return &TaskList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *TaskService) doAuthorization(data []model.Task, userID string) ([]model.Task, error) {
	var res []model.Task
	for _, t := range data {
		if t.GetUserID() == userID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (svc *TaskService) doAuthorizationByIDs(ids []string, userID string) ([]model.Task, error) {
	var res []model.Task
	for _, id := range ids {
		var t model.Task
		t, err := svc.taskCache.Get(id)
		if err != nil {
			return nil, err
		}
		if t.GetUserID() == userID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (svc *TaskService) doSorting(data []model.Task, sortBy string, sortOrder string) []model.Task {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetName() > data[j].GetName()
			} else {
				return data[i].GetName() < data[j].GetName()
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

func (svc *TaskService) GetCount(userID string) (int64, error) {
	var res int64
	var err error
	if res, err = svc.taskRepo.GetCount(userID); err != nil {
		return -1, err
	}
	return res, nil
}

func (svc *TaskService) doPagination(data []model.Task, page, size uint) ([]model.Task, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Task{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

func (svc *TaskService) Delete(id string, userID string) error {
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
	if err := svc.taskRepo.Delete(id); err != nil {
		return err
	}
	if err := svc.taskCache.Delete(task.GetID()); err != nil {
		return err
	}
	if err := svc.taskSearch.Delete([]string{task.GetID()}); err != nil {
		return err
	}
	return nil
}

type taskMapper struct {
	groupCache *cache.TaskCache
}

func newTaskMapper() *taskMapper {
	return &taskMapper{
		groupCache: cache.NewTaskCache(),
	}
}

func (mp *taskMapper) mapOne(m model.Task) (*Task, error) {
	return &Task{
		ID:              m.GetID(),
		Name:            m.GetName(),
		Error:           m.GetError(),
		Percentage:      m.GetPercentage(),
		IsIndeterminate: m.GetIsIndeterminate(),
		UserID:          m.GetUserID(),
	}, nil
}

func (mp *taskMapper) mapMany(orgs []model.Task) ([]*Task, error) {
	res := make([]*Task, 0)
	for _, task := range orgs {
		t, err := mp.mapOne(task)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}