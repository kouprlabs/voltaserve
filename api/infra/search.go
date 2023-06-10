package infra

import (
	"voltaserve/config"

	"github.com/meilisearch/meilisearch-go"
)

var searchClient *meilisearch.Client

const (
	FileSearchIndex         = "file"
	GroupSearchIndex        = "group"
	WorkspaceSearchIndex    = "workspace"
	OrganizationSearchIndex = "organization"
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
	}
	return &SearchManager{
		config: config.GetConfig().Search,
	}
}

func (mgr *SearchManager) Query(index string, query string) ([]interface{}, error) {
	res, err := searchClient.Index(index).Search(query, &meilisearch.SearchRequest{})
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
