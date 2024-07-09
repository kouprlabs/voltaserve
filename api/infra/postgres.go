// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/config"
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
