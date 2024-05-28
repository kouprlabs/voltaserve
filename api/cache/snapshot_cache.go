package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type SnapshotCache struct {
	redis        *infra.RedisManager
	snapshotRepo repo.SnapshotRepo
	keyPrefix    string
}

func NewSnapshotCache() *SnapshotCache {
	return &SnapshotCache{
		snapshotRepo: repo.NewSnapshotRepo(),
		redis:        infra.NewRedisManager(),
		keyPrefix:    "snapshot:",
	}
}

func (c *SnapshotCache) Set(file model.Snapshot) error {
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

func (c *SnapshotCache) Get(id string) (model.Snapshot, error) {
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
