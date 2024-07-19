// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import (
	"github.com/meilisearch/meilisearch-go"

	"github.com/kouprlabs/voltaserve/api/config"
)

var searchClient *meilisearch.Client

const (
	FileSearchIndex         = "file"
	GroupSearchIndex        = "group"
	WorkspaceSearchIndex    = "workspace"
	OrganizationSearchIndex = "organization"
	TaskSearchIndex         = "task"
	UserSearchIndex         = "user"
)

type SearchModel interface {
	GetID() string
}

type SearchManager struct {
	config config.SearchConfig
}

func NewSearchManager() *SearchManager {
	if searchClient == nil {
		searchClient = meilisearch.NewClient(meilisearch.ClientConfig{
			Host: config.GetConfig().Search.URL,
		})
		/* Configure file index */
		if _, err := searchClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        FileSearchIndex,
			PrimaryKey: "id",
		}); err != nil {
			panic(err)
		}
		if _, err := searchClient.Index(FileSearchIndex).UpdateSettings(&meilisearch.Settings{
			SearchableAttributes: []string{"name", "text"},
		}); err != nil {
			panic(err)
		}
		/* Configure group index */
		if _, err := searchClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        GroupSearchIndex,
			PrimaryKey: "id",
		}); err != nil {
			panic(err)
		}
		if _, err := searchClient.Index(GroupSearchIndex).UpdateSettings(&meilisearch.Settings{
			SearchableAttributes: []string{"name"},
		}); err != nil {
			panic(err)
		}
		/* Configure workspace index */
		if _, err := searchClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        WorkspaceSearchIndex,
			PrimaryKey: "id",
		}); err != nil {
			panic(err)
		}
		if _, err := searchClient.Index(WorkspaceSearchIndex).UpdateSettings(&meilisearch.Settings{
			SearchableAttributes: []string{"name"},
		}); err != nil {
			panic(err)
		}
		/* Configure organization index */
		if _, err := searchClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        OrganizationSearchIndex,
			PrimaryKey: "id",
		}); err != nil {
			panic(err)
		}
		if _, err := searchClient.Index(OrganizationSearchIndex).UpdateSettings(&meilisearch.Settings{
			SearchableAttributes: []string{"name"},
		}); err != nil {
			panic(err)
		}
		/* Configure user index */
		if _, err := searchClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        UserSearchIndex,
			PrimaryKey: "id",
		}); err != nil {
			panic(err)
		}
		if _, err := searchClient.Index(UserSearchIndex).UpdateSettings(&meilisearch.Settings{
			SearchableAttributes: []string{"fullName", "email"},
		}); err != nil {
			panic(err)
		}
		/* Configure task index */
		if _, err := searchClient.CreateIndex(&meilisearch.IndexConfig{
			Uid:        TaskSearchIndex,
			PrimaryKey: "id",
		}); err != nil {
			panic(err)
		}
		if _, err := searchClient.Index(TaskSearchIndex).UpdateSettings(&meilisearch.Settings{
			SearchableAttributes: []string{"name"},
		}); err != nil {
			panic(err)
		}
	}
	return &SearchManager{
		config: config.GetConfig().Search,
	}
}

type QueryOptions struct {
	Limit int64
}

func (mgr *SearchManager) Query(index string, query string, opts QueryOptions) ([]interface{}, error) {
	res, err := searchClient.Index(index).Search(query, &meilisearch.SearchRequest{Limit: opts.Limit})
	if err != nil {
		return nil, err
	}
	return res.Hits, nil
}

func (mgr *SearchManager) Index(index string, models []SearchModel) error {
	_, err := searchClient.Index(index).AddDocuments(models)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *SearchManager) Update(index string, m []SearchModel) error {
	_, err := searchClient.Index(index).UpdateDocuments(m)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *SearchManager) Delete(index string, ids []string) error {
	_, err := searchClient.Index(index).DeleteDocuments(ids)
	if err != nil {
		return err
	}
	return nil
}
