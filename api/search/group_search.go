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

	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type GroupSearch struct {
	index     string
	search    infra.SearchManager
	groupRepo *repo.GroupRepo
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

func NewGroupSearch() *GroupSearch {
	return &GroupSearch{
		index:     infra.GroupSearchIndex,
		search:    infra.NewSearchManager(config.GetConfig().Search, config.GetConfig().Environment),
		groupRepo: repo.NewGroupRepo(),
	}
}

func (s *GroupSearch) Index(groups []model.Group) error {
	if len(groups) == 0 {
		return nil
	}
	var models []infra.SearchModel
	for _, g := range groups {
		models = append(models, s.mapEntity(g))
	}
	if err := s.search.Index(s.index, models); err != nil {
		return err
	}
	return nil
}

func (s *GroupSearch) Update(groups []model.Group) error {
	if len(groups) == 0 {
		return nil
	}
	var models []infra.SearchModel
	for _, g := range groups {
		models = append(models, s.mapEntity(g))
	}
	if err := s.search.Update(s.index, models); err != nil {
		return err
	}
	return nil
}

func (s *GroupSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *GroupSearch) Query(query string, opts infra.SearchQueryOptions) ([]model.Group, error) {
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
		group := repo.NewGroupModel()
		if err = json.Unmarshal(b, &group); err != nil {
			return nil, err
		}
		res = append(res, group)
	}
	return res, nil
}

func (s *GroupSearch) mapEntity(group model.Group) *groupEntity {
	return &groupEntity{
		ID:             group.GetID(),
		Name:           group.GetName(),
		OrganizationID: group.GetOrganizationID(),
		Members:        group.GetMembers(),
		CreateTime:     group.GetCreateTime(),
		UpdateTime:     group.GetUpdateTime(),
	}
}
