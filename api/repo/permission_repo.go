package repo

import (
	"voltaserve/infra"

	"gorm.io/gorm"
)

type UserPermission struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	ResourceId string `json:"resourceId"`
	Permission string `json:"permission"`
	CreateTime string `json:"createTime"`
}

func (UserPermission) TableName() string {
	return "userpermission"
}

type GroupPermission struct {
	Id         string `json:"id"`
	GroupId    string `json:"groupId"`
	ResourceId string `json:"resourceId"`
	Permission string `json:"permission"`
	CreateTime string `json:"createTime"`
}

func (GroupPermission) TableName() string {
	return "grouppermission"
}

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo() *PermissionRepo {
	return &PermissionRepo{
		db: infra.GetDb(),
	}
}

func (repo *PermissionRepo) GetUserPermissions(id string) ([]*UserPermission, error) {
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

func (repo *PermissionRepo) GetGroupPermissions(id string) ([]*GroupPermission, error) {
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
