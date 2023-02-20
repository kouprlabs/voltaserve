package infra

import (
	"context"
	"voltaserve/config"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	config config.RedisConfig
	client *redis.Client
}

func NewRedisManager() *RedisManager {
	mgr := new(RedisManager)
	mgr.config = config.GetConfig().Redis
	return mgr
}

func (mgr *RedisManager) Set(key string, value interface{}) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	if err := mgr.client.Set(context.Background(), key, value, 0); err != nil {
		return nil
	}
	return nil
}

func (mgr *RedisManager) Get(key string) (string, error) {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return "", err
		}
	}
	value, err := mgr.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (mgr *RedisManager) Delete(key string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	mgr.client.Del(context.Background(), key)
	return nil
}

func (mgr *RedisManager) Close() error {
	if mgr.client != nil {
		if err := mgr.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *RedisManager) connect() error {
	if mgr.client != nil {
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     mgr.config.Addr,
		Password: mgr.config.Password,
		DB:       mgr.config.Db,
	})
	mgr.client = client
	return nil
}
