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
	"os"

	"github.com/alicebob/miniredis/v2"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"

	"github.com/kouprlabs/voltaserve/api/config"
)

func setupPostgres() (*embeddedpostgres.EmbeddedPostgres, error) {
	os.Setenv("POSTGRES_URL", "postgres://postgres:postgres@localhost:15432/postgres?sslmode=disable")
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(15432).Logger(nil))
	if err := postgres.Start(); err != nil {
		return nil, err
	}
	m, err := migrate.New("file://../test/migrations", config.GetConfig().DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil {
		return nil, err
	}
	return postgres, nil
}

func setupRedis() error {
	s, err := miniredis.Run()
	if err != nil {
		return err
	}
	os.Setenv("REDIS_ADDRESS", s.Addr())
	return nil
}
