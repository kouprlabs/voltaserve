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

type OrganizationSearch struct {
	index   string
	search  *infra.SearchManager
	orgRepo repo.OrganizationRepo
}

type organizationEntity struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Members    []string `json:"members"`
	CreateTime string   `json:"createTime"`
	UpdateTime *string  `json:"updateTime,omitempty"`
}

func (o organizationEntity) GetID() string {
	return o.ID
}

func NewOrganizationSearch() *OrganizationSearch {
	return &OrganizationSearch{
		index:   infra.OrganizationSearchIndex,
		search:  infra.NewSearchManager(),
		orgRepo: repo.NewOrganizationRepo(),
	}
}

func (s *OrganizationSearch) Index(orgs []model.Organization) error {
	if len(orgs) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, o := range orgs {
		res = append(res, s.mapEntity(o))
	}
	if err := s.search.Index(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *OrganizationSearch) Update(orgs []model.Organization) error {
	if len(orgs) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, o := range orgs {
		res = append(res, s.mapEntity(o))
	}
	if err := s.search.Update(s.index, res); err != nil {
		return err
	}
	return nil
}

func (s *OrganizationSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.search.Delete(s.index, ids); err != nil {
		return err
	}
	return nil
}

func (s *OrganizationSearch) Query(query string, opts infra.QueryOptions) ([]model.Organization, error) {
	hits, err := s.search.Query(s.index, query, opts)
	if err != nil {
		return nil, err
	}
	var res []model.Organization
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		org := repo.NewOrganization()
		if err = json.Unmarshal(b, &org); err != nil {
			return nil, err
		}
		res = append(res, org)
	}
	return res, nil
}

func (s *OrganizationSearch) mapEntity(org model.Organization) *organizationEntity {
	return &organizationEntity{
		ID:         org.GetID(),
		Name:       org.GetName(),
		Members:    org.GetMembers(),
		CreateTime: org.GetCreateTime(),
		UpdateTime: org.GetUpdateTime(),
	}
}
