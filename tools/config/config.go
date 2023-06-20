package config

import (
	"os"
	"strconv"
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
		readSecurity(config)
		readLimits(config)
	}
	return *config
}

func readSecurity(config *Config) {
	config.Security.APIKey = os.Getenv("SECURITY_API_KEY")
}

func readLimits(config *Config) {
	if len(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("LIMITS_EXTERNAL_COMMAND_TIMEOUT_SECONDS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Limits.ExternalCommandTimeoutSeconds = int(v)
	}
}
