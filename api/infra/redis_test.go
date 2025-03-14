// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra_test

import (
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"
	"github.com/kouprlabs/voltaserve/shared/repo"

	"github.com/kouprlabs/voltaserve/api/config"
)

type RedisSuite struct {
	suite.Suite
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}

func (s *RedisSuite) TestSetAndGet() {
	opts := repo.FileNewModelOptions{
		ID:   helper.NewID(),
		Name: "file",
		Type: model.FileTypeFile,
	}
	err := cache.NewFileCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Set(repo.NewFileModelWithOptions(opts))
	s.Require().NoError(err)

	file, err := cache.NewFileCache(
		config.GetConfig().Postgres,
		config.GetConfig().Redis,
		config.GetConfig().Environment,
	).Get(opts.ID)
	s.Require().NoError(err)
	s.Equal(opts.ID, file.GetID())
	s.Equal(opts.Name, file.GetName())
	s.Equal(opts.Type, file.GetType())
}
