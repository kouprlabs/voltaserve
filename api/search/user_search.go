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

type UserSearch struct {
	index  string
	search *infra.SearchManager
}

func NewUserSearch() *UserSearch {
	return &UserSearch{
		index:  infra.UserSearchIndex,
		search: infra.NewSearchManager(),
	}
}

func (s *UserSearch) Query(query string) ([]model.User, error) {
	hits, err := s.search.Query(s.index, query)
	if err != nil {
		return nil, err
	}
	res := []model.User{}
	for _, v := range hits {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		user := repo.NewUser()
		if err := json.Unmarshal(b, &user); err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}
