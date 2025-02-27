// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/kouprlabs/voltaserve/shared/config"
)

type RedisManager struct {
	redisConfig   config.RedisConfig
	client        *redis.Client
	clusterClient *redis.ClusterClient
}

func NewRedisManager(redisConfig config.RedisConfig) *RedisManager {
	return &RedisManager{
		redisConfig: redisConfig,
	}
}

func (mgr *RedisManager) Set(key string, value interface{}) error {
	if err := mgr.Connect(); err != nil {
		return err
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
	if err := mgr.Connect(); err != nil {
		return "", err
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
	if err := mgr.Connect(); err != nil {
		return err
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

func (mgr *RedisManager) Connect() error {
	if mgr.client != nil || mgr.clusterClient != nil {
		return nil
	}
	addresses := strings.Split(mgr.redisConfig.Address, ";")
	if len(addresses) > 1 {
		mgr.clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addresses,
			Password: mgr.redisConfig.Password,
		})
		if err := mgr.clusterClient.Ping(context.Background()).Err(); err != nil {
			return err
		}
	} else {
		mgr.client = redis.NewClient(&redis.Options{
			Addr:     mgr.redisConfig.Address,
			Password: mgr.redisConfig.Password,
			DB:       mgr.redisConfig.DB,
		})
		if err := mgr.client.Ping(context.Background()).Err(); err != nil {
			return err
		}
	}
	return nil
}
