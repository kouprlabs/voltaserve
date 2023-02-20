package config

import (
	"io"
	"os"
	"strconv"
	"strings"

	"sigs.k8s.io/yaml"
)

type LimitsConfig struct {
	ExternalCommandTimeoutSec  int `json:"externalCommandTimeoutSec"`
	FileProcessingMaxSizeMb    int `json:"fileProcessingMaxSizeMb"`
	ImagePreviewMaxWidth       int `json:"imagePreviewMaxWidth"`
	ImagePreviewMaxHeight      int `json:"imagePreviewMaxHeight"`
	MultipartBodyLengthLimitMb int `json:"multipartBodyLengthLimitMb"`
}

type TokenConfig struct {
	AccessTokenLifetime  int    `json:"accessTokenLifetime"`
	RefreshTokenLifetime int    `json:"refreshTokenLifetime"`
	TokenAudience        string `json:"tokenAudience"`
	TokenIssuer          string `json:"tokenIssuer"`
}

type SearchConfig struct {
	Url string `json:"url"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

type S3Config struct {
	Url       string `json:"url"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region"`
	Secure    bool   `json:"secure"`
}

type SecurityConfig struct {
	JwtSigningKey string   `json:"jwtSigningKey"`
	CorsOrigins   []string `json:"corsOrigins"`
}

type SmtpConfig struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Secure        bool   `json:"secure"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	SenderAddress string `json:"senderAddress"`
	SenderName    string `json:"senderName"`
}

type Config struct {
	Url         string         `json:"url"`
	WebUrl      string         `json:"webUrl"`
	DatabaseUrl string         `json:"databaseUrl"`
	Search      SearchConfig   `json:"search"`
	Redis       RedisConfig    `json:"redis"`
	S3          S3Config       `json:"s3"`
	Limits      LimitsConfig   `json:"limits"`
	Security    SecurityConfig `json:"security"`
	Smtp        SmtpConfig     `json:"smtp"`
}

var config *Config

func GetConfig() Config {
	if config != nil {
		return *config
	}
	var filename string
	if _, err := os.Stat("./config.local.yml"); err == nil {
		filename = "./config.local.yml"
	} else {
		filename = "./config.yml"
	}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	config = &Config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		panic(err)
	}
	if len(os.Getenv("URL")) > 0 {
		config.Url = os.Getenv("URL")
	}
	if len(os.Getenv("WEB_URL")) > 0 {
		config.WebUrl = os.Getenv("WEB_URL")
	}
	if len(os.Getenv("DATABASE_URL")) > 0 {
		config.DatabaseUrl = os.Getenv("DATABASE_URL")
	}
	overrideSecurityConfig(config)
	overrideS3Config(config)
	overrideSearchConfig(config)
	overrideRedisConfig(config)
	overrideSmtpConfig(config)
	overrideLimits(config)
	return *config
}

func overrideSecurityConfig(config *Config) {
	if len(os.Getenv("SECURITY_JWT_SIGNING_KEY")) > 0 {
		config.Security.JwtSigningKey = os.Getenv("SECURITY_JWT_SIGNING_KEY")
	}
	if len(os.Getenv("SECURITY_CORS_ORIGINS")) > 0 {
		config.Security.CorsOrigins = strings.Split(os.Getenv("SECURITY_CORS_ORIGINS"), ",")
	}
}

func overrideS3Config(config *Config) {
	if len(os.Getenv("S3_URL")) > 0 {
		config.S3.Url = os.Getenv("S3_URL")
	}
	if len(os.Getenv("S3_ACCESS_KEY")) > 0 {
		config.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	}
	if len(os.Getenv("S3_SECRET_KEY")) > 0 {
		config.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	}
	if len(os.Getenv("S3_REGION")) > 0 {
		config.S3.Region = os.Getenv("S3_REGION")
	}
	if len(os.Getenv("S3_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("S3_SECURE"))
		if err != nil {
			panic(err)
		}
		config.S3.Secure = v
	}
}

func overrideSearchConfig(config *Config) {
	if len(os.Getenv("SEARCH_URL")) > 0 {
		config.Search.Url = os.Getenv("SEARCH_URL")
	}
}

func overrideRedisConfig(config *Config) {
	if len(os.Getenv("REDIS_ADDR")) > 0 {
		config.Redis.Addr = os.Getenv("REDIS_ADDR")
	}
	if len(os.Getenv("REDIS_PASSWORD")) > 0 {
		config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	}
	if len(os.Getenv("REDIS_DB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Redis.Db = int(v)
	}
}

func overrideSmtpConfig(config *Config) {
	if len(os.Getenv("SMTP_HOST")) > 0 {
		config.Smtp.Host = os.Getenv("SMTP_HOST")
	}
	if len(os.Getenv("SMTP_PORT")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Smtp.Port = int(v)
	}
	if len(os.Getenv("SMTP_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("SMTP_SECURE"))
		if err != nil {
			panic(err)
		}
		config.Smtp.Secure = v
	}
	if len(os.Getenv("SMTP_USERNAME")) > 0 {
		config.Smtp.Username = os.Getenv("SMTP_USERNAME")
	}
	if len(os.Getenv("SMTP_PASSWORD")) > 0 {
		config.Smtp.Password = os.Getenv("SMTP_PASSWORD")
	}
	if len(os.Getenv("SMTP_SENDER_ADDRESS")) > 0 {
		config.Smtp.SenderAddress = os.Getenv("SMTP_SENDER_ADDRESS")
	}
	if len(os.Getenv("SMTP_SENDER_NAME")) > 0 {
		config.Smtp.SenderName = os.Getenv("SMTP_SENDER_NAME")
	}
}

func overrideLimits(config *Config) {
	if len(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SEC")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SEC"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ExternalCommandTimeoutSec = int(v)
	}
	if len(os.Getenv("LIMITS_FILE_PROCESSING_MAX_SIZE_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_FILE_PROCESSING_MAX_SIZE_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.FileProcessingMaxSizeMb = int(v)
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
		config.Limits.MultipartBodyLengthLimitMb = int(v)
	}
}
