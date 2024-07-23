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

	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type GroupCache struct {
	redis     *infra.RedisManager
	groupRepo repo.GroupRepo
	keyPrefix string
}

func NewGroupCache() *GroupCache {
	return &GroupCache{
		redis:     infra.NewRedisManager(),
		groupRepo: repo.NewGroupRepo(),
		keyPrefix: "group:",
	}
}

func (c *GroupCache) Set(workspace model.Group) error {
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

func (c *GroupCache) Get(id string) (model.Group, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	res := repo.NewGroup()
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *GroupCache) Refresh(id string) (model.Group, error) {
	res, err := c.groupRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *GroupCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return err
	}
	return nil
}
