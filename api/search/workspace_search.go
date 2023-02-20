package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type WorkspaceSearch struct {
	index         string
	search        *infra.SearchManager
	workspaceRepo *repo.WorkspaceRepo
}

func NewWorkspaceSearch() *WorkspaceSearch {
	return &WorkspaceSearch{
		index:         infra.WorkspaceSearchIndex,
		search:        infra.NewSearchManager(),
		workspaceRepo: repo.NewWorkspaceRepo(),
	}
}

func (search *WorkspaceSearch) Index(workspaces []model.WorkspaceModel) error {
	if len(workspaces) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, w := range workspaces {
		res = append(res, w)
	}
	if err := search.search.Index(search.index, res); err != nil {
		return err
	}
	return nil
}

func (search *WorkspaceSearch) Update(workspaces []model.WorkspaceModel) error {
	if len(workspaces) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, w := range workspaces {
		res = append(res, w)
	}
	if err := search.search.Update(search.index, res); err != nil {
		return err
	}
	return nil
}

func (search *WorkspaceSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := search.search.Delete(search.index, ids); err != nil {
		return err
	}
	return nil
}

func (search *WorkspaceSearch) Query(query string) ([]model.WorkspaceModel, error) {
	hits, err := search.search.Query(search.index, query)
	if err != nil {
		return nil, err
	}
	var res []model.WorkspaceModel
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var workspace repo.WorkspaceEntity
		if err = json.Unmarshal(b, &workspace); err != nil {
			return nil, err
		}
		res = append(res, &workspace)
	}
	return res, nil
}
