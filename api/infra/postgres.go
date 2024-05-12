package infra

import (
	"voltaserve/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type PostgresManager struct {
}

func NewPostgresManager() *PostgresManager {
	return &PostgresManager{}
}

func (mgr *PostgresManager) Connect(ignoreExisting bool) error {
	if !ignoreExisting && db != nil {
		return nil
	}
	var err error
	db, err = gorm.Open(postgres.Open(config.GetConfig().DatabaseURL), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func (mgr *PostgresManager) GetDB() (*gorm.DB, error) {
	if err := mgr.Connect(false); err != nil {
		return nil, err
	}
	return db, nil
}

func (mgr *PostgresManager) GetDBOrPanic() *gorm.DB {
	if err := mgr.Connect(false); err != nil {
		panic(err)
	}
	return db
}
