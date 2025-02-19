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
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Setup struct {
	env      *Env
	redis    *Redis
	postgres *Postgres
	m        *testing.M
}

func NewSetup(m *testing.M) *Setup {
	return &Setup{
		m:        m,
		env:      NewEnv(),
		redis:    NewRedis(),
		postgres: NewPostgres(),
	}
}

type SetupOptions struct {
	Postgres PostgresOptions
}

func (s *Setup) Up(opts SetupOptions) error {
	if err := s.env.Apply(); err != nil {
		return err
	}
	if err := s.redis.Start(); err != nil {
		return err
	}
	err := s.postgres.Start(opts.Postgres)
	if err != nil {
		return err
	}
	return nil
}

func (s *Setup) Down() error {
	if err := s.postgres.Stop(); err != nil {
		return err
	}
	return nil
}
