package config

import (
	"os"
	"strconv"
	"strings"
)

var config *Config

func GetConfig() Config {
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
	return *config
}

func readURLs(config *Config) {
	config.PublicUIURL = os.Getenv("PUBLIC_UI_URL")
	config.ConversionURL = os.Getenv("CONVERSION_URL")
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
	if len(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ExternalCommandTimeoutSeconds = int(v)
	}
	if len(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.MultipartBodyLengthLimitMB = int(v)
	}
}

func readDefaults(config *Config) {
	if len(os.Getenv("DEFAULT_WORKSPACE_STORAGE_CAPACITY_BYTES")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("DEFAULT_WORKSPACE_STORAGE_CAPACITY_BYTES"), 10, 64)
		if err != nil {
			panic(err)
		}
		config.Defaults.WorkspaceStorageCapacityBytes = v
	}
}
