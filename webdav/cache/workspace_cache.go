// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package cache

import (
	"encoding/json"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

type WorkspaceCache struct {
	redis     *infra.RedisManager
	keyPrefix string
}

func NewWorkspaceCache() *WorkspaceCache {
	return &WorkspaceCache{
		redis:     infra.NewRedisManager(),
		keyPrefix: "workspace:",
	}
}

type Workspace struct {
	ID              string  `json:"id," gorm:"column:id;size:36"`
	Name            string  `json:"name" gorm:"column:name;size:255"`
	StorageCapacity int64   `json:"storageCapacity" gorm:"column:storage_capacity"`
	RootID          string  `json:"rootId" gorm:"column:root_id;size:36"`
	OrganizationID  string  `json:"organizationId" gorm:"column:organization_id;size:36"`
	Bucket          string  `json:"bucket" gorm:"column:bucket;size:255"`
	CreateTime      string  `json:"createTime" gorm:"column:create_time"`
	UpdateTime      *string `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (c *WorkspaceCache) Get(id string) (*Workspace, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return nil, err
	}
	var res Workspace
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return &res, nil
}
