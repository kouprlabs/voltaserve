package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port   int
	Limits LimitsConfig
	S3     S3Config
}

type LimitsConfig struct {
	MultipartBodyLengthLimitMB int
}

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
}

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
		readS3(config)
		readLimits(config)
	}
	return config
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
	if len(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_MULTIPART_BODY_LENGTH_LIMIT_MB"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.MultipartBodyLengthLimitMB = int(v)
	}
}
