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
	blevemapping "github.com/blevesearch/bleve/v2/mapping"
	blevequery "github.com/blevesearch/bleve/v2/search/query"
	bleveindex "github.com/blevesearch/bleve_index_api"
)

var indices map[string]bleve.Index

type bleveSearchManager struct{}

func newBleveSearchManager() SearchManager {
	mgr := &bleveSearchManager{}
	if indices == nil {
		indices = make(map[string]bleve.Index)
		if err := mgr.createFileIndex(); err != nil {
			panic(err)
		}
		if err := mgr.createGroupIndex(); err != nil {
			panic(err)
		}
		if err := mgr.createWorkspaceIndex(); err != nil {
			panic(err)
		}
		if err := mgr.createOrganizationIndex(); err != nil {
			panic(err)
		}
		if err := mgr.createUserIndex(); err != nil {
			panic(err)
		}
		if err := mgr.createTaskIndex(); err != nil {
			panic(err)
		}
	}
	return mgr
}

func (mgr *bleveSearchManager) Query(indexName string, query string, opts SearchQueryOptions) ([]interface{}, error) {
	index, ok := indices[indexName]
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
		res[i] = mgr.documentToMap(doc)
	}
	return res, nil
}

func (mgr *bleveSearchManager) Index(indexName string, models []SearchModel) error {
	index, ok := indices[indexName]
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
	index, ok := indices[indexName]
	if !ok {
		return errors.New("index not found")
	}
	batch := index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	return index.Batch(batch)
}

func (mgr *bleveSearchManager) createFileIndex() error {
	mapping := bleve.NewIndexMapping()
	mgr.appendCommonFields(mapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}
	indices[FileSearchIndex] = index
	return nil
}

func (mgr *bleveSearchManager) createGroupIndex() error {
	mapping := bleve.NewIndexMapping()
	mgr.appendCommonFields(mapping)
	mgr.disableField("members", mapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}
	indices[GroupSearchIndex] = index
	return nil
}

func (mgr *bleveSearchManager) createWorkspaceIndex() error {
	mapping := bleve.NewIndexMapping()
	mgr.appendCommonFields(mapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}
	indices[WorkspaceSearchIndex] = index
	return nil
}

func (mgr *bleveSearchManager) createOrganizationIndex() error {
	mapping := bleve.NewIndexMapping()
	mgr.appendCommonFields(mapping)
	mgr.disableField("members", mapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}
	indices[OrganizationSearchIndex] = index
	return nil
}

func (mgr *bleveSearchManager) createUserIndex() error {
	mapping := bleve.NewIndexMapping()
	mgr.appendCommonFields(mapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}
	indices[UserSearchIndex] = index
	return nil
}

func (mgr *bleveSearchManager) createTaskIndex() error {
	mapping := bleve.NewIndexMapping()
	mgr.appendCommonFields(mapping)
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}
	indices[TaskSearchIndex] = index
	return nil
}

func (mgr *bleveSearchManager) appendCommonFields(mapping *blevemapping.IndexMappingImpl) {
	mapping.DefaultMapping.AddFieldMappingsAt("createTime", bleve.NewTextFieldMapping())
	mapping.DefaultMapping.AddFieldMappingsAt("updateTime", bleve.NewTextFieldMapping())
}

func (mgr *bleveSearchManager) disableField(name string, mapping *blevemapping.IndexMappingImpl) {
	fieldMapping := bleve.NewTextFieldMapping()
	fieldMapping.Index = false
	fieldMapping.Store = false
	mapping.DefaultMapping.AddFieldMappingsAt(name, fieldMapping)
}

func (mgr *bleveSearchManager) buildFilter(filter interface{}) []*blevequery.MatchQuery {
	res := make([]*blevequery.MatchQuery, 0)
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

func (mgr *bleveSearchManager) documentToMap(doc bleveindex.Document) map[string]interface{} {
	res := make(map[string]interface{})
	doc.VisitFields(func(field bleveindex.Field) {
		fieldName := field.Name()
		fieldValue := field.Value()
		switch field.(type) {
		case bleveindex.TextField:
			res[fieldName] = string(fieldValue)
		case bleveindex.NumericField:
			parsedValue, err := strconv.ParseFloat(string(fieldValue), 64)
			if err == nil {
				res[fieldName] = parsedValue
			}
		case bleveindex.BooleanField:
			parsedValue, err := strconv.ParseBool(string(fieldValue))
			if err == nil {
				res[fieldName] = parsedValue
			}
		}
	})
	return res
}
