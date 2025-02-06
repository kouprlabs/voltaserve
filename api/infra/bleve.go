// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve/v2"
	bleve_query "github.com/blevesearch/bleve/v2/search/query"
	bleve_index "github.com/blevesearch/bleve_index_api"
)

type bleveSearchManager struct {
	indexes map[string]bleve.Index
}

func newBleveSearchManager() SearchManager {
	manager := &bleveSearchManager{
		indexes: make(map[string]bleve.Index),
	}
	manager.createIndex(FileSearchIndex)
	manager.createIndex(GroupSearchIndex)
	manager.createIndex(WorkspaceSearchIndex)
	manager.createIndex(OrganizationSearchIndex)
	manager.createIndex(UserSearchIndex)
	manager.createIndex(TaskSearchIndex)
	return manager
}

func (mgr *bleveSearchManager) Query(indexName string, query string, opts QueryOptions) ([]interface{}, error) {
	index, ok := mgr.indexes[indexName]
	if !ok {
		return nil, errors.New("index not found")
	}
	var searchRequest *bleve.SearchRequest
	var err error
	if opts.Filter == nil {
		searchRequest = bleve.NewSearchRequestOptions(
			bleve.NewQueryStringQuery(query),
			int(opts.Limit), 0, false,
		)
	} else {
		filterQueries := mgr.buildFilter(opts.Filter)
		conjunctionQuery := bleve.NewConjunctionQuery(bleve.NewQueryStringQuery(query))
		for _, v := range filterQueries {
			conjunctionQuery.AddQuery(v)
		}
		searchRequest = bleve.NewSearchRequestOptions(conjunctionQuery, int(opts.Limit), 0, false)
	}
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	res := make([]interface{}, len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
		doc, err := index.Document(hit.ID)
		if err != nil {
			return nil, err
		}
		raw := make(map[string]interface{})
		doc.VisitFields(func(field bleve_index.Field) {
			fieldName := field.Name()
			fieldValue := field.Value()
			switch field.(type) {
			case bleve_index.TextField:
				raw[fieldName] = string(fieldValue)
			case bleve_index.NumericField:
				num, err := strconv.ParseFloat(string(fieldValue), 64)
				if err == nil {
					raw[fieldName] = num
				}
			case bleve_index.BooleanField:
				boolVal, err := strconv.ParseBool(string(fieldValue))
				if err == nil {
					raw[fieldName] = boolVal
				}
			default:
				raw[fieldName] = string(fieldValue)
			}
		})
		res[i] = raw
	}
	return res, nil
}

func (mgr *bleveSearchManager) Index(indexName string, models []SearchModel) error {
	index, ok := mgr.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	batch := index.NewBatch()
	for _, model := range models {
		if err := batch.Index(model.GetID(), model); err != nil {
			return err
		}
	}
	return index.Batch(batch)
}

func (mgr *bleveSearchManager) Update(indexName string, models []SearchModel) error {
	return mgr.Index(indexName, models)
}

func (mgr *bleveSearchManager) Delete(indexName string, ids []string) error {
	index, ok := mgr.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	batch := index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	return index.Batch(batch)
}

func (mgr *bleveSearchManager) createIndex(indexName string) {
	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		panic(err)
	}
	mgr.indexes[indexName] = index
}

func (mgr *bleveSearchManager) buildFilter(filter interface{}) []*bleve_query.MatchQuery {
	res := make([]*bleve_query.MatchQuery, 0)
	expression, ok := filter.(string)
	if !ok {
		return nil
	}
	expression = regexp.MustCompile(`\s+`).ReplaceAllString(expression, "")
	parts := strings.Split(expression, "AND")
	for _, part := range parts {
		condition := strings.Split(part, "=")
		matchQuery := bleve.NewMatchQuery(condition[1])
		matchQuery.SetField(condition[0])
		res = append(res, matchQuery)
	}
	return res
}
