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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
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
	var bleveQuery string
	var err error
	if opts.Filter != nil {
		filterQuery, err := mgr.parseFilter(opts.Filter)
		if err != nil {
			return nil, err
		}
		bleveQuery = query + " " + filterQuery
	} else {
		bleveQuery = query
	}
	searchResult, err := index.Search(
		bleve.NewSearchRequestOptions(
			bleve.NewQueryStringQuery(bleveQuery),
			int(opts.Limit), 0, false),
	)
	if err != nil {
		return nil, err
	}
	res := make([]interface{}, len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
		doc, err := index.Document(hit.ID)
		if err != nil {
			return nil, err
		}
		docMap := make(map[string]interface{})
		doc.VisitFields(func(field bleve_index.Field) {
			fieldName := field.Name()
			fieldValue := field.Value()
			switch field.(type) {
			case bleve_index.TextField:
				docMap[fieldName] = string(fieldValue)
			case bleve_index.NumericField:
				num, err := strconv.ParseFloat(string(fieldValue), 64)
				if err == nil {
					docMap[fieldName] = num
				}
			case bleve_index.DateTimeField:
				dateTime, err := time.Parse(time.RFC3339, string(fieldValue))
				if err == nil {
					docMap[fieldName] = dateTime
				}
			case bleve_index.BooleanField:
				boolVal, err := strconv.ParseBool(string(fieldValue))
				if err == nil {
					docMap[fieldName] = boolVal
				}
			default:
				docMap[fieldName] = string(fieldValue)
			}
		})
		res[i] = docMap
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

func (mgr *bleveSearchManager) parseFilter(filter interface{}) (string, error) {
	filterStr, ok := filter.(string)
	if !ok {
		return "", errors.New("filter must be a string")
	}
	re := regexp.MustCompile(`(\w+)\s*=\s*("[^"]*"|\d+)`)
	matches := re.FindAllStringSubmatch(filterStr, -1)
	if len(matches) == 0 {
		return "", errors.New("invalid filter format")
	}
	var conditions []string
	for _, match := range matches {
		field, value := match[1], match[2]
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			conditions = append(conditions, fmt.Sprintf(`%s:%s`, field, value))
		} else {
			conditions = append(conditions, fmt.Sprintf(`%s:%s`, field, value))
		}
	}
	operators := re.ReplaceAllString(filterStr, "")
	operatorParts := strings.Fields(operators)
	finalConditions := make([]string, 0, len(conditions)+len(operatorParts))
	for i, condition := range conditions {
		finalConditions = append(finalConditions, condition)
		if i < len(operatorParts) {
			finalConditions = append(finalConditions, operatorParts[i])
		}
	}
	queryString := strings.Join(finalConditions, " ")
	return queryString, nil
}
