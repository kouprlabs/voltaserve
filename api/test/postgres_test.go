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
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

func TestPostgres(t *testing.T) {
	orgRepo := repo.NewOrganizationRepo()
	org, err := orgRepo.Insert(repo.OrganizationInsertOptions{
		ID:   helper.NewID(),
		Name: "organization",
	})
	if err != nil {
		t.Fatal(err)
	}
	workspaceRepo := repo.NewWorkspaceRepo()
	workspace, err := workspaceRepo.Insert(repo.WorkspaceInsertOptions{
		ID:             helper.NewID(),
		Name:           "wokrspace",
		OrganizationID: org.GetID(),
	})
	if err != nil {
		t.Fatal(err)
	}
	fileRepo := repo.NewFileRepo()
	file, err := fileRepo.Insert(repo.FileInsertOptions{
		Name:        "file",
		Type:        model.FileTypeFile,
		WorkspaceID: workspace.GetID(),
	})
	if err != nil {
		t.Fatal(err)
	}
	foundFile, err := fileRepo.Find(file.GetID())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, file.GetID(), foundFile.GetID())
}
