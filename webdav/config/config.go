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

type Config struct {
	Host     string
	Port     string
	APIURL   string
	IdPURL   string
	S3       S3Config
	Redis    RedisConfig
	Security SecurityConfig
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type SecurityConfig struct {
	APIKey string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Port: os.Getenv("PORT"),
		}
		readURLs(config)
		readS3(config)
		readRedis(config)
		readSecurity(config)
	}
	return config
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

func readRedis(config *Config) {
	config.Redis.Address = os.Getenv("REDIS_ADDRESS")
	config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	if len(os.Getenv("REDIS_DB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Redis.DB = int(v)
	}
}

func readSecurity(config *Config) {
	config.Security.APIKey = os.Getenv("SECURITY_API_KEY")
}
