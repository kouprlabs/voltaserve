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

	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type TaskCache struct {
	redis     *infra.RedisManager
	taskRepo  *repo.TaskRepo
	keyPrefix string
}

func NewTaskCache() *TaskCache {
	return &TaskCache{
		taskRepo:  repo.NewTaskRepo(),
		redis:     infra.NewRedisManager(),
		keyPrefix: "task:",
	}
}

func (c *TaskCache) Set(file model.Task) error {
	b, err := json.Marshal(file)
	if err != nil {
		return err
	}
	err = c.redis.Set(c.keyPrefix+file.GetID(), string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c *TaskCache) Get(id string) (model.Task, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	task := repo.NewTaskModel()
	if err = json.Unmarshal([]byte(value), &task); err != nil {
		return nil, err
	}
	return task, nil
}

func (c *TaskCache) GetOrNil(id string) model.Task {
	res, err := c.Get(id)
	if err != nil {
		return nil
	}
	return res
}

func (c *TaskCache) Refresh(id string) (model.Task, error) {
	res, err := c.taskRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *TaskCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return nil
	}
	return nil
}
