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

type GroupSearch interface {
	Index(groups []model.Group) error
	Update(groups []model.Group) error
	Delete(ids []string) error
	Query(query string, opts infra.QueryOptions) ([]model.Group, error)
}

func NewGroupSearch() GroupSearch {
	return newGroupSearch()
}

type groupSearch struct {
	index     string
	search    infra.SearchManager
	groupRepo repo.GroupRepo
}

type groupEntity struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	OrganizationID string   `json:"organizationId"`
	Members        []string `json:"members"`
	CreateTime     string   `json:"createTime"`
	UpdateTime     *string  `json:"updateTime"`
}

func (g groupEntity) GetID() string {
	return g.ID
}

func newGroupSearch() *groupSearch {
	return &groupSearch{
		index:     infra.GroupSearchIndex,
		search:    infra.NewSearchManager(),
		groupRepo: repo.NewGroupRepo(),
	}
}

func (s *groupSearch) Index(groups []model.Group) error {
	if len(groups) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, g := range groups {
		res = append(res, s.mapEntity(g))
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *groupSearch) Update(groups []model.Group) error {
	if len(groups) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, g := range groups {
		res = append(res, s.mapEntity(g))
	}
	if err := s.search.Update(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *groupSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *groupSearch) Query(query string, opts infra.QueryOptions) ([]model.Group, error) {
	hits, err := s.search.Query(s.index, query, opts)
	if err != nil {
		return nil, err
	}
	var res []model.Group
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		group := repo.NewGroup()
		if err = json.Unmarshal(b, &group); err != nil {
			return nil, err
		}
		res = append(res, group)
	}
	return res, nil
}

func (s *groupSearch) mapEntity(group model.Group) *groupEntity {
	return &groupEntity{
		ID:             group.GetID(),
		Name:           group.GetName(),
		OrganizationID: group.GetOrganizationID(),
		Members:        group.GetMembers(),
		CreateTime:     group.GetCreateTime(),
		UpdateTime:     group.GetUpdateTime(),
	}
}
