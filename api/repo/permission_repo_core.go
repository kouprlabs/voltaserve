package repo

import (
	"fmt"
	"voltaserve/config"
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

type CorePermissionRepo interface {
	GetUserPermissions(id string) ([]*UserPermission, error)
	GetGroupPermissions(id string) ([]*GroupPermission, error)
}

func NewPermissionRepo() CorePermissionRepo {
	if config.GetConfig().DatabaseType == config.DATABASE_TYPE_POSTGRES {
		return NewPostgresPermissionRepo()
	}
	panic(fmt.Sprintf("database type %s repo not implemented", config.GetConfig().DatabaseType))
}
