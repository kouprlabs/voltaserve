package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type GroupSearch struct {
	index     string
	search    *infra.SearchManager
	groupRepo repo.CoreGroupRepo
}

func NewGroupSearch() *GroupSearch {
	return &GroupSearch{
		index:     infra.GroupSearchIndex,
		search:    infra.NewSearchManager(),
		groupRepo: repo.NewPostgresGroupRepo(),
	}
}

func (search *GroupSearch) Index(groups []model.GroupModel) error {
	if len(groups) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, g := range groups {
		res = append(res, g)
	}
	if err := search.search.Index(search.index, res); err != nil {
		return err
	}
	return nil
}

func (search *GroupSearch) Update(groups []model.GroupModel) error {
	if len(groups) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, g := range groups {
		res = append(res, g)
	}
	if err := search.search.Update(search.index, res); err != nil {
		return err
	}
	return nil
}

func (search *GroupSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := search.search.Delete(search.index, ids); err != nil {
		return err
	}
	return nil
}

func (search *GroupSearch) Query(query string) ([]model.GroupModel, error) {
	hits, err := search.search.Query(search.index, query)
	if err != nil {
		return nil, err
	}
	var res []model.GroupModel
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var group repo.PostgresGroup
		if err = json.Unmarshal(b, &group); err != nil {
			return nil, err
		}
		res = append(res, &group)
	}
	return res, nil
}
