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

import "github.com/kouprlabs/voltaserve/api/config"

type SearchManager interface {
	Query(index string, query string, opts QueryOptions) ([]interface{}, error)
	Index(index string, models []SearchModel) error
	Update(index string, models []SearchModel) error
	Delete(index string, ids []string) error
}

func NewSearchManager() SearchManager {
	if config.GetConfig().Environment.IsTest {
		return newBleveSearchManager()
	} else {
		return newMeilisearchManager()
	}
}

type SearchModel interface {
	GetID() string
}

type QueryOptions struct {
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
