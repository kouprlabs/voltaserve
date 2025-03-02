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
	Host        string
	Port        int
	APIURL      string
	IdPURL      string
	S3          config.S3Config
	Security    config.SecurityConfig
	Environment config.EnvironmentConfig
}

func GetConfig() *Config {
	cfg := &Config{}
	readPort(cfg)
	readURLs(cfg)
	config.ReadS3(&cfg.S3)
	config.ReadSecurity(&cfg.Security)
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

func readURLs(config *Config) {
	config.APIURL = os.Getenv("API_URL")
	config.IdPURL = os.Getenv("IDP_URL")
}
