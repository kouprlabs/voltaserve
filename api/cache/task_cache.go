package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type TaskCache struct {
	redis     *infra.RedisManager
	taskRepo  repo.TaskRepo
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
	task := repo.NewTask()
	if err = json.Unmarshal([]byte(value), &task); err != nil {
		return nil, err
	}
	return task, nil
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
