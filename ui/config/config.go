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

func GetConfig() Config {
	config := &Config{}
	readPort(config)
	readURLs(config)
	return *config
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
	config.APIURL = os.Getenv("API_URL")
	config.IDPURL = os.Getenv("IDP_URL")
	config.ConsoleURL = os.Getenv("CONSOLE_URL")
}
