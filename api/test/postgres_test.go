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
	"github.com/stretchr/testify/suite"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type PostgresSuite struct {
	suite.Suite
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}

func (s *PostgresSuite) TestInsertAndFind() {
	org, err := repo.NewOrganizationRepo().Insert(repo.OrganizationInsertOptions{
		ID:   helper.NewID(),
		Name: "organization",
	})
	s.Require().NoError(err)

	workspace, err := repo.NewWorkspaceRepo().Insert(repo.WorkspaceInsertOptions{
		ID:             helper.NewID(),
		Name:           "workspace",
		OrganizationID: org.GetID(),
	})
	s.Require().NoError(err)

	file, err := repo.NewFileRepo().Insert(repo.FileInsertOptions{
		Name:        "file",
		Type:        model.FileTypeFile,
		WorkspaceID: workspace.GetID(),
	})
	s.Require().NoError(err)

	foundFile, err := repo.NewFileRepo().Find(file.GetID())
	s.Require().NoError(err)
	s.Equal(file.GetID(), foundFile.GetID())
}
