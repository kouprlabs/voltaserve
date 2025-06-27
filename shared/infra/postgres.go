// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kouprlabs/voltaserve/shared/config"
)

var db *gorm.DB

type PostgresManager struct {
	postgresConfig config.PostgresConfig
	envConfig      config.EnvironmentConfig
}

func NewPostgresManager(postgresConfig config.PostgresConfig, envConfig config.EnvironmentConfig) *PostgresManager {
	return &PostgresManager{
		postgresConfig: postgresConfig,
		envConfig:      envConfig,
	}
}

func (mgr *PostgresManager) Connect(ignoreExisting bool) error {
	if !ignoreExisting && db != nil {
		return nil
	}

	sqlDB, err := sql.Open("pgx", mgr.postgresConfig.URL)
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(mgr.postgresConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(mgr.postgresConfig.MaxOpenConnections)
	sqlDB.SetConnMaxIdleTime(time.Duration(mgr.postgresConfig.ConnectionMaxIdleTimeMinutes) * time.Minute)

	opts := &gorm.Config{}
	if mgr.envConfig.IsTest {
		opts.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), opts)
	if err != nil {
		return err
	}

	go func() {
		t := time.NewTicker(time.Duration(mgr.postgresConfig.KeepAliveIntervalMinutes) * time.Minute)
		for range t.C {
			db.Exec("SELECT 1")
		}
	}()

	return nil
}

func (mgr *PostgresManager) GetDB() (*gorm.DB, error) {
	if db == nil {
		if err := mgr.Connect(false); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func (mgr *PostgresManager) GetDBOrPanic() *gorm.DB {
	db, err := mgr.GetDB()
	if err != nil {
		panic(err)
	}
	return db
}
