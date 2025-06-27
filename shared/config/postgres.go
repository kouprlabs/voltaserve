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

type PostgresConfig struct {
	URL                          string
	MaxIdleConnections           int
	MaxOpenConnections           int
	ConnectionMaxIdleTimeMinutes int
}

func ReadPostgres(config *PostgresConfig) {
	config.URL = os.Getenv("POSTGRES_URL")
	if len(os.Getenv("POSTGRES_MAX_IDLE_CONNECTIONS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("POSTGRES_MAX_IDLE_CONNECTIONS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.MaxIdleConnections = int(v)
	}
	if len(os.Getenv("POSTGRES_MAX_OPEN_CONNECTIONS")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("POSTGRES_MAX_OPEN_CONNECTIONS"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.MaxOpenConnections = int(v)
	}
	if len(os.Getenv("POSTGRES_CONNECTION_MAX_IDLE_TIME_MINUTES")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("POSTGRES_CONNECTION_MAX_IDLE_TIME_MINUTES"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.ConnectionMaxIdleTimeMinutes = int(v)
	}
}
