package search

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
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

func (svc *UserSearch) Query(query string) ([]model.User, error) {
	hits, err := svc.search.Query(svc.index, query)
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
