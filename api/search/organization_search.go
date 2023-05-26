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
		orgRepo: repo.NewPostgresOrganizationRepo(),
	}
}

func (search *OrganizationSearch) Index(orgs []model.CoreOrganization) error {
	if len(orgs) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, o := range orgs {
		res = append(res, o)
	}
	if err := search.search.Index(search.index, res); err != nil {
		return err
	}
	return nil
}

func (search *OrganizationSearch) Update(orgs []model.CoreOrganization) error {
	if len(orgs) == 0 {
		return nil
	}
	var res []infra.SearchModel
	for _, o := range orgs {
		res = append(res, o)
	}
	if err := search.search.Update(search.index, res); err != nil {
		return err
	}
	return nil
}

func (search *OrganizationSearch) Delete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := search.search.Delete(search.index, ids); err != nil {
		return err
	}
	return nil
}

func (search *OrganizationSearch) Query(query string) ([]model.CoreOrganization, error) {
	hits, err := search.search.Query(search.index, query)
	if err != nil {
		return nil, err
	}
	var res []model.CoreOrganization
	for _, v := range hits {
		var b []byte
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var org repo.OrganizationEntity
		if err = json.Unmarshal(b, &org); err != nil {
			return nil, err
		}
		res = append(res, &org)
	}
	return res, nil
}
