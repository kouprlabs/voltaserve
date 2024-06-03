package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type TaskSearch struct {
	index    string
	search   *infra.SearchManager
	taskRepo repo.TaskRepo
}

func NewTaskSearch() *TaskSearch {
	return &TaskSearch{
		index:    infra.TaskSearchIndex,
		search:   infra.NewSearchManager(),
		taskRepo: repo.NewTaskRepo(),
	}
}

func (s *TaskSearch) Index(tasks []model.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, o := range tasks {
		res = append(res, o)
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *TaskSearch) Update(orgs []model.Task) error {
	if len(orgs) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, o := range orgs {
		res = append(res, o)
	}
	if err := s.search.Update(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *TaskSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *TaskSearch) Query(query string) ([]model.Task, error) {
	hits, err := s.search.Query(s.index, query)
	if err != nil {
		return nil, err
	}
	var res []model.Task
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		org := repo.NewTask()
		if err = json.Unmarshal(b, &org); err != nil {
			return nil, err
		}
		res = append(res, org)
	}
	return res, nil
}
