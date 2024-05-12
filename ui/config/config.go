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
		readURLs(config)
	}
	return *config
}

func readURLs(config *Config) {
	config.APIURL = os.Getenv("API_URL")
	config.IDPURL = os.Getenv("IDP_URL")
}
