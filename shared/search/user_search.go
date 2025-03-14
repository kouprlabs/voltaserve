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

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
)

type UserSearch struct {
	index  string
	search infra.SearchManager
}

func NewUserSearch(search config.SearchConfig, environment config.EnvironmentConfig) *UserSearch {
	return &UserSearch{
		index:  infra.UserSearchIndex,
		search: infra.NewSearchManager(search, environment),
	}
}

func (s *UserSearch) Query(query string, opts infra.SearchQueryOptions) ([]model.User, error) {
	hits, err := s.search.Query(s.index, query, opts)
	if err != nil {
		return nil, err
	}
	res := make([]model.User, 0)
	for _, v := range hits {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		user := repo.NewUserModel()
		if err := json.Unmarshal(b, &user); err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}
