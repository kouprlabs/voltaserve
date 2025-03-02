// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package config

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func ReadRedis(config *RedisConfig) {
	config.Address = os.Getenv("REDIS_ADDRESS")
	config.Password = os.Getenv("REDIS_PASSWORD")
	if len(os.Getenv("REDIS_DB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.DB = int(v)
	}
}
