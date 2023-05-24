package cache

import (
	"encoding/json"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type WorkspaceCache struct {
	redis         *infra.RedisManager
	workspaceRepo repo.CoreWorkspaceRepo
	keyPrefix     string
}

func NewWorkspaceCache() *WorkspaceCache {
	return &WorkspaceCache{
		redis:         infra.NewRedisManager(),
		workspaceRepo: repo.NewPostgresWorkspaceRepo(),
		keyPrefix:     "workspace:",
	}
}

func (c *WorkspaceCache) Set(workspace model.WorkspaceModel) error {
	b, err := json.Marshal(workspace)
	if err != nil {
		return err
	}
	err = c.redis.Set(c.keyPrefix+workspace.GetId(), string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c *WorkspaceCache) Get(id string) (model.WorkspaceModel, error) {
	value, err := c.redis.Get(c.keyPrefix + id)
	if err != nil {
		return c.Refresh(id)
	}
	var workspace = repo.PostgresWorkspace{}
	if err = json.Unmarshal([]byte(value), &workspace); err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (c *WorkspaceCache) Refresh(id string) (model.WorkspaceModel, error) {
	res, err := c.workspaceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err = c.Set(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *WorkspaceCache) Delete(id string) error {
	if err := c.redis.Delete(c.keyPrefix + id); err != nil {
		return err
	}
	return nil
}
