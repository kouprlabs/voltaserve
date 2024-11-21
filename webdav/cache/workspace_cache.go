// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

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
	ID              string  `gorm:"column:id;size:36"              json:"id"`
	Name            string  `gorm:"column:name;size:255"           json:"name"`
	StorageCapacity int64   `gorm:"column:storage_capacity"        json:"storageCapacity"`
	RootID          string  `gorm:"column:root_id;size:36"         json:"rootId"`
	OrganizationID  string  `gorm:"column:organization_id;size:36" json:"organizationId"`
	Bucket          string  `gorm:"column:bucket;size:255"         json:"bucket"`
	CreateTime      string  `gorm:"column:create_time"             json:"createTime"`
	UpdateTime      *string `gorm:"column:update_time"             json:"updateTime,omitempty"`
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
