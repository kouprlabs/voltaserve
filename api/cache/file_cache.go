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

type FileCache struct {
	redis     *infra.RedisManager
	fileRepo  *repo.FileRepo
	keyPrefix string
}

func NewFileCache() *FileCache {
	return &FileCache{
		fileRepo:  repo.NewFileRepo(),
		redis:     infra.NewRedisManager(),
		keyPrefix: "file:",
	}
}

func (c *FileCache) Set(file model.File) error {
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

func (c *FileCache) Get(id string) (model.File, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	res := repo.NewFileModel()
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *FileCache) GetOrNil(id string) (model.File, error) {
	res, err := c.Get(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *FileCache) Refresh(id string) (model.File, error) {
	res, err := c.fileRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *FileCache) RefreshWithExisting(file model.File, userID string) (model.File, error) {
	err := c.fileRepo.PopulateModelFieldsForUser([]model.File{file}, userID)
	if err != nil {
		return nil, err
	}
	if err = c.Set(file); err != nil {
		return nil, err
	}
	return file, nil
}

func (c *FileCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return nil
	}
	return nil
}
