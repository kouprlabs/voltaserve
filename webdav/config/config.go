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
	Security    SecurityConfig
	Environment config.EnvironmentConfig
}

type SecurityConfig struct {
	APIKey string
}

func GetConfig() *Config {
	cfg := &Config{}
	readPort(cfg)
	readURLs(cfg)
	readS3(cfg)
	readSecurity(cfg)
	readEnvironment(cfg)
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

func readS3(config *Config) {
	config.S3.URL = os.Getenv("S3_URL")
	config.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	config.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	config.S3.Region = os.Getenv("S3_REGION")
	if len(os.Getenv("S3_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("S3_SECURE"))
		if err != nil {
			panic(err)
		}
		config.S3.Secure = v
	}
}

func readSecurity(config *Config) {
	config.Security.APIKey = os.Getenv("SECURITY_API_KEY")
}

func readEnvironment(config *Config) {
	if os.Getenv("TEST") == "true" {
		config.Environment.IsTest = true
	}
}
