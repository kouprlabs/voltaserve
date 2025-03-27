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
	"strings"

	"github.com/kouprlabs/voltaserve/shared/config"
)

type Config struct {
	Port             int
	PublicUIURL      string
	ConversionURL    string
	LanguageURL      string
	MosaicURL        string
	Postgres         config.PostgresConfig
	Search           config.SearchConfig
	Redis            config.RedisConfig
	S3               config.S3Config
	Security         config.SecurityConfig
	Environment      config.EnvironmentConfig
	SMTP             config.SMTPConfig
	Limits           LimitsConfig
	Defaults         DefaultsConfig
	SnapshotWebhook  string
	WorkspaceWebhook string
}

type LimitsConfig struct {
	FileUploadMB     int
	FileProcessingMB map[string]int
}

type DefaultsConfig struct {
	WorkspaceStorageCapacityMB int
}

type WebhookConfig struct {
	Snapshot string
}

const (
	FileTypePDF            = "pdf"
	FileTypeOffice         = "office"
	FileTypePlainText      = "plain_text"
	FileTypeImage          = "image"
	FileTypeVideo          = "video"
	FileTypeAudio          = "audio"
	FileTypeGLB            = "glb"
	FileTypeZIP            = "zip"
	FileTypeGLTF           = "gltf"
	FileTypeEverythingElse = "*"
)

func GetConfig() *Config {
	cfg := &Config{}
	readPort(cfg)
	readURLs(cfg)
	readLimits(cfg)
	readDefaults(cfg)
	readWebhooks(cfg)
	config.ReadSecurity(&cfg.Security)
	config.ReadPostgres(&cfg.Postgres)
	config.ReadS3(&cfg.S3)
	config.ReadSearch(&cfg.Search)
	config.ReadRedis(&cfg.Redis)
	config.ReadSMTP(&cfg.SMTP)
	config.ReadEnvironment(&cfg.Environment)
	return cfg
}

func (l *LimitsConfig) GetFileProcessingMB(fileType string) int {
	v, ok := l.FileProcessingMB[fileType]
	if !ok {
		return l.FileProcessingMB[FileTypeEverythingElse]
	}
	return v
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
	config.PublicUIURL = os.Getenv("PUBLIC_UI_URL")
	config.ConversionURL = os.Getenv("CONVERSION_URL")
	config.LanguageURL = os.Getenv("LANGUAGE_URL")
	config.MosaicURL = os.Getenv("MOSAIC_URL")
}

func readLimits(config *Config) {
	if len(os.Getenv("LIMITS_FILE_UPLOAD_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_FILE_UPLOAD_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.FileUploadMB = int(v)
	}
	if len(os.Getenv("LIMITS_FILE_PROCESSING_MB")) > 0 {
		raw := os.Getenv("LIMITS_FILE_PROCESSING_MB")
		parts := strings.Split(raw, ",")
		config.Limits.FileProcessingMB = make(map[string]int)
		for _, part := range parts {
			limit := strings.Split(part, ":")
			if len(limit) != 2 {
				panic("invalid LIMITS_FILE_PROCESSING_MB format")
			}
			v, err := strconv.ParseInt(limit[1], 10, 32)
			if err != nil {
				panic(err)
			}
			config.Limits.FileProcessingMB[limit[0]] = int(v)
		}
	}
}

func readDefaults(config *Config) {
	if len(os.Getenv("DEFAULTS_WORKSPACE_STORAGE_CAPACITY_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("DEFAULTS_WORKSPACE_STORAGE_CAPACITY_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Defaults.WorkspaceStorageCapacityMB = int(v)
	}
}

func readWebhooks(config *Config) {
	config.SnapshotWebhook = os.Getenv("SNAPSHOT_WEBHOOK")
	config.WorkspaceWebhook = os.Getenv("WORKSPACE_WEBHOOK")
}
