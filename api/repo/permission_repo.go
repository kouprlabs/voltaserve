package repo

import (
	"voltaserve/infra"

	"gorm.io/gorm"
)

type UserPermission struct {
	ID         string `json:"id" gorm:"column:id"`
	UserID     string `json:"userId" gorm:"column:user_id"`
	ResourceID string `json:"resourceId" gorm:"column:resource_id"`
	Permission string `json:"permission" gorm:"column:permission"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}

type GroupPermission struct {
	ID         string `json:"id" gorm:"column:id"`
	GroupID    string `json:"groupId" gorm:"column:group_id"`
	ResourceID string `json:"resourceId" gorm:"column:resource_id"`
	Permission string `json:"permission" gorm:"column:permission"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}

type PermissionRepo interface {
	GetUserPermissions(id string) ([]*UserPermission, error)
	GetGroupPermissions(id string) ([]*GroupPermission, error)
}

func NewPermissionRepo() PermissionRepo {
	return newPermissionRepo()
}

func (UserPermission) TableName() string {
	return "userpermission"
}

func (GroupPermission) TableName() string {
	return "grouppermission"
}

type userPermissionValue struct {
	UserID string `json:"userId,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (p userPermissionValue) GetUserID() string {
	return p.UserID
}
func (p userPermissionValue) GetValue() string {
	return p.Value
}

type groupPermissionValue struct {
	GroupID string `json:"groupId,omitempty"`
	Value   string `json:"value,omitempty"`
}

func (p groupPermissionValue) GetGroupID() string {
	return p.GroupID
}
func (p groupPermissionValue) GetValue() string {
	return p.Value
}

type permissionRepo struct {
	db *gorm.DB
}

func newPermissionRepo() *permissionRepo {
	return &permissionRepo{
		db: infra.GetDb(),
	}
}

func (repo *permissionRepo) GetUserPermissions(id string) ([]*UserPermission, error) {
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

func (repo *permissionRepo) GetGroupPermissions(id string) ([]*GroupPermission, error) {
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
