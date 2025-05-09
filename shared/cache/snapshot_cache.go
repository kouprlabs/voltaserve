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

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"
)

type SnapshotCache struct {
	redis        *infra.RedisManager
	snapshotRepo *repo.SnapshotRepo
	keyPrefix    string
}

func NewSnapshotCache(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *SnapshotCache {
	return &SnapshotCache{
		redis:        infra.NewRedisManager(redis),
		snapshotRepo: repo.NewSnapshotRepo(postgres, environment),
		keyPrefix:    "snapshot:",
	}
}

func (c *SnapshotCache) Set(snapshot model.Snapshot) error {
	b, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	err = c.redis.Set(c.keyPrefix+snapshot.GetID(), string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c *SnapshotCache) Get(id string) (model.Snapshot, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	res := repo.NewSnapshotModel()
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *SnapshotCache) GetOrNil(id string) model.Snapshot {
	res, err := c.Get(id)
	if err != nil {
		return nil
	}
	return res
}

func (c *SnapshotCache) Refresh(id string) (model.Snapshot, error) {
	res, err := c.snapshotRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *SnapshotCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return nil
	}
	return nil
}
