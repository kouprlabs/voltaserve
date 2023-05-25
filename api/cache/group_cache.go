package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type GroupCache struct {
	redis     *infra.RedisManager
	groupRepo repo.CoreGroupRepo
	keyPrefix string
}

func NewGroupCache() *GroupCache {
	return &GroupCache{
		redis:     infra.NewRedisManager(),
		groupRepo: repo.NewPostgresGroupRepo(),
		keyPrefix: "group:",
	}
}

func (c *GroupCache) Set(workspace model.GroupModel) error {
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

func (c *GroupCache) Get(id string) (model.GroupModel, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	var group = repo.PostgresGroup{}
	if err = json.Unmarshal([]byte(value), &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *GroupCache) Refresh(id string) (model.GroupModel, error) {
	res, err := c.groupRepo.Find(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}
