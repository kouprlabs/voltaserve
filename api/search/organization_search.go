package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type OrganizationSearch struct {
	index   string
	search  *infra.SearchManager
	orgRepo repo.OrganizationRepo
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
		res = append(res, o)
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
		res = append(res, o)
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

func (s *OrganizationSearch) Query(query string) ([]model.Organization, error) {
	hits, err := s.search.Query(s.index, query)
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
