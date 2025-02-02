package repo_test

import (
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

func TestPostgres(t *testing.T) {
	t.Setenv("PORT", "0")
	t.Setenv("POSTGRES_URL", "postgres://postgres:postgres@localhost:15432/postgres?sslmode=disable")
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(15432))
	if err := postgres.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := postgres.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
	m, err := migrate.New("file://./migrations", "postgres://postgres:postgres@localhost:15432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}
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
		ID:              helper.NewID(),
		Name:            "wokrspace",
		StorageCapacity: 100000,
		OrganizationID:  org.GetID(),
		Bucket:          "bucket",
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
