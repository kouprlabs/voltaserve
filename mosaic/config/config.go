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

	"github.com/kouprlabs/voltaserve/shared/config"
)

type Config struct {
	Port        int
	Limits      LimitsConfig
	S3          config.S3Config
	Environment config.EnvironmentConfig
}

type LimitsConfig struct {
	MultipartBodyLengthLimitMB int64
}

func GetConfig() *Config {
	cfg := &Config{}
	readPort(cfg)
	readLimits(cfg)
	config.ReadS3(&cfg.S3)
	config.ReadEnvironment(&cfg.Environment)
	return cfg
}

func readPort(config *Config) {
	if len(os.Getenv("PORT")) > 0 {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			config.Port = port
		}
	}
}

func readLimits(config *Config) {
	if len(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB"), 10, 64)
		if err != nil {
			panic(err)
		}
		config.Limits.MultipartBodyLengthLimitMB = v
	}
}
