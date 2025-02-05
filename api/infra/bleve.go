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
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	bleve_index "github.com/blevesearch/bleve_index_api"
)

type bleveSearchManager struct {
	indexes              map[string]bleve.Index
	searchableAttributes map[string][]string
	filterableAttributes map[string][]string
}

func newBleveSearchManager() SearchManager {
	manager := &bleveSearchManager{
		indexes:              make(map[string]bleve.Index),
		searchableAttributes: make(map[string][]string),
		filterableAttributes: make(map[string][]string),
	}
	manager.createIndex(
		FileSearchIndex,
		[]string{"name", "text"},
		[]string{"id", "workspaceId", "type", "parentId", "snapshotId", "createTime", "updateTime"},
	)
	manager.createIndex(
		GroupSearchIndex,
		[]string{"name"},
		[]string{"id", "organizationId", "members", "createTime", "updateTime"},
	)
	manager.createIndex(
		WorkspaceSearchIndex,
		[]string{"name"},
		[]string{"id", "storageCapacity", "rootId", "organizationId", "bucket", "createTime", "updateTime"},
	)
	manager.createIndex(
		OrganizationSearchIndex,
		[]string{"name"},
		[]string{"id", "members", "createTime", "updateTime"},
	)
	manager.createIndex(
		UserSearchIndex,
		[]string{"fullName", "username", "email"},
		[]string{"id", "isEmailConfirmed", "createTime", "updateTime"},
	)
	manager.createIndex(
		TaskSearchIndex,
		[]string{"name"},
		[]string{"id", "error", "percentage", "isIndeterminate", "userId", "status", "createTime", "updateTime"},
	)
	return manager
}

func (mgr *bleveSearchManager) Query(indexName string, query string, opts QueryOptions) ([]interface{}, error) {
	index, ok := mgr.indexes[indexName]
	if !ok {
		return nil, errors.New("index not found")
	}
	var searchRequest *bleve.SearchRequest
	var err error
	if query != "" && opts.Filter != nil {
		filterQuery, err := mgr.buildFilter(opts.Filter, mgr.filterableAttributes[indexName])
		if err != nil {
			return nil, err
		}
		searchRequest = bleve.NewSearchRequestOptions(
			bleve.NewConjunctionQuery(
				bleve.NewQueryStringQuery(mgr.buildQuery(query, mgr.searchableAttributes[indexName])),
				bleve.NewQueryStringQuery(filterQuery),
			),
			int(opts.Limit), 0, false,
		)
	} else if query == "" && opts.Filter != nil {
		filterQuery, err := mgr.buildFilter(opts.Filter, mgr.filterableAttributes[indexName])
		if err != nil {
			return nil, err
		}
		searchRequest = bleve.NewSearchRequestOptions(
			bleve.NewQueryStringQuery(filterQuery),
			int(opts.Limit), 0, false,
		)
	} else {
		searchRequest = bleve.NewSearchRequestOptions(
			bleve.NewQueryStringQuery(mgr.buildQuery(query, mgr.searchableAttributes[indexName])),
			int(opts.Limit), 0, false,
		)
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

func (mgr *bleveSearchManager) createIndex(indexName string, searchableAttributes []string, filterableAttributes []string) {
	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		panic(err)
	}
	mgr.indexes[indexName] = index
	mgr.searchableAttributes[indexName] = searchableAttributes
	mgr.filterableAttributes[indexName] = filterableAttributes
}

func (mgr *bleveSearchManager) buildQuery(query string, searchableAttributes []string) string {
	var conditions []string
	for _, attr := range searchableAttributes {
		conditions = append(conditions, fmt.Sprintf(`%s:"%s"`, attr, query))
	}
	return strings.Join(conditions, " OR ")
}

func (mgr *bleveSearchManager) buildFilter(filter interface{}, filterableAttributes []string) (string, error) {
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
		if !slices.Contains(filterableAttributes, field) {
			return "", errors.New(field + " is a non filterable attribute")
		}
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
