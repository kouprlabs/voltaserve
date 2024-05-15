package repo

import (
	"time"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type userPermissionEntity struct {
	ID         string `json:"id" gorm:"column:id"`
	UserID     string `json:"userId" gorm:"column:user_id"`
	ResourceID string `json:"resourceId" gorm:"column:resource_id"`
	Permission string `json:"permission" gorm:"column:permission"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}

func (*userPermissionEntity) TableName() string {
	return "userpermission"
}

func (u *userPermissionEntity) BeforeCreate(*gorm.DB) (err error) {
	u.CreateTime = time.Now().UTC().Format(time.RFC3339)
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
	ID         string `json:"id" gorm:"column:id"`
	GroupID    string `json:"groupId" gorm:"column:group_id"`
	ResourceID string `json:"resourceId" gorm:"column:resource_id"`
	Permission string `json:"permission" gorm:"column:permission"`
	CreateTime string `json:"createTime" gorm:"column:create_time"`
}

func (*groupPermissionEntity) TableName() string {
	return "grouppermission"
}

func (g *groupPermissionEntity) BeforeCreate(*gorm.DB) (err error) {
	g.CreateTime = time.Now().UTC().Format(time.RFC3339)
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

type PermissionRepo interface {
	GetUserPermissions(id string) ([]model.UserPermission, error)
	GetGroupPermissions(id string) ([]model.GroupPermission, error)
}

func NewPermissionRepo() PermissionRepo {
	return newPermissionRepo()
}

func NewUserPermission() model.UserPermission {
	return &userPermissionEntity{}
}

func NewGroupPermission() model.GroupPermission {
	return &groupPermissionEntity{}
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

type permissionRepo struct {
	db *gorm.DB
}

func newPermissionRepo() *permissionRepo {
	return &permissionRepo{
		db: infra.NewPostgresManager().GetDBOrPanic(),
	}
}

func (repo *permissionRepo) GetUserPermissions(id string) ([]model.UserPermission, error) {
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

func (repo *permissionRepo) GetGroupPermissions(id string) ([]model.GroupPermission, error) {
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
