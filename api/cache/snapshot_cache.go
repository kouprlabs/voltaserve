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

type SnapshotCache interface {
	Set(file model.Snapshot) error
	Get(id string) (model.Snapshot, error)
	Refresh(id string) (model.Snapshot, error)
	Delete(id string) error
}

func NewSnapshotCache() SnapshotCache {
	return newSnapshotCache()
}

type snapshotCache struct {
	redis        *infra.RedisManager
	snapshotRepo repo.SnapshotRepo
	keyPrefix    string
}

func newSnapshotCache() *snapshotCache {
	return &snapshotCache{
		snapshotRepo: repo.NewSnapshotRepo(),
		redis:        infra.NewRedisManager(),
		keyPrefix:    "snapshot:",
	}
}

func (c *snapshotCache) Set(file model.Snapshot) error {
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

func (c *snapshotCache) Get(id string) (model.Snapshot, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	res := repo.NewSnapshot()
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *snapshotCache) Refresh(id string) (model.Snapshot, error) {
	res, err := c.snapshotRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *snapshotCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return nil
	}
	return nil
}
