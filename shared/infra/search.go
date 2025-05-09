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

import "github.com/kouprlabs/voltaserve/shared/config"

type SearchManager interface {
	Query(index string, query string, opts SearchQueryOptions) ([]interface{}, error)
	Index(index string, models []SearchModel) error
	Update(index string, models []SearchModel) error
	Delete(index string, ids []string) error
}

func NewSearchManager(searchConfig config.SearchConfig, envConfig config.EnvironmentConfig) SearchManager {
	if envConfig.IsTest {
		return newBleveSearchManager()
	} else {
		return newMeilisearchManager(searchConfig)
	}
}

type SearchModel interface {
	GetID() string
}

type SearchQueryOptions struct {
	Limit  int64
	Filter interface{}
}

const (
	FileSearchIndex         = "file"
	GroupSearchIndex        = "group"
	WorkspaceSearchIndex    = "workspace"
	OrganizationSearchIndex = "organization"
	TaskSearchIndex         = "task"
	UserSearchIndex         = "user"
)
