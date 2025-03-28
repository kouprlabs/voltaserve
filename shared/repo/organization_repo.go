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
	"errors"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type organizationEntity struct {
	ID               string                  `gorm:"column:id"          json:"id"`
	Name             string                  `gorm:"column:name"        json:"name"`
	UserPermissions  []*UserPermissionValue  `gorm:"-"                  json:"userPermissions"`
	GroupPermissions []*GroupPermissionValue `gorm:"-"                  json:"groupPermissions"`
	Members          []string                `gorm:"-"                  json:"members"`
	CreateTime       string                  `gorm:"column:create_time" json:"createTime"`
	UpdateTime       *string                 `gorm:"column:update_time" json:"updateTime,omitempty"`
}

func (*organizationEntity) TableName() string {
	return "organization"
}

func (o *organizationEntity) BeforeCreate(*gorm.DB) (err error) {
	o.CreateTime = helper.NewTimeString()
	return nil
}

func (o *organizationEntity) BeforeSave(*gorm.DB) (err error) {
	o.UpdateTime = helper.ToPtr(helper.NewTimeString())
	return nil
}

func (o *organizationEntity) GetID() string {
	return o.ID
}

func (o *organizationEntity) GetName() string {
	return o.Name
}

func (o *organizationEntity) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range o.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (o *organizationEntity) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range o.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (o *organizationEntity) GetMembers() []string {
	return o.Members
}

func (o *organizationEntity) GetCreateTime() string {
	return o.CreateTime
}

func (o *organizationEntity) GetUpdateTime() *string {
	return o.UpdateTime
}

func (o *organizationEntity) SetID(id string) {
	o.ID = id
}

func (o *organizationEntity) SetName(name string) {
	o.Name = name
}

func (o *organizationEntity) SetUserPermissions(permissions []model.CoreUserPermission) {
	o.UserPermissions = make([]*UserPermissionValue, len(permissions))
	for i, p := range permissions {
		o.UserPermissions[i] = p.(*UserPermissionValue)
	}
}

func (o *organizationEntity) SetGroupPermissions(permissions []model.CoreGroupPermission) {
	o.GroupPermissions = make([]*GroupPermissionValue, len(permissions))
	for i, p := range permissions {
		o.GroupPermissions[i] = p.(*GroupPermissionValue)
	}
}

func (o *organizationEntity) SetCreateTime(createTime string) {
	o.CreateTime = createTime
}

func (o *organizationEntity) SetUpdateTime(updateTime *string) {
	o.UpdateTime = updateTime
}

func NewOrganizationModel() model.Organization {
	return &organizationEntity{}
}

type OrganizationNewModelOptions struct {
	ID               string
	Name             string
	UserPermissions  []model.CoreUserPermission
	GroupPermissions []model.CoreGroupPermission
	Members          []string
	CreateTime       string
	UpdateTime       *string
}

func NewOrganizationModelWithOptions(opts OrganizationNewModelOptions) model.Organization {
	res := &organizationEntity{
		ID:         opts.ID,
		Name:       opts.Name,
		Members:    opts.Members,
		CreateTime: opts.CreateTime,
		UpdateTime: opts.UpdateTime,
	}
	res.SetUserPermissions(opts.UserPermissions)
	res.SetGroupPermissions(opts.GroupPermissions)
	return res
}

type OrganizationRepo struct {
	db             *gorm.DB
	groupRepo      *GroupRepo
	permissionRepo *PermissionRepo
}

func NewOrganizationRepo(postgres config.PostgresConfig, environment config.EnvironmentConfig) *OrganizationRepo {
	return &OrganizationRepo{
		db:             infra.NewPostgresManager(postgres, environment).GetDBOrPanic(),
		groupRepo:      NewGroupRepo(postgres, environment),
		permissionRepo: NewPermissionRepo(postgres, environment),
	}
}

type OrganizationInsertOptions struct {
	ID   string
	Name string
}

func (repo *OrganizationRepo) Insert(opts OrganizationInsertOptions) (model.Organization, error) {
	org := organizationEntity{
		ID:   opts.ID,
		Name: opts.Name,
	}
	if db := repo.db.Create(&org); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.Find(opts.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *OrganizationRepo) Find(id string) (model.Organization, error) {
	org, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*organizationEntity{org}); err != nil {
		return nil, err
	}
	return org, nil
}

func (repo *OrganizationRepo) FindOrNil(id string) model.Organization {
	res, err := repo.Find(id)
	if err != nil {
		return nil
	}
	return res
}

func (repo *OrganizationRepo) FindIDs() ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw("SELECT id result FROM organization ORDER BY create_time DESC").Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *OrganizationRepo) FindIDsByOwner(userID string) ([]string, error) {
	type IDResult struct {
		Result string
	}
	var ids []IDResult
	db := repo.db.
		Raw(`SELECT id result FROM organization 
			 WHERE id IN (SELECT resource_id FROM userpermission WHERE user_id = ? AND permission = 'owner')`,
			userID).Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := make([]string, 0)
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *OrganizationRepo) FindMembers(id string) ([]model.User, error) {
	var entities []*userEntity
	db := repo.db.
		Raw(`SELECT u.* FROM "user" u INNER JOIN userpermission up on
             u.id = up.user_id AND up.resource_id = ?`,
			id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.User
	for _, u := range entities {
		res = append(res, u)
	}
	return res, nil
}

func (repo *OrganizationRepo) FindGroups(id string) ([]model.Group, error) {
	var entities []*groupEntity
	db := repo.db.
		Raw(`SELECT * FROM "group" g WHERE g.organization_id = ? ORDER BY g.name`, id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.groupRepo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.Group
	for _, g := range entities {
		res = append(res, g)
	}
	return res, nil
}

func (repo *OrganizationRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&organizationEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *OrganizationRepo) Save(org model.Organization) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *OrganizationRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM organization WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	db = repo.db.Exec("DELETE FROM userpermission WHERE resource_id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	db = repo.db.Exec("DELETE FROM grouppermission WHERE resource_id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *OrganizationRepo) CountOwners(id string) (int64, error) {
	var count int64
	db := repo.db.Model(&userPermissionEntity{}).
		Where("resource_id = ?", id).
		Where("permission = ?", model.PermissionOwner).
		Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *OrganizationRepo) GrantUserPermission(id string, userID string, permission string) error {
	db := repo.db.
		Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
              VALUES (?, ?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?`,
			helper.NewID(), userID, id, permission, helper.NewTimeString(), permission)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *OrganizationRepo) RevokeUserPermission(id string, userID string) error {
	db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userID, id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *OrganizationRepo) find(id string) (*organizationEntity, error) {
	res := organizationEntity{}
	db := repo.db.Where("id = ?", id).First(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewOrganizationNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *OrganizationRepo) populateModelFields(organizations []*organizationEntity) error {
	for _, o := range organizations {
		o.UserPermissions = make([]*UserPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.FindUserPermissions(o.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			o.UserPermissions = append(o.UserPermissions, &UserPermissionValue{
				UserID: p.GetUserID(),
				Value:  p.GetPermission(),
			})
		}
		o.GroupPermissions = make([]*GroupPermissionValue, 0)
		groupPermissions, err := repo.permissionRepo.FindGroupPermissions(o.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			o.GroupPermissions = append(o.GroupPermissions, &GroupPermissionValue{
				GroupID: p.GetGroupID(),
				Value:   p.GetPermission(),
			})
		}
		members, err := repo.FindMembers(o.ID)
		if err != nil {
			return nil
		}
		o.Members = make([]string, 0)
		for _, u := range members {
			o.Members = append(o.Members, u.GetID())
		}
	}
	return nil
}
