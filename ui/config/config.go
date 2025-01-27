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
	config.ConsoleURL = os.Getenv("CONSOLE_URL")
}
