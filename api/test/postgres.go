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
	"path"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kouprlabs/voltaserve/api/helper"
)

type Postgres struct {
	postgres *embeddedpostgres.EmbeddedPostgres
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

type PostgresOptions struct {
	Port uint32
}

func (p *Postgres) Start(opts PostgresOptions) error {
	url := fmt.Sprintf("postgres://postgres:postgres@localhost:%d/postgres?sslmode=disable", opts.Port)
	if err := os.Setenv("POSTGRES_URL", url); err != nil {
		return err
	}
	p.postgres = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(opts.Port).
		Logger(nil).
		RuntimePath(path.Join(os.TempDir(), helper.NewID())),
	)
	if err := p.postgres.Start(); err != nil {
		return err
	}
	m, err := migrate.New("file://fixtures/migrations", url)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Stop() error {
	if err := p.postgres.Stop(); err != nil {
		return err
	}
	return nil
}
