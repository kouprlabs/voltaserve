package repo

import (
	"voltaserve/infra"

	"gorm.io/gorm"
)

func (UserPermission) TableName() string {
	return "userpermission"
}

func (GroupPermission) TableName() string {
	return "grouppermission"
}

type PostgresPermissionRepo struct {
	db *gorm.DB
}

func NewPostgresPermissionRepo() *PostgresPermissionRepo {
	return &PostgresPermissionRepo{
		db: infra.GetDb(),
	}
}

func (repo *PostgresPermissionRepo) GetUserPermissions(id string) ([]*UserPermission, error) {
	var res []*UserPermission
	if db := repo.db.
		Raw("SELECT * FROM userpermission WHERE resource_id = ?", id).
		Scan(&res); db.Error != nil {
		return nil, db.Error
	}
	if len(res) > 0 {
		return res, nil
	} else {
		return []*UserPermission{}, nil
	}
}

func (repo *PostgresPermissionRepo) GetGroupPermissions(id string) ([]*GroupPermission, error) {
	var res []*GroupPermission
	if db := repo.db.
		Raw("SELECT * FROM grouppermission WHERE resource_id = ?", id).
		Scan(&res); db.Error != nil {
		return nil, db.Error
	}
	if len(res) > 0 {
		return res, nil
	} else {
		return []*GroupPermission{}, nil
	}
}
