package repo

import (
	"voltaserve/infra"

	"gorm.io/gorm"
)

type UserPermission struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	ResourceID string `json:"resourceId"`
	Permission string `json:"permission"`
	CreateTime string `json:"createTime"`
}

type GroupPermission struct {
	ID         string `json:"id"`
	GroupID    string `json:"groupId"`
	ResourceId string `json:"resourceId"`
	Permission string `json:"permission"`
	CreateTime string `json:"createTime"`
}

type PermissionRepo interface {
	GetUserPermissions(id string) ([]*UserPermission, error)
	GetGroupPermissions(id string) ([]*GroupPermission, error)
}

func NewPermissionRepo() PermissionRepo {
	return NewPostgresPermissionRepo()
}

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
