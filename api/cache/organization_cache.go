package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type OrganizationCache struct {
	redis     *infra.RedisManager
	orgRepo   repo.OrganizationRepo
	keyPrefix string
}

func NewOrganizationCache() *OrganizationCache {
	return &OrganizationCache{
		redis:     infra.NewRedisManager(),
		orgRepo:   repo.NewOrganizationRepo(),
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
	var org = repo.NewOrganization()
	if err = json.Unmarshal([]byte(value), &org); err != nil {
		return nil, err
	}
	return org, nil
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
