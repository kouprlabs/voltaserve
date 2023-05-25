package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type OrganizationCache struct {
	redis     *infra.RedisManager
	orgRepo   repo.CoreOrganizationRepo
	keyPrefix string
}

func NewOrganizationCache() *OrganizationCache {
	return &OrganizationCache{
		redis:     infra.NewRedisManager(),
		orgRepo:   repo.NewOrganizationRepo(),
		keyPrefix: "organization:",
	}
}

func (c *OrganizationCache) Set(organization model.CoreOrganization) error {
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

func (c *OrganizationCache) Get(id string) (model.CoreOrganization, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	var org = repo.OrganizationEntity{}
	if err = json.Unmarshal([]byte(value), &org); err != nil {
		return nil, err
	}
	return &org, nil
}

func (c *OrganizationCache) Refresh(id string) (model.CoreOrganization, error) {
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
