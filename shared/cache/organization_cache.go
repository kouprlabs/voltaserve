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

type OrganizationCache struct {
	redis     *infra.RedisManager
	orgRepo   *repo.OrganizationRepo
	keyPrefix string
}

func NewOrganizationCache(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *OrganizationCache {
	return &OrganizationCache{
		redis:     infra.NewRedisManager(redis),
		orgRepo:   repo.NewOrganizationRepo(postgres, environment),
		keyPrefix: "organization:",
	}
}

func (c *OrganizationCache) Set(organization model.Organization) error {
	b, err := json.Marshal(organization)
	if err != nil {
		return err
	}
	err = c.redis.Set(c.keyPrefix+organization.GetID(), string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c *OrganizationCache) Get(id string) (model.Organization, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	res := repo.NewOrganizationModel()
	if err = json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *OrganizationCache) GetOrNil(id string) model.Organization {
	res, err := c.Get(id)
	if err != nil {
		return nil
	}
	return res
}

func (c *OrganizationCache) Refresh(id string) (model.Organization, error) {
	res, err := c.orgRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *OrganizationCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return err
	}
	return nil
}
