package repo_test

import (
	"fmt"
	"os"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
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
	gormDB := infra.NewPostgresManager().GetDBOrPanic()
	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cwd)
	if err := goose.Up(sqlDB, "./migrations"); err != nil {
		t.Fatal(err)
	}

	orgRepo := repo.NewOrganizationRepoWithDB(gormDB)
	org, err := orgRepo.Insert(repo.OrganizationInsertOptions{
		ID:   helper.NewID(),
		Name: "organization",
	})
	if err != nil {
		t.Fatal(err)
	}

	workspaceRepo := repo.NewWorkspaceRepoWithDB(gormDB)
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

	fileRepo := repo.NewFileRepoWithDB(gormDB)
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
