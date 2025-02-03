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
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/blevesearch/bleve/v2"
)

type bleveSearchManager struct {
	indexes map[string]bleve.Index
}

func NewBleveSearchManager() SearchManager {
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

func (mgr *bleveSearchManager) Query(index string, query string, opts QueryOptions) ([]interface{}, error) {
	indexInstance, ok := mgr.indexes[index]
	if !ok {
		return nil, errors.New("index not found")
	}
	queryString, err := parseFilter(opts.Filter)
	if err != nil {
		return nil, err
	}
	bleveQuery := bleve.NewQueryStringQuery(query + " " + queryString)
	searchRequest := bleve.NewSearchRequestOptions(bleveQuery, int(opts.Limit), 0, false)
	searchResult, err := indexInstance.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	results := make([]interface{}, len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
		results[i] = hit
	}
	return results, nil
}

func (mgr *bleveSearchManager) Index(index string, models []SearchModel) error {
	indexInstance, ok := mgr.indexes[index]
	if !ok {
		return errors.New("index not found")
	}
	batch := indexInstance.NewBatch()
	for _, model := range models {
		err := batch.Index(model.GetID(), model)
		if err != nil {
			return err
		}
	}
	return indexInstance.Batch(batch)
}

func (mgr *bleveSearchManager) Update(index string, models []SearchModel) error {
	return mgr.Index(index, models)
}

func (mgr *bleveSearchManager) Delete(index string, ids []string) error {
	indexInstance, ok := mgr.indexes[index]
	if !ok {
		return errors.New("index not found")
	}
	batch := indexInstance.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	return indexInstance.Batch(batch)
}

func (mgr *bleveSearchManager) createIndex(indexName string) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexName+".bleve", mapping)
	if err != nil {
		panic(err)
	}
	mgr.indexes[indexName] = index
}

func parseFilter(filter interface{}) (string, error) {
	filterStr, ok := filter.(string)
	if !ok {
		return "", errors.New("filter must be a string")
	}
	expression, err := govaluate.NewEvaluableExpression(filterStr)
	if err != nil {
		return "", err
	}
	tokens := expression.Tokens()
	var queryParts []string
	for _, token := range tokens {
		//nolint:exhaustive
		switch token.Kind {
		case govaluate.VARIABLE:
			queryParts = append(queryParts, token.Value.(string))
		case govaluate.STRING, govaluate.LOGICALOP, govaluate.COMPARATOR:
			queryParts = append(queryParts, strings.ToUpper(token.Value.(string)))
		default:
			queryParts = append(queryParts, fmt.Sprintf("%v", token.Value))
		}
	}
	return strings.Join(queryParts, " "), nil
}
