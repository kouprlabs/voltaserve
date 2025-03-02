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
	Port            int
	APIURL          string
	LanguageURL     string
	MosaicURL       string
	EnableInstaller bool
	Security        config.SecurityConfig
	S3              config.S3Config
	Environment     config.EnvironmentConfig
	Limits          LimitsConfig
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds int
	ImagePreviewMaxWidth          int
	ImagePreviewMaxHeight         int
	MultipartBodyLengthLimitMB    int
}

func GetConfig() *Config {
	cfg := &Config{}
	readPort(cfg)
	readEnableInstaller(cfg)
	readURLs(cfg)
	readLimits(cfg)
	config.ReadSecurity(&cfg.Security)
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

func readEnableInstaller(config *Config) {
	if len(os.Getenv("ENABLE_INSTALLER")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("ENABLE_INSTALLER"))
		if err != nil {
			panic(err)
		}
		config.EnableInstaller = v
	}
}

func readURLs(config *Config) {
	config.APIURL = os.Getenv("API_URL")
	config.LanguageURL = os.Getenv("LANGUAGE_URL")
	config.MosaicURL = os.Getenv("MOSAIC_URL")
}

func readLimits(config *Config) {
	if len(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ExternalCommandTimeoutSeconds = int(v)
	}
	if len(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_WIDTH")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_WIDTH"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ImagePreviewMaxWidth = int(v)
	}
	if len(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_HEIGHT")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_IMAGE_PREVIEW_MAX_HEIGHT"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ImagePreviewMaxHeight = int(v)
	}
	if len(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.MultipartBodyLengthLimitMB = int(v)
	}
}
