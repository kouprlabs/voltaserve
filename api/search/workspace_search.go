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
	workspaceRepo repo.WorkspaceRepo
}

func NewWorkspaceSearch() *WorkspaceSearch {
	return &WorkspaceSearch{
		index:         infra.WorkspaceSearchIndex,
		search:        infra.NewSearchManager(),
		workspaceRepo: repo.NewWorkspaceRepo(),
	}
}

func (s *WorkspaceSearch) Index(workspaces []model.Workspace) error {
	if len(workspaces) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, w := range workspaces {
		res = append(res, w)
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceSearch) Update(workspaces []model.Workspace) error {
	if len(workspaces) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, w := range workspaces {
		res = append(res, w)
	}
	if err := s.search.Update(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceSearch) Query(query string) ([]model.Workspace, error) {
	hits, err := s.search.Query(s.index, query)
	if err != nil {
		return nil, err
	}
	var res []model.Workspace
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		workspace := repo.NewWorkspace()
		if err = json.Unmarshal(b, &workspace); err != nil {
			return nil, err
		}
		res = append(res, workspace)
	}
	return res, nil
}
