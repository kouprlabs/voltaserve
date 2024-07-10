// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port          int
	PublicUIURL   string
	ConversionURL string
	LanguageURL   string
	MosaicURL     string
	DatabaseURL   string
	Search        SearchConfig
	Redis         RedisConfig
	S3            S3Config
	Limits        LimitsConfig
	Security      SecurityConfig
	SMTP          SMTPConfig
	Defaults      DefaultsConfig
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

type SearchConfig struct {
	URL string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}

type SecurityConfig struct {
	JWTSigningKey string
	CORSOrigins   []string
	APIKey        string
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

var config *Config

func GetConfig() *Config {
	if config == nil {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		config = &Config{
			Port: port,
		}
		readURLs(config)
		readSecurity(config)
		readS3(config)
		readSearch(config)
		readRedis(config)
		readSMTP(config)
		readLimits(config)
		readDefaults(config)
	}
	return config
}

func (l *LimitsConfig) GetFileProcessingMB(fileType string) int {
	v, ok := l.FileProcessingMB[fileType]
	if !ok {
		return l.FileProcessingMB[FileTypeEverythingElse]
	}
	return v
}

func readURLs(config *Config) {
	config.PublicUIURL = os.Getenv("PUBLIC_UI_URL")
	config.ConversionURL = os.Getenv("CONVERSION_URL")
	config.LanguageURL = os.Getenv("LANGUAGE_URL")
	config.MosaicURL = os.Getenv("MOSAIC_URL")
	config.DatabaseURL = os.Getenv("POSTGRES_URL")
}

func readSecurity(config *Config) {
	config.Security.JWTSigningKey = os.Getenv("SECURITY_JWT_SIGNING_KEY")
	config.Security.CORSOrigins = strings.Split(os.Getenv("SECURITY_CORS_ORIGINS"), ",")
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
