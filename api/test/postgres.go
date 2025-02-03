package test

import (
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"

	"github.com/kouprlabs/voltaserve/api/config"
)

func SetupPostgres(t *testing.T) *embeddedpostgres.EmbeddedPostgres {
	t.Helper()
	t.Setenv("POSTGRES_URL", "postgres://postgres:postgres@localhost:15432/postgres?sslmode=disable")
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(15432).Logger(nil))
	if err := postgres.Start(); err != nil {
		t.Fatal(err)
	}
	m, err := migrate.New("file://../test/migrations", config.GetConfig().DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Up(); err != nil {
		t.Fatal(err)
	}
	return postgres
}
