package infra

import (
	"context"
	"strings"
	"voltaserve/config"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	config        config.RedisConfig
	client        *redis.Client
	clusterClient *redis.ClusterClient
}

func NewRedisManager() *RedisManager {
	mgr := new(RedisManager)
	mgr.config = config.GetConfig().Redis
	return mgr
}

func (mgr *RedisManager) Set(key string, value interface{}) error {
	if mgr.client == nil && mgr.clusterClient == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	if mgr.clusterClient != nil {
		if _, err := mgr.clusterClient.Set(context.Background(), key, value, 0).Result(); err != nil {
			return err
		}
	} else {
		if _, err := mgr.client.Set(context.Background(), key, value, 0).Result(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *RedisManager) Get(key string) (string, error) {
	if mgr.client == nil && mgr.clusterClient == nil {
		if err := mgr.connect(); err != nil {
			return "", err
		}
	}
	if mgr.clusterClient != nil {
		value, err := mgr.clusterClient.Get(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		return value, nil
	} else {
		value, err := mgr.client.Get(context.Background(), key).Result()
		if err != nil {
			return "", err
		}
		return value, nil
	}
}

func (mgr *RedisManager) Delete(key string) error {
	if mgr.client == nil && mgr.clusterClient == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	if mgr.clusterClient != nil {
		if _, err := mgr.clusterClient.Del(context.Background(), key).Result(); err != nil {
			return err
		}
	} else {
		if _, err := mgr.client.Del(context.Background(), key).Result(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *RedisManager) Close() error {
	if mgr.client != nil {
		if err := mgr.client.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *RedisManager) connect() error {
	if mgr.client != nil || mgr.clusterClient != nil {
		return nil
	}
	addresses := strings.Split(mgr.config.Address, ";")
	if len(addresses) > 1 {
		mgr.clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addresses,
			Password: mgr.config.Password,
		})
		if err := mgr.clusterClient.Ping(context.Background()).Err(); err != nil {
			return err
		}
	} else {
		mgr.client = redis.NewClient(&redis.Options{
			Addr:     mgr.config.Address,
			Password: mgr.config.Password,
			DB:       mgr.config.DB,
		})
		if err := mgr.client.Ping(context.Background()).Err(); err != nil {
			return err
		}
	}
	return nil
}
