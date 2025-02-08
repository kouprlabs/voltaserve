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
	Port            int
	APIURL          string
	LanguageURL     string
	MosaicURL       string
	EnableInstaller bool
	Security        SecurityConfig
	Limits          LimitsConfig
	S3              S3Config
}

type SecurityConfig struct {
	APIKey string
}

type LimitsConfig struct {
	ExternalCommandTimeoutSeconds     int
	ImagePreviewMaxWidth              int
	ImagePreviewMaxHeight             int
	MultipartBodyLengthLimitMB        int
	ImageMosaicTriggerThresholdPixels int
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}

func GetConfig() *Config {
	config := &Config{}
	readPort(config)
	readEnableInstaller(config)
	readURLs(config)
	readSecurity(config)
	readS3(config)
	readLimits(config)
	return config
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

func readSecurity(config *Config) {
	config.Security.APIKey = os.Getenv("SECURITY_API_KEY")
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
	if len(os.Getenv("LIMITS_IMAGE_MOSAIC_TRIGGER_THRESHOLD_PIXELS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_IMAGE_MOSAIC_TRIGGER_THRESHOLD_PIXELS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ImageMosaicTriggerThresholdPixels = int(v)
	}
}
