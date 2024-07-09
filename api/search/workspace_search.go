// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package search

import (
	"encoding/json"

	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
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
