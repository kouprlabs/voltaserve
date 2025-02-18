// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package repo

import (
	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type userPermissionEntity struct {
	ID         string `gorm:"column:id"          json:"id"`
	UserID     string `gorm:"column:user_id"     json:"userId"`
	ResourceID string `gorm:"column:resource_id" json:"resourceId"`
	Permission string `gorm:"column:permission"  json:"permission"`
	CreateTime string `gorm:"column:create_time" json:"createTime"`
}

func (*userPermissionEntity) TableName() string {
	return "userpermission"
}

func (u *userPermissionEntity) BeforeCreate(*gorm.DB) (err error) {
	u.CreateTime = helper.NewTimeString()
	return nil
}

func (u *userPermissionEntity) GetID() string {
	return u.ID
}

func (u *userPermissionEntity) GetUserID() string {
	return u.UserID
}

func (u *userPermissionEntity) GetResourceID() string {
	return u.ResourceID
}

func (u *userPermissionEntity) GetPermission() string {
	return u.Permission
}

func (u *userPermissionEntity) GetCreateTime() string {
	return u.CreateTime
}

func (u *userPermissionEntity) SetID(id string) {
	u.ID = id
}

func (u *userPermissionEntity) SetUserID(userID string) {
	u.UserID = userID
}

func (u *userPermissionEntity) SetResourceID(resourceID string) {
	u.ResourceID = resourceID
}

func (u *userPermissionEntity) SetPermission(permission string) {
	u.Permission = permission
}

func (u *userPermissionEntity) SetCreateTime(createTime string) {
	u.CreateTime = createTime
}

type groupPermissionEntity struct {
	ID         string `gorm:"column:id"          json:"id"`
	GroupID    string `gorm:"column:group_id"    json:"groupId"`
	ResourceID string `gorm:"column:resource_id" json:"resourceId"`
	Permission string `gorm:"column:permission"  json:"permission"`
	CreateTime string `gorm:"column:create_time" json:"createTime"`
}

func (*groupPermissionEntity) TableName() string {
	return "grouppermission"
}

func (g *groupPermissionEntity) BeforeCreate(*gorm.DB) (err error) {
	g.CreateTime = helper.NewTimeString()
	return nil
}

func (g *groupPermissionEntity) GetID() string {
	return g.ID
}

func (g *groupPermissionEntity) GetGroupID() string {
	return g.GroupID
}

func (g *groupPermissionEntity) GetResourceID() string {
	return g.ResourceID
}

func (g *groupPermissionEntity) GetPermission() string {
	return g.Permission
}

func (g *groupPermissionEntity) GetCreateTime() string {
	return g.CreateTime
}

func (g *groupPermissionEntity) SetID(id string) {
	g.ID = id
}

func (g *groupPermissionEntity) SetGroupID(groupID string) {
	g.GroupID = groupID
}

func (g *groupPermissionEntity) SetResourceID(resourceID string) {
	g.ResourceID = resourceID
}

func (g *groupPermissionEntity) SetPermission(permission string) {
	g.Permission = permission
}

func (g *groupPermissionEntity) SetCreateTime(createTime string) {
	g.CreateTime = createTime
}

type UserPermissionValue struct {
	UserID string `json:"userId,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (p UserPermissionValue) GetUserID() string {
	return p.UserID
}

func (p UserPermissionValue) GetValue() string {
	return p.Value
}

type GroupPermissionValue struct {
	GroupID string `json:"groupId,omitempty"`
	Value   string `json:"value,omitempty"`
}

func (p GroupPermissionValue) GetGroupID() string {
	return p.GroupID
}

func (p GroupPermissionValue) GetValue() string {
	return p.Value
}

func NewUserPermissionModel() model.UserPermission {
	return &userPermissionEntity{}
}

type UserPermissionNewModelOptions struct {
	ID         string
	UserID     string
	ResourceID string
	Permission string
	CreateTime string
}

func NewUserPermissionNewModel(opts UserPermissionNewModelOptions) model.UserPermission {
	return &userPermissionEntity{
		ID:         opts.ID,
		UserID:     opts.UserID,
		ResourceID: opts.ResourceID,
		Permission: opts.Permission,
		CreateTime: opts.CreateTime,
	}
}

func NewGroupPermissionModel() model.GroupPermission {
	return &groupPermissionEntity{}
}

type GroupPermissionNewModelOptions struct {
	ID         string
	GroupID    string
	ResourceID string
	Permission string
	CreateTime string
}

func NewGroupPermissionModelWithOptions(opts GroupPermissionNewModelOptions) model.GroupPermission {
	return &groupPermissionEntity{
		ID:         opts.ID,
		GroupID:    opts.GroupID,
		ResourceID: opts.ResourceID,
		Permission: opts.Permission,
		CreateTime: opts.CreateTime,
	}
}

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo() *PermissionRepo {
	return &PermissionRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *PermissionRepo) FindUserPermissions(id string) ([]model.UserPermission, error) {
	var entities []*userPermissionEntity
	if db := repo.db.
		Raw("SELECT * FROM userpermission WHERE resource_id = ?", id).
		Scan(&entities); db.Error != nil {
		return nil, db.Error
	}
	if len(entities) > 0 {
		var res []model.UserPermission
		for _, entity := range entities {
			res = append(res, entity)
		}
		return res, nil
	} else {
		return nil, nil
	}
}

func (repo *PermissionRepo) FindGroupPermissions(id string) ([]model.GroupPermission, error) {
	var entities []*groupPermissionEntity
	if db := repo.db.
		Raw("SELECT * FROM grouppermission WHERE resource_id = ?", id).
		Scan(&entities); db.Error != nil {
		return nil, db.Error
	}
	if len(entities) > 0 {
		var res []model.GroupPermission
		for _, entity := range entities {
			res = append(res, entity)
		}
		return res, nil
	} else {
		return nil, nil
	}
}
