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

type GroupSearch struct {
	index     string
	search    *infra.SearchManager
	groupRepo repo.GroupRepo
}

func NewGroupSearch() *GroupSearch {
	return &GroupSearch{
		index:     infra.GroupSearchIndex,
		search:    infra.NewSearchManager(),
		groupRepo: repo.NewGroupRepo(),
	}
}

func (s *GroupSearch) Index(groups []model.Group) error {
	if len(groups) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, g := range groups {
		res = append(res, g)
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *GroupSearch) Update(groups []model.Group) error {
	if len(groups) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, g := range groups {
		res = append(res, g)
	}
	if err := s.search.Update(s.index, res); err != nil {
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

func (s *GroupSearch) Query(query string) ([]model.Group, error) {
	hits, err := s.search.Query(s.index, query)
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
