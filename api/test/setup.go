// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package test

import (
	"fmt"
	"os"

	"github.com/alicebob/miniredis/v2"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func SetupPostgres(port uint32) (*embeddedpostgres.EmbeddedPostgres, error) {
	if err := os.Setenv("DEFAULTS_WORKSPACE_STORAGE_CAPACITY_MB", "100000"); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("postgres://postgres:postgres@localhost:%d/postgres?sslmode=disable", port)
	if err := os.Setenv("POSTGRES_URL", url); err != nil {
		return nil, err
	}
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(port).Logger(nil))
	if err := postgres.Start(); err != nil {
		return nil, err
	}
	m, err := migrate.New("file://fixtures/migrations", url)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil {
		return nil, err
	}
	return postgres, nil
}

func SetupRedis() error {
	s, err := miniredis.Run()
	if err != nil {
		return err
	}
	if err := os.Setenv("REDIS_ADDRESS", s.Addr()); err != nil {
		return err
	}
	return nil
}
