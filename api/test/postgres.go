package test

import (
	"fmt"
	"os"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"

	"github.com/kouprlabs/voltaserve/api/config"
)

func SetupPostgres(port uint32) (*embeddedpostgres.EmbeddedPostgres, error) {
	os.Setenv("POSTGRES_URL", fmt.Sprintf("postgres://postgres:postgres@localhost:%d/postgres?sslmode=disable", port))
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(port).Logger(nil))
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
