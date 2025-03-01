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
	Port          int
	PublicUIURL   string
	ConversionURL string
	LanguageURL   string
	MosaicURL     string
	Postgres      config.PostgresConfig
	Search        config.SearchConfig
	Redis         config.RedisConfig
	S3            config.S3Config
	Limits        LimitsConfig
	Security      config.SecurityConfig
	SMTP          SMTPConfig
	Defaults      DefaultsConfig
	Webhook       WebhookConfig
	Environment   config.EnvironmentConfig
}

type LimitsConfig struct {
	FileUploadMB     int
	FileProcessingMB map[string]int
}

type DefaultsConfig struct {
	WorkspaceStorageCapacityMB int
}

type TokenConfig struct {
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	TokenAudience        string
	TokenIssuer          string
}

type SMTPConfig struct {
	Host          string
	Port          int
	Secure        bool
	Username      string
	Password      string
	SenderAddress string
	SenderName    string
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
	readSecurity(cfg)
	readPostgres(cfg)
	readS3(cfg)
	readSearch(cfg)
	readRedis(cfg)
	readSMTP(cfg)
	readLimits(cfg)
	readDefaults(cfg)
	readWebhook(cfg)
	readEnvironment(cfg)
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

func readSecurity(config *Config) {
	config.Security.JWTSigningKey = os.Getenv("SECURITY_JWT_SIGNING_KEY")
	config.Security.CORSOrigins = strings.Split(os.Getenv("SECURITY_CORS_ORIGINS"), ",")
	config.Security.APIKey = os.Getenv("SECURITY_API_KEY")
}

func readPostgres(config *Config) {
	config.Postgres.URL = os.Getenv("POSTGRES_URL")
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

func readSearch(config *Config) {
	config.Search.URL = os.Getenv("SEARCH_URL")
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

func readSMTP(config *Config) {
	config.SMTP.Host = os.Getenv("SMTP_HOST")
	if len(os.Getenv("SMTP_PORT")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.SMTP.Port = int(v)
	}
	if len(os.Getenv("SMTP_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("SMTP_SECURE"))
		if err != nil {
			panic(err)
		}
		config.SMTP.Secure = v
	}
	config.SMTP.Username = os.Getenv("SMTP_USERNAME")
	config.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	config.SMTP.SenderAddress = os.Getenv("SMTP_SENDER_ADDRESS")
	config.SMTP.SenderName = os.Getenv("SMTP_SENDER_NAME")
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

func readWebhook(config *Config) {
	config.Webhook.Snapshot = os.Getenv("WEBHOOK_SNAPSHOT")
}

func readEnvironment(config *Config) {
	if os.Getenv("TEST") == "true" {
		config.Environment.IsTest = true
	}
}
