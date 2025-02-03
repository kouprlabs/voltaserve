// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package search

import (
	"encoding/json"

	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type WorkspaceSearch interface {
	Index(workspaces []model.Workspace) error
	Update(workspaces []model.Workspace) error
	Delete(ids []string) error
	Query(query string, opts infra.QueryOptions) ([]model.Workspace, error)
}

func NewWorkspaceSearch() WorkspaceSearch {
	return newWorkspaceSearch()
}

type workspaceSearch struct {
	index         string
	search        infra.SearchManager
	workspaceRepo repo.WorkspaceRepo
}

type workspaceEntity struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	StorageCapacity int64   `json:"storageCapacity"`
	RootID          *string `json:"rootId"`
	OrganizationID  string  `json:"organizationId"`
	Bucket          string  `json:"bucket"`
	CreateTime      string  `json:"createTime"`
	UpdateTime      *string `json:"updateTime,omitempty"`
}

func (w workspaceEntity) GetID() string {
	return w.ID
}

func newWorkspaceSearch() *workspaceSearch {
	return &workspaceSearch{
		index:         infra.WorkspaceSearchIndex,
		search:        infra.NewSearchManager(),
		workspaceRepo: repo.NewWorkspaceRepo(),
	}
}

func (s *workspaceSearch) Index(workspaces []model.Workspace) error {
	if len(workspaces) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, w := range workspaces {
		res = append(res, s.mapEntity(w))
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *workspaceSearch) Update(workspaces []model.Workspace) error {
	if len(workspaces) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, w := range workspaces {
		res = append(res, s.mapEntity(w))
	}
	if err := s.search.Update(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *workspaceSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *workspaceSearch) Query(query string, opts infra.QueryOptions) ([]model.Workspace, error) {
	hits, err := s.search.Query(s.index, query, opts)
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

func (s *workspaceSearch) mapEntity(workspace model.Workspace) *workspaceEntity {
	entity := &workspaceEntity{
		ID:              workspace.GetID(),
		Name:            workspace.GetName(),
		StorageCapacity: workspace.GetStorageCapacity(),
		OrganizationID:  workspace.GetOrganizationID(),
		Bucket:          workspace.GetBucket(),
		CreateTime:      workspace.GetCreateTime(),
		UpdateTime:      workspace.GetUpdateTime(),
	}
	if workspace.GetRootID() != "" {
		rootID := workspace.GetRootID()
		entity.RootID = &rootID
	}
	return entity
}
