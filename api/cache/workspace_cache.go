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
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type WorkspaceCache struct {
	redis         *infra.RedisManager
	workspaceRepo repo.WorkspaceRepo
	keyPrefix     string
}

func NewWorkspaceCache() *WorkspaceCache {
	return &WorkspaceCache{
		redis:         infra.NewRedisManager(),
		workspaceRepo: repo.NewWorkspaceRepo(),
		keyPrefix:     "workspace:",
	}
}

func (c *WorkspaceCache) Set(workspace model.Workspace) error {
	b, err := json.Marshal(workspace)
	if err != nil {
		return err
	}
	err = c.redis.Set(c.keyPrefix+workspace.GetID(), string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c *WorkspaceCache) Get(id string) (model.Workspace, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	res := repo.NewWorkspace()
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *WorkspaceCache) Refresh(id string) (model.Workspace, error) {
	res, err := c.workspaceRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *WorkspaceCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return err
	}
	return nil
}
