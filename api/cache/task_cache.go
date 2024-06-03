package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type TaskCache struct {
	redis       *infra.RedisManager
	processRepo repo.SnapshotRepo
	keyPrefix   string
}

func NewTaskCacheCache() *TaskCache {
	return &TaskCache{
		processRepo: repo.NewSnapshotRepo(),
		redis:       infra.NewRedisManager(),
		keyPrefix:   "task:",
	}
}

func (c *TaskCache) Set(file model.Snapshot) error {
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

func (c *TaskCache) Get(id string) (model.Snapshot, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	snapshot := repo.NewSnapshot()
	if err = json.Unmarshal([]byte(value), &snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (c *TaskCache) Refresh(id string) (model.Snapshot, error) {
	res, err := c.processRepo.Find(id)
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
